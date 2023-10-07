#!/bin/bash

# Stop all running containers
if [ "$(docker ps -q)" ]; then
    docker stop $(docker ps -q)
fi

# Remove all containers
if [ "$(docker ps -aq)" ]; then
    docker rm $(docker ps -aq)
fi

# Remove all images
if [ "$(docker images -q)" ]; then
    docker rmi $(docker images -q)
fi

# Remove all volumes
if [ "$(docker volume ls -q)" ]; then
    docker volume rm $(docker volume ls -q)
fi

echo "All Docker containers, images and volumes have been removed."
