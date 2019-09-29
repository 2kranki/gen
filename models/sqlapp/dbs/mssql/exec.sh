#!/usr/bin/env bash

name="mssql1"
user="sa"
pw="Passw0rd"

echo "Remember: /opt/mssql-tools/bin/sqlcmd -U ${user} -P ${pw}"
docker container exec -it  ${name} bash


