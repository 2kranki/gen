#!/usr/bin/env bash

docker image rm -f mysql_new

docker image build -t mysql_new .


