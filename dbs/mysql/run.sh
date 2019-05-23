#!/usr/bin/env bash

docker container rm -f mysql1

docker container run --name mysql1 -p 3300:3306 -v mysql_db1:/var/lib/mysql -d mysql_new 


