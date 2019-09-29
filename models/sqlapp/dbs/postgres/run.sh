#!/usr/bin/env bash

name="postgres1"
user="postgres"
pw="Passw0rd"
server="localhost"
port=5432
dockerName="postgres"
dockerTag="latest"

imageName="${dockerName}"
if [ -n "${dockerTag}" ]; then
    imageName="${dockerName}:${dockerTag}"
fi
echo "Image Name: ${imageName}"

echo "Deleting Container: ${name}..."
echo "...Ignore message: Error: No such container: ${name}"
docker container rm -f ${name}

echo "Pulling Image: ${imageName} if needed..."
if docker image ls ${imageName} | tail -n 1 | grep "${dockerName}"; then
    echo "...Image: ${imageName} present."
else
    echo "...Pulling Image: ${imageName}:"
    docker image pull "${imageName}"
fi

echo "Running Container: ${name}..."
containerID=`docker container run --name ${name} -e "POSTGRES_PASSWORD=${pw}" -p ${port}:5432  -d postgres`
echo "...Container ID: ${containerID: -10}"

echo "Waiting for Container: ${name} to initialize..."
while ! `nc -z ${server} ${port}`; do sleep 3; done

echo ..."Postgres SQL Server, ${name}:${containerID: -10}, has started with user:${user} pw:${pw} on ${server}:${port}!"

