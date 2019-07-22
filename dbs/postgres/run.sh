#!/usr/bin/env bash

name="postgres1"
user="root"
pw="Passw0rd!"
server="localhost"
port=5430

echo "Ignore message: Error: No such container: mssql1"
docker container rm -f ${name}

containerID=`docker container run --name ${name} -e "POSTGRES_PASSWORD=${pw}" -p ${port}:3306  -d postgres`

#echo "Container ID: ${containerID: -10}"

while ! `nc -z ${server} ${port}`; do sleep 3; done

echo ..."Postgres SQL Server, ${name}:${containerID: -10}, has started with user:${user} pw:${pw} on ${server}:${port}!"

