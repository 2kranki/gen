#!/usr/bin/env bash
# add -xv above after bash to debug

name="mssql1"
user="sa"
pw="Passw0rd"
server="localhost"
port=1401
dockerName="mcr.microsoft.com/mssql/server"
dockerTag="2017-latest-ubuntu"

imageName="${dockerName}"
if [ -n "${dockerTag}" ]; then
    imageName="${dockerName}:${dockerTag}"
fi
echo "Image Name: ${imageName}"

echo "Deleting Container: ${name}..."
echo "...Ignore message: Error: No such container: mssql1"
s=`docker container rm -f ${name} 2>&1`

echo "Pulling Image: ${imageName} if needed..."
if docker image ls ${imageName} | tail -n 1 | grep "${dockerName}"; then
    echo "...Image: ${imageName} present."
else
    echo "...Pulling Image: ${imageName}:"
    docker image pull "${imageName}"
fi

echo "Running Container: ${name}..."
containerID=`docker container run --name ${name} -e "ACCEPT_EULA=Y" -e "MSSQL_SA_PASSWORD=${pw}" -p ${port}:1433  -d "${imageName}"`
echo "...Container ID: ${containerID: -10}"

echo "Waiting for Container: ${name} to initialize..."
while ! `nc -z ${server} ${port}`; do sleep 3; done

echo ..."MSSQL Server, ${name}:${containerID: -10}, has started with user:${user} pw:${pw} on ${server}:${port}!"
