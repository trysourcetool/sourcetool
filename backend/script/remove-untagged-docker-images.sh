#!/bin/sh

set -e

processes=`docker ps -q -f status=exited`
if [ -n "$processes" ]; then
  docker rm $processes
fi

images=`docker images -q -f dangling=true`
if [ -n "$images" ]; then
  docker rmi -f $images
fi