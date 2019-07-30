#!/bin/sh

echo "Generating MariaDB Version of 01:"
/tmp/bin/genapp -mdldir ./src/models -x misc/test01ma.exec.json.txt

echo "Generating MS SQL Version of 01:"
/tmp/bin/genapp -mdldir ./src/models -x misc/test01ms.exec.json.txt

echo "Generating MySQL Version of 01:"
/tmp/bin/genapp -mdldir ./src/models -x misc/test01my.exec.json.txt

echo "Generating PostGres Version of 01:"
/tmp/bin/genapp -mdldir ./src/models -x misc/test01pg.exec.json.txt

echo "Generating SQLite Version of 01:"
/tmp/bin/genapp -mdldir ./src/models -x misc/test01sq.exec.json.txt


