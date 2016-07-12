#!/usr/bin/env bash

if test ! $(which docker-machine)
  SERVER_ADDR=`docker run cimple-agent --server-addr $(docker-machine ip $DOCKER_MACHINE_NAME)`
else
  SERVER_ADDR="127.0.0.1"
fi

docker run cimple-agent --server-addr $SERVER_ADDR
