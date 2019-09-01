#!/bin/bash -x

ROOT_DIR=$(cd $(dirname $0)/..; pwd)
DB_DIR="$ROOT_DIR/db"
BENCH_DIR="$ROOT_DIR/bench"
MYSQL_HOST=127.0.0.1
MYSQL_PORT=43306

export MYSQL_PWD=isucon

mysql -uisucon -h $MYSQL_HOST --port $MYSQL_PORT -e "DROP DATABASE IF EXISTS torb; CREATE DATABASE torb;"
mysql -uisucon -h $MYSQL_HOST --port $MYSQL_PORT torb < "$DB_DIR/schema.sql"

if [ ! -f "$DB_DIR/isucon8q-initial-dataset.sql.gz" ]; then
  echo "Run the following command beforehand." 1>&2
  echo "$ ( cd \"$BENCH_DIR\" && bin/gen-initial-dataset )" 1>&2
  exit 1
fi

mysql -uisucon -h $MYSQL_HOST --port $MYSQL_PORT torb -e 'ALTER TABLE reservations DROP KEY event_id_and_sheet_id_idx'
gzip -dc "$DB_DIR/isucon8q-initial-dataset.sql.gz" | mysql -uisucon -h $MYSQL_HOST --port $MYSQL_PORT torb
mysql -uisucon -h $MYSQL_HOST --port $MYSQL_PORT torb -e 'ALTER TABLE reservations ADD KEY event_id_and_sheet_id_idx (event_id, sheet_id)'
