#!/usr/bin/env bash
# add -xv above after bash to debug

name="mssql1"
user="SA"
pw="Passw0rd!"
server="localhost"
port=1401
dockerName="mcr.microsoft.com/mssql/server"
dockerTag="2017-latest-ubuntu"

imageName="${dockerName}"
if [ -n "${dockerTag}" ]; then
    imageName="${dockerName}:${dockerTag}"
fi
#echo "Image Name: ${imageName}"

#echo "Ignore message: Error: No such container: mssql1"
s=`docker container rm -f ${name} 2>1`

if docker image ls ${imageName} | tail -n 1 | grep "${dockerName}"; then
    :
else
    docker image pull "${imageName}"
fi

containerID=`docker container run --name ${name} -e "ACCEPT_EULA=Y" -e "SA_PASSWORD=${pw}" -p ${port}:1433  -d "${imageName}"`
#echo "Container ID: ${containerID: -10}"

while ! `nc -z ${server} ${port}`; do sleep 3; done

echo ..."MSSQL Server, ${name}:${containerID: -10}, has started with user:${user} pw:${pw} on ${server}:${port}!"
