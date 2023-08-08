#!/bin/bash

docker compose -f test/docker-compose.yml --env-file .env up mariadb -d

docker compose -f test/docker-compose.yml --env-file .env run migrate

res=$(docker compose -f test/docker-compose.yml --env-file .env run test_app)

echo -e $res

docker compose -f test/docker-compose.yml --env-file .env down

if [[ $res == *"Test Suite Passed"* ]]; then
  echo "Test Suite Passed"  
  exit 0
else
  echo "Test Suite Failed"
  exit 1
fi