#!/bin/sh

echo "Setting up Model Directory with auxilliary data:"

if [[ -d "./src/models/sqlapp/util" ]]; then
    rm -fr ./src/models/sqlapp/util
fi
cp -R src/util ./src/models/sqlapp/

if [[ -d "./src/models/sqlapp/docker" ]]; then
    rm -fr ./src/models/sqlapp/docker
fi
cp -R src/docker ./src/models/sqlapp/

echo "...Model Directory is set up!"

