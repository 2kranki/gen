#!/bin/sh

MODELDIR="./models"

echo "Generating MariaDB Version of 01:"
if /tmp/bin/genapp -mdldir $MODELDIR -x misc/test01ma.exec.json.txt; then
    :
else
    echo "\tMariaDB Gen did not work!"
    exit 1
fi

echo "Generating MS SQL Version of 01:"
if /tmp/bin/genapp -mdldir $MODELDIR -x misc/test01ms.exec.json.txt; then
    :
else
    echo "\tMS SQL Gen did not work!"
    exit 1
fi

echo "Generating MySQL Version of 01:"
if /tmp/bin/genapp -mdldir $MODELDIR -x misc/test01my.exec.json.txt; then
    :
else
    echo "\tMySQL Gen did not work!"
    exit 1
fi

echo "Generating PostGres Version of 01:"
if /tmp/bin/genapp -mdldir $MODELDIR -x misc/test01pg.exec.json.txt; then
    :
else
    echo "\tPostGres Gen did not work!"
    exit 1
fi

echo "Generating SQLite Version of 01:"
if /tmp/bin/genapp -mdldir $MODELDIR -x misc/test01sq.exec.json.txt; then
    :
else
    echo "\tSQLite Gen did not work!"
    exit 1
fi

cd /tmp/app01
go fmt ./...
cd -

