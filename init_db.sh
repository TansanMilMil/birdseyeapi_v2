#!/bin/bash -eu

setup() {
    cd `dirname $0`
    if [ -f .env ]; then
        source .env
    else
        echo ".env file not found."
    fi
    docker compose up -d
}

setup

MYSQL="docker compose exec -T mysql mysql -u root -p$MYSQL_ROOT_PASSWORD"
$MYSQL < mysql/create_db.sql
#$MYSQL show databases

echo 'init db completed!'

