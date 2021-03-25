#!/bin/bash
set -x

if [ -z $DISCORD_TOKEN ]; then
    echo "discord token not set"
    exit 1
fi

if [ -z $GIPHY_API_KEY ]; then
    echo "discord token not set"
    exit 1
fi


sudo docker stop celexacreams && sudo docker rm celexacreams
sudo docker run -d --env DISCORD_TOKEN=$DISCORD_TOKEN --env GIPHY_API_KEY=$GIPHY_API_KEY --name celexacreams celexacreams:latest