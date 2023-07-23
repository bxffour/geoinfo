#!/bin/bash

set -e

if [ -z "${POSTGRES_DB_DSN}" ] || [ -z "${GEOINFO_API_USER}" ] || [ -z "${GEOINFO_API_PASSWORD}" ]; then
    echo "the POSTGRES_DB_DSN, GEOINFO_API_USER and GEOINFO_API_PASSWORD envars must be set"
    exit 15
fi

migrate -path=/migrations -database=$POSTGRES_DB_DSN up

psql -v ON_ERROR_STOP=1 $POSTGRES_DB_DSN <<-EOSQL
    DO \$\$
    BEGIN
        IF NOT EXISTS (SELECT 1 FROM pg_roles WHERE rolname= '$GEOINFO_API_USER') THEN
            CREATE ROLE $GEOINFO_API_USER WITH LOGIN PASSWORD '$GEOINFO_API_PASSWORD';
            GRANT SELECT ON countries TO $GEOINFO_API_USER;
        ELSE
            RAISE NOTICE 'Role $GEOINFO_API_USER already exists, skipping...';
        END IF;
    END 
    \$\$
    ;        

EOSQL

bootstrap -d "$POSTGRES_DB_DSN"