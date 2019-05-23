#!/usr/bin/env sh

# Allow the SQL Server to start
sleep 30s

/usr/src/app/load.sh
/usr/src/app/teachLoad.sh
#/usr/src/app/tsqlLoad.sh

