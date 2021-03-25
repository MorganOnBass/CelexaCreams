#!/bin/bash
TS="$(date +%s)"
TAG="celexacreams:$TS"
SCRIPT_PATH=$(readlink -f $0)
PROJECT_ROOT=$(dirname $SCRIPT_PATH)/..

CGO_ENABLED=0 GOPATH=$PROJECT_ROOT/../../../../ /usr/local/go/bin/go build -o $PROJECT_ROOT/bin/celexacreams github.com/morganonbass/celexacreams/cmd/celexacreams
docker build -t $TAG $PROJECT_ROOT
docker tag $TAG celexacreams:latest
echo "$TAG"