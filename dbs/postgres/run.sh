#!/usr/bin/env bash

docker container rm -f postgres1

docker container run --name postgres1 -e POSTGRES_PASSWORD='Passw0rd!' -p 5430:5432 -d postgres


