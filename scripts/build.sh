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
DOCKER_HUB_USER="morganonbass"
TAG="$DOCKER_HUB_USER/celexacreams:$TS"
SCRIPT_PATH=$(readlink -f $0)
PROJECT_ROOT=$(dirname $SCRIPT_PATH)/..

docker build -t $TAG $PROJECT_ROOT
docker tag $TAG $DOCKER_HUB_USER/celexacreams:latest
echo "$TAG"