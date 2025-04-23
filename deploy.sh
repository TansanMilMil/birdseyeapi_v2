#!/bin/bash -eu

cd `dirname $0`

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

echo "build ----------------------------"
docker compose down
docker compose -f ./docker-compose.ci.yml up -d go
docker compose exec go ./build.sh --no-docker-compose
tar cvfz ./go/dist.tgz ./go/dist/

echo "docker compose teardown ----------------------------"
docker compose down

echo "scp files to venus ----------------------------"
ssh $VENUS_SSH_HOST mkdir -p $VENUS_TARGET_DIR/go
scp ./go/dist.tgz $VENUS_SSH_HOST:$VENUS_TARGET_DIR/go/dist.tgz
ssh $VENUS_SSH_HOST tar xvfz $VENUS_TARGET_DIR/go/dist.tgz -C $VENUS_TARGET_DIR

scp ./docker-compose.yml $VENUS_SSH_HOST:$VENUS_TARGET_DIR/docker-compose.yml
scp ./go-entrypoint.sh $VENUS_SSH_HOST:$VENUS_TARGET_DIR/go-entrypoint.sh
scp ./init_db.sh $VENUS_SSH_HOST:$VENUS_TARGET_DIR/init_db.sh
scp -r ./nginx $VENUS_SSH_HOST:$VENUS_TARGET_DIR/
scp -r ./mysql $VENUS_SSH_HOST:$VENUS_TARGET_DIR/
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
