#!/bin/bash

# Start the first process
./cloud_sql_proxy -instances=${1}=tcp:3306 -credential_file=${2} &

# Start the second process
./main -db-password=${3} -log-level=-1 -log-time-format=2006-01-02T15:04:05.999999999Z07:00

