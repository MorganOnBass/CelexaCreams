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
TAGBASE="$DOCKER_HUB_USER/celexacreams:$TS"
SCRIPT_PATH=$(readlink -f $0)
PROJECT_ROOT=$(dirname $SCRIPT_PATH)/..
ARCHES=("amd64" "arm64v8")

docker manifest create $DOCKER_HUB_USER/celexacreams:latest
for ARCH in ${ARCHES[@]}
do
  TAG=$TAGBASE-$ARCH
  docker build -t $TAG $PROJECT_ROOT --build-arg ARCH=$ARCH
  docker push $TAG
  docker manifest create $DOCKER_HUB_USER/celexacreams:latest --amend $TAG
  docker manifest annotate $DOCKER_HUB_USER/celexacreams:latest $TAG --arch $ARCH
done
docker manifest push $DOCKER_HUB_USER/celexacreams:latest
