#!/usr/bin/env bash

name="mysql1"
user="root"
pw="Passw0rd!"
server="localhost"
port=3306

echo "Ignore message: Error: No such container: mssql1"
docker container rm -f ${name}

containerID=`docker container run --name ${name} -e "MYSQL_ROOT_PASSWORD=${pw}" -e "MYSQL_DATABASE='Finances'" -p ${port}:3306  -d mysql:5.7`

#echo "Container ID: ${containerID: -10}"

while ! `nc -z ${server} ${port}`; do sleep 3; done

echo ..."MySQL Server, ${name}:${containerID: -10}, has started with user:${user} pw:${pw} on ${server}:${port}!"

