#!/bin/bash
set -ue

make

export DB_DATABASE=torb
export DB_HOST=127.0.0.1
export DB_PORT=43306
export DB_USER=isucon
export DB_PASS=isucon

cat /dev/null > ./echo.log
./torb
