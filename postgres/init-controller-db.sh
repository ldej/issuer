#!/bin/bash
set -e

psql -v ON_ERROR_STOP=1 --username "$POSTGRES_USER" --dbname "$POSTGRES_DB" <<-EOSQL
    CREATE USER controller;
    CREATE DATABASE controller;
    GRANT ALL PRIVILEGES ON DATABASE controller TO controller;
EOSQL