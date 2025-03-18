#!/bin/sh

sleep 5 # wait for mysql db started
cat /source/data.sql | mysql -h mysql -P 3306 -u root -ppass mysql
