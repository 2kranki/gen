#!/usr/bin/env bash

docker container rm -f mariadb1

docker container run --name mariadb1 -e MYSQL_ROOT_PASSWORD='Passw0rd!' -e MYSQL_DATABASE='Finances' -p 4306:3306  -d mariadb


