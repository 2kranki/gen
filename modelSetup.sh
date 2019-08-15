#!/bin/sh

# This script is used to set up the model directory which the genapp binary needs
# to generate programs.  You would execute this script if you wanted a copy of the
# model directory outside of the normal genapp directory.

echo "Setting up Model Directory with auxilliary data:"

if [[ -d "./src/models/sqlapp/util" ]]; then
    rm -fr ./src/models/sqlapp/util
fi
cp -R src/util ./src/models/sqlapp/

echo "...Model Directory is set up!"

