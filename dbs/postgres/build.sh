#!/usr/bin/env bash

docker image rm -f postgres_new

docker image build -t postgres_new .


