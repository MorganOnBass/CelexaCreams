#!/bin/bash
shopt -s expand_aliases
[[ `uname` == 'Darwin' ]] && {
	which greadlink gsed gzcat > /dev/null && {
		unalias readlink sed zcat
		alias readlink=greadlink sed=gsed zcat=gzcat
	} || {
		echo 'ERROR: GNU utils required for Mac. You may use homebrew to install them: brew install coreutils gnu-sed'
		exit 1
	}
}

TS="$(date +%s)"
TAG="celexacreams:$TS"
SCRIPT_PATH=$(readlink -f $0)
PROJECT_ROOT=$(dirname $SCRIPT_PATH)/..

CGO_ENABLED=0 GOOS=linux GOARCH=amd64 GOPATH=$PROJECT_ROOT/go $(which go) build -o $PROJECT_ROOT/bin/celexacreams $PROJECT_ROOT
docker build -t $TAG $PROJECT_ROOT
docker tag $TAG celexacreams:latest
echo "$TAG"