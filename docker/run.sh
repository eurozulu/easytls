#!/bin/bash
export TAG=$(git describe --abbrev=0 --tags)
export LABEL=eurospoofer/easytls
export CLAB=${LABEL}_container
echo Executing image $LABEL:$TAG as $CLAB
docker run -l $CLAB -p 443:443 $LABEL:$TAG
