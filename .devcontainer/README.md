# DevContainer 使用方法

## 開発環境の起動

1. VS Codeでプロジェクトを開く
2. コマンドパレット（Ctrl+Shift+P）で「Dev Containers: Reopen in Container」を選択
3. DevContainerが起動するまで待機

## サービスの起動

DevContainer内で以下のコマンドを実行してインフラサービスを起動：

```bash
docker compose up -d
```

これにより以下のサービスが起動します：
- MySQL (ポート: 3307)
- Selenium (ポート: 4444)
- Nginx (ポート: 1111)

## Goアプリケーションの実行

```bash
# 依存関係のインストール
go mod download

# アプリケーションの起動
go run go/main.go
```

または既存のスクリプトを使用：

```bash
./run.sh
```

## ポートフォワーディング

以下のポートが自動的にフォワードされます：
- 8080: Go API サーバー
- 1111: Nginx
- 3307: MySQL
- 4444: Selenium

## 開発環境の特徴

- Go 1.22 開発環境
- Docker-in-Docker サポート
- Go言語拡張機能（LSP、デバッガー、フォーマッター）
- 自動保存時のコードフォーマット
- MySQL クライアントツール

## 環境変数

`.env`ファイルを作成して必要な環境変数を設定してください。
`.env.example`を参考にしてください。