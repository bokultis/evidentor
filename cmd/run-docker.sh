IMAGE_NAME="evidentor"

docker stop ${IMAGE_NAME}
docker rm ${IMAGE_NAME}
docker rmi ${IMAGE_NAME}

docker run -it -p --name evidentor -p 3001:3001  ${IMAGE_NAME}:v0.0.1