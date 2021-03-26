#!/bin/bash
set -x

if [ -z $DISCORD_TOKEN ]; then
    echo "discord token not set"
    exit 1
fi

if [ -z $GIPHY_API_KEY ]; then
    echo "giphy api key not set"
    exit 1
fi


docker stop celexacreams && docker rm celexacreams
docker run -d --env DISCORD_TOKEN=$DISCORD_TOKEN --env GIPHY_API_KEY=$GIPHY_API_KEY --name celexacreams morganonbass/celexacreams:latest