#!/bin/bash

function exec_psql() {
    PGPASSWORD=$POSTGRES_PASSWORD psql -h $POSTGRES_HOST -U $POSTGRES_USER -d $POSTGRES_DB -q
}

# Set dev container directory as safe to fix dubious ownership warning.
if ! git config --global --get-all safe.directory | grep -Fxq "$0"; then
    git config --global --add safe.directory "$0"
fi

# Get number of public tables, and extract the count value from the 3rd line of the output. 
TABLE_COUNT=$(echo "SELECT COUNT(*) FROM information_schema.tables t WHERE t.table_schema = 'public';" | exec_psql | sed -n 3p | xargs)

# Import database schema if no tables are found.
if [ $TABLE_COUNT -eq 0 ]; then
    echo "Importing PostgreSQL database..."

    # Wait for database container to start before creating tables.
    until echo "SELECT 1" | exec_psql | &> /dev/null; do
            sleep 2
    done

    cat ./.devcontainer/init.sql | exec_psql &> /dev/null && echo "PostgreSQL database is ready!"
fi