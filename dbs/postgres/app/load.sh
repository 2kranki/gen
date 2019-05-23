#!/usr/bin/env sh

/opt/mssql-tools/bin/sqlcmd -S localhost -U SA -P 'Passw0rd!' -i /usr/src/app/load.sql.txt 
