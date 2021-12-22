#!/bin/bash
export TAG=$(git describe --abbrev=0 --tags)
export LABEL=eurospoofer/easytls
echo Building $LABEL:$TAG
docker build -t $LABEL:$TAG -f docker/Dockerfile .
echo Pushing $LABEL:$TAG
docker push $LABEL:$TAG