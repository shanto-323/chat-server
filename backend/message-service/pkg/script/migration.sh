#!/usr/bin/env bash

set -e
cd "$(dirname "$0")/../migrations"
CASSANDRA_HOST=${CASSANDRA_HOST:-scylladb}
CQL_FILES=(
  "001_keyspace.cql"
  "002_chat_table.cql"
)

until cqlsh $CASSANDRA_HOST -e "DESC KEYSPACES";do
  sleep 2
done

for cql in "${CQL_FILES[@]}";do
  cqlsh $CASSANDRA_HOST -f "$cql"
done

echo "Migration have been applied."
  