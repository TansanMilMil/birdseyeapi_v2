#!/bin/bash -eu

echo "load .env ----------------------------"
if [ ! -f .env ]; then
    echo ".env file not found!"
else 
    source .env
fi

ssh $VENUS_SSH_HOST curl -wv localhost:1111/HealthCheck
echo "----------------------------------------"
ssh $VENUS_SSH_HOST curl -wv localhost:1111/news/today-news | jq .
