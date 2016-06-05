#!/usr/bin/env bash

DB=$1
OUT="app/models/db_struct.go"
rm $OUT
echo "package models" >"$OUT"
psql -t -A -d "$DB" -U $2 -h localhost -f extra/pg2go.sql >>"$OUT"
goimports -w "$OUT" || gofmt -w "$OUT"