#!/usr/bin/env bash

name="mssql1"
user="SA"
pw="Passw0rd!"
server="localhost"
port=1433

echo "Ignore message: Error: No such container: mssql1"
docker container rm -f ${name}

containerID=`docker container run --name ${name} -e "ACCEPT_EULA=Y" -e "SA_PASSWORD=${pw}" -p ${port}:1433  -d mcr.microsoft.com/mssql/server:2017-latest`
#echo "Container ID: ${containerID: -10}"

while ! `nc -z ${server} ${port}`; do sleep 3; done

echo ..."MSSQL Server, ${name}:${containerID: -10}, has started with user:${user} pw:${pw} on ${server}:${port}!"
