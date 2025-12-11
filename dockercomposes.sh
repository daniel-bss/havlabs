#!/bin/bash

if [[ "$1" != "up" && "$1" != "upbuild" && "$1" != "down" ]]; then
    echo "Unknown argument provided: $1"
    exit 1
fi

yamls=(
    "./internal/auth/docker-compose.yaml"
    "./internal/news/docker-compose.yaml"
    "./internal/z/docker-compose.yaml"
)

for yaml in "${yamls[@]}"; do
    if [ "$1" == "up" ]; then
        sudo docker compose -f "$yaml" up -d
    elif [ "$1" == "upbuild" ]; then
        sudo docker compose -f "$yaml" up --build -d
    elif [ "$1" == "down" ]; then
        sudo docker compose -f "$yaml" down
    else
        echo "Unknown argument provided: $1"
    fi
done
