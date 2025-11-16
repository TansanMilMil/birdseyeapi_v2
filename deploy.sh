#!/bin/bash -eu

cd $(dirname $0)

NO_SERVE=false
INIT_DB=false

for ARG in "$@"; do
    case "$ARG" in
    --no-serve)
        NO_SERVE=true
        ;;
    --init-db)
        INIT_DB=true
        ;;
    esac
done

echo "load .env ----------------------------"
if [ ! -f .env ]; then
    echo ".env file not found!"
else
    source .env
fi

VENUS_TARGET_DIR=$VENUS_HOME/birdseyeapi_v2

echo "docker compose setup ----------------------------"
docker compose down
docker compose -f ./docker-compose.ci.yml up -d go

echo "run all tests ----------------------------"
docker compose exec go ./test.sh

echo "build ----------------------------"
docker compose exec go ./build.sh

echo "docker compose teardown ----------------------------"
docker compose down

echo "create deploy package ----------------------------"
rm -rf ./deploy-temp
mkdir -p ./deploy-temp
cp -r ./go/dist ./deploy-temp/
cp ./docker-compose.yml ./deploy-temp/
cp ./go-entrypoint.sh ./deploy-temp/
cp ./init_db.sh ./deploy-temp/
cp -r ./nginx ./deploy-temp/
cp -r ./mysql ./deploy-temp/
tar czf ./deploy.tgz -C ./deploy-temp .
rm -rf ./deploy-temp

echo "scp files to venus ----------------------------"
ssh $VENUS_SSH_HOST mkdir -p $VENUS_TARGET_DIR
scp ./deploy.tgz $VENUS_SSH_HOST:$VENUS_TARGET_DIR/deploy.tgz
ssh $VENUS_SSH_HOST "tar xzf $VENUS_TARGET_DIR/deploy.tgz -C $VENUS_TARGET_DIR && rm $VENUS_TARGET_DIR/deploy.tgz"
rm ./deploy.tgz
echo ""
echo "current files in $VENUS_SSH_HOST:$VENUS_TARGET_DIR"
ssh $VENUS_SSH_HOST ls -alh $VENUS_TARGET_DIR

echo "stop birdseyeapi ----------------------------"
ssh $VENUS_SSH_HOST docker compose -f $VENUS_TARGET_DIR/docker-compose.yml down

if [ "$NO_SERVE" = "false" ]; then
    echo "start birdseyeapi ----------------------------"
    ssh $VENUS_SSH_HOST docker compose -f $VENUS_TARGET_DIR/docker-compose.yml up -d
fi

if $INIT_DB; then
    echo "init db ----------------------------"
    ssh $VENUS_SSH_HOST docker compose -f $VENUS_TARGET_DIR/docker-compose.yml up -d
    # MySQLが起動するまでウェイトする
    echo "Waiting for MySQL to start..."
    sleep 5
    ssh $VENUS_SSH_HOST $VENUS_TARGET_DIR/init_db.sh
fi

echo ""
echo "current docker containers"
echo "Waiting for 10 seconds to get the latest status..."
sleep 10
ssh $VENUS_SSH_HOST docker ps

echo "deploy finished!! ----------------------------"

echo "curl health check ----------------------------"
echo "Waiting for 20 seconds to get the latest status..."
sleep 20
./check-remote-curl.sh
