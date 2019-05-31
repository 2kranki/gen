#!/usr/bin/env bash



docker container rm -f mysql1

docker container run --name mysql1 -e  MYSQL_ROOT_PASSWORD="Passw0rd!" -p 3306:3306  -d mysql:5.7

ipAddr=`docker inspect -f '{{range .NetworkSettings.Networks}}{{.IPAddress}}{{end}}' mysql1`
echo "mysql1 is  on tcp ip:${ipAddr}"


while ! `nc -z ${ipAddr} 3306`; do sleep 3; done
echo ..."MySQL Server has started!"
