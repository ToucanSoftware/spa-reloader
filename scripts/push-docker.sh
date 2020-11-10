#!/bin/bash

set -ex

docker tag docker.io/toucansoftware/spa-reloader:0.0.${TRAVIS_BUILD_NUMBER} docker.io/toucansoftware/spa-reloader:latest
docker push docker.io/toucansoftware/spa-reloader:0.0.${TRAVIS_BUILD_NUMBER}
docker push docker.io/toucansoftware/spa-reloader:latest
