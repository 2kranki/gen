#!/bin/bash

ipAddr=`docker inspect -f '{{range .NetworkSettings.Networks}}{{.IPAddress}}{{end}}' $@`
echo "Starting $@ app on tcp ip:${ipAddr}"

