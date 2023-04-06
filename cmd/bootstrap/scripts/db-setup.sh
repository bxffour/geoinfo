#!/bin/bash

set -e

migrate -path=/migrations -database=$POSTGRES_DB_DSN up

psql -v ON_ERROR_STOP=1 $POSTGRES_DB_DSN <<-EOSQL
    DO \$\$
    BEGIN
        IF NOT EXISTS (SELECT 1 FROM pg_roles WHERE rolname= 'geoinfo') THEN
            CREATE ROLE geoinfo WITH LOGIN PASSWORD '$GEOINFO_API_PASSWORD';
            GRANT SELECT ON countries TO geoinfo;
        ELSE
            RAISE NOTICE 'Role geoinfo already exists, skipping...';
        END IF;
    END 
    \$\$
    ;        

EOSQL

bootstrap -d "$POSTGRES_DB_DSN"