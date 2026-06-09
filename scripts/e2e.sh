#!/bin/bash -eu

# ============================================================================
# ローカル全コンテナ起動 E2E テスト
#
# 全7サービスを PRODUCTION モードで docker compose up し、以下を assert する:
#   1. 全サービスが running で居続けている (再起動ループしていない)
#   2. HTTP 経由でアプリが疎通する (nginx /HealthCheck と go+DB /news/today-news)
#   3. 各コンテナのログに致命的エラーが出ていない
#
# 既存 check-remote-curl.sh の作法 (.env source + curl) に倣う。
# ============================================================================

# --- 準備 -------------------------------------------------------------------
# リポジトリルート (このスクリプトの一つ上) へ移動
cd "$(dirname "$0")/.."

echo "load .env ----------------------------"
if [ ! -f .env ]; then
    echo ".env file not found!"
    exit 1
fi
source .env

# ビルド済みバイナリが無いと PRODUCTION モードの entrypoint が起動できない
if [ ! -f go/dist/birdseyeapi_v2 ]; then
    echo "ERROR: go/dist/birdseyeapi_v2 が見つかりません。"
    echo "       先に 'task build' を実行してバイナリをビルドしてください。"
    exit 1
fi

# 期待サービス (7つ)
EXPECTED_SERVICES="selenium go nginx mysql loki promtail grafana"

# --- ログダンプ (失敗時の調査用) --------------------------------------------
dump_logs() {
    echo ""
    echo "########## E2E FAILED: dumping container logs ##########"
    for svc in $EXPECTED_SERVICES; do
        echo ""
        echo "===== logs: ${svc} ====="
        docker compose logs --no-color "$svc" 2>&1 || true
    done
    echo ""
    echo "===== docker compose ps ====="
    docker compose ps 2>&1 || true
    echo "#######################################################"
}

# --- cleanup 登録 -----------------------------------------------------------
# EXIT 時に必ず後始末。失敗 (rc != 0) の場合は down 前にログをダンプ。
# down のみでボリューム (db-store 等) は温存し、ローカル開発データを壊さない。
cleanup() {
    rc=$?
    if [ "$rc" -ne 0 ]; then
        dump_logs
    fi
    echo ""
    echo "==== teardown: docker compose down ===="
    docker compose down || true
    exit "$rc"
}
trap cleanup EXIT

# --- 全サービス running の assert -------------------------------------------
assert_all_running() {
    echo ""
    echo "==== assert: all services running ===="

    # compose v2 の出力差 (NDJSON / 配列) を jq -s 'flatten' で吸収
    local ps_json
    ps_json="$(docker compose ps --format json | jq -s 'flatten')"

    for svc in $EXPECTED_SERVICES; do
        local state
        state="$(echo "$ps_json" | jq -r --arg s "$svc" \
            '.[] | select(.Service == $s) | .State' | head -n1)"

        if [ -z "$state" ]; then
            echo "  [FAIL] service '${svc}' がコンテナ一覧に存在しません"
            return 1
        fi
        if [ "$state" != "running" ]; then
            echo "  [FAIL] service '${svc}' の State == '${state}' (expected: running)"
            return 1
        fi

        # クラッシュループ検知: restart: unless-stopped で running に見えても
        # 再起動を繰り返しているケースを RestartCount で除外する
        local cid restarts
        cid="$(docker compose ps -q "$svc")"
        if [ -z "$cid" ]; then
            echo "  [FAIL] service '${svc}' のコンテナIDを取得できません"
            return 1
        fi
        restarts="$(docker inspect -f '{{.RestartCount}}' "$cid")"
        if [ "$restarts" != "0" ]; then
            echo "  [FAIL] service '${svc}' が再起動を繰り返しています (RestartCount=${restarts})"
            return 1
        fi

        echo "  [OK]   ${svc}: running (RestartCount=0)"
    done
}

# --- HTTP 疎通待ち (healthcheck が無いのでリトライ付きポーリング) -----------
# usage: poll_http <url> <description> <success-test-cmd>
#   success-test-cmd は curl の出力(stdin)を受け取り、成功時に exit 0 を返す
poll_http() {
    local url="$1" desc="$2"
    local interval=2 max_attempts=30 attempt body

    echo ""
    echo "==== wait HTTP: ${desc} (${url}) ===="
    for attempt in $(seq 1 "$max_attempts"); do
        if body="$(curl -fsS "$url" 2>/dev/null)"; then
            # 第3引数があれば本文を検証 (HealthCheck の 'ok' 等)
            if [ "$#" -lt 3 ] || echo "$body" | eval "$3"; then
                echo "  [OK]   ${desc} 疎通 (attempt ${attempt}/${max_attempts})"
                return 0
            fi
        fi
        printf '  ... waiting (%d/%d)\n' "$attempt" "$max_attempts"
        sleep "$interval"
    done

    echo "  [FAIL] ${desc} がタイムアウトしました (${url})"
    return 1
}

wait_http() {
    # nginx 疎通: /HealthCheck が 'ok' を返すまで
    poll_http "localhost:1111/HealthCheck" "nginx /HealthCheck" "grep -q '^ok$'"

    # go+DB 疎通 (本命): /news/today-news が HTTP 200 を返せば
    # 「go 起動済み + MySQL 接続成功 + GORM AutoMigrate 成功」が同時に保証される
    poll_http "localhost:1111/news/today-news" "go+DB /news/today-news"
}

# --- 致命ログの assert (致命のみ) -------------------------------------------
# usage: assert_no_pattern <service> <egrep-pattern>
assert_no_pattern() {
    local svc="$1" pattern="$2" hits
    hits="$(docker compose logs --no-color "$svc" 2>&1 | grep -E "$pattern" || true)"
    if [ -n "$hits" ]; then
        echo "  [FAIL] ${svc} のログに致命パターン (${pattern}) を検出:"
        echo "$hits" | sed 's/^/         /'
        return 1
    fi
    echo "  [OK]   ${svc}: 致命ログ無し"
}

assert_no_fatal_logs() {
    echo ""
    echo "==== assert: no fatal logs ===="
    # go: panic / Fatal / log.Fatalf 由来の起動失敗メッセージ
    assert_no_pattern go 'panic:|Fatal|Failed to connect to database|Failed to start server'
    # nginx: 致命寄りの emerg / alert
    assert_no_pattern nginx '\[emerg\]|\[alert\]'
    # mysql: ERROR は致命寄り
    assert_no_pattern mysql '\[ERROR\]'
    # その他: panic / fatal のみ (無害な error/warn は対象外)
    assert_no_pattern selenium 'panic|fatal'
    assert_no_pattern loki 'panic|fatal'
    assert_no_pattern promtail 'panic|fatal'
    assert_no_pattern grafana 'panic|fatal'
}

# --- 起動 -------------------------------------------------------------------
echo ""
echo "==== boot: clean start in PRODUCTION mode ===="
# クリーンに開始
docker compose down
# .env ファイル自体は書き換えず、環境変数のインライン上書きで PRODUCTION を強制
BIRDSEYEAPI_EXECUTION_MODE=PRODUCTION docker compose up -d

# --- 各 assert 実行 ---------------------------------------------------------
assert_all_running
wait_http
assert_no_fatal_logs

# --- 成功 -------------------------------------------------------------------
echo ""
echo "########## E2E PASSED ##########"
# この後 trap cleanup が docker compose down を実行する
