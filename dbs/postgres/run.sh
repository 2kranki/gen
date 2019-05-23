#!/usr/bin/env bash

docker container rm -f postgres1

docker container run --name postgres1 -e -p 5432:5432 -d postgres_new


