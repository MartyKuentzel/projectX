#!/bin/bash

PW=$(cat ${CLOUDSQL_ROOT_CREDENTIALS})
# Start the first process
./cloud_sql_proxy -instances=${INSTANCE}=tcp:3306 -credential_file=${GOOGLE_APPLICATION_CREDENTIALS} &

# Start the second process
./main -db-password=${PW} -log-level=-1 -log-time-format=2006-01-02T15:04:05.999999999Z07:00

