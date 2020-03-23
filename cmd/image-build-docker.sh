#!/bin/bash

set -ex

# PARENT_DIR=$(basename "${PWD%/*}")
# CURRENT_DIR="${PWD##*/}"
IMAGE_NAME="evidentor"
TAG="${1}"

docker build -t ${IMAGE_NAME}:${TAG} ../cmd
docker tag ${IMAGE_NAME}:${TAG} ${IMAGE_NAME}:latest
