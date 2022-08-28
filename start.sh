#!/bin/sh


set -e

echo "run db migration"
source /app/app.env
/app/migrate -path /app/migrations -database "$DB_SOURCE" -verbose up

echo "start the app"
exec "$@"

