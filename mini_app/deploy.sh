#!/bin/bash

set -e

docker compose -f docker/docker-compose.yml down

docker-compose -f docker/docker-compose.yml up --build -d