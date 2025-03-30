set dotenv-load
set dotenv-required
set dotenv-filename := ".env.development"

# Lists all available recipes.
_default:
  @just --list

# Imports database schema.
import:
        #!/usr/bin/env bash
        echo "Importing PostgreSQL database..."
        
        # Wait for database container to start before running pg_restore.
        until PGPASSWORD=$POSTGRES_PASSWORD psql -h $POSTGRES_HOST -U $POSTGRES_USER -d $POSTGRES_DB -c "SELECT 1" &> /dev/null; do
                sleep 2
        done

        PGPASSWORD=$POSTGRES_PASSWORD psql -h $POSTGRES_HOST -U $POSTGRES_USER -d $POSTGRES_DB -q -f init.sql && \
        echo "PostgreSQL database is ready!"
