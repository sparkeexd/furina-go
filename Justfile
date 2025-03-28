set dotenv-load
set dotenv-required
set dotenv-filename := ".env.development"

# Lists all available recipes.
_default:
  @just --list

# Starts mimo stack.
start:
        #!/usr/bin/env bash
        echo "Starting mimo..."
        just import
        COMPOSE_BAKE=true docker compose up --watch --force-recreate bot database

# Stops mimo stack.
stop:
        #!/usr/bin/env bash
        echo "Stopping mimo..."
        docker compose stop bot database 

# Imports database.
import:
        #!/usr/bin/env bash
        if [ -z "$(docker volume ls -q -f name=db)" ]; then
                COMPOSE_BAKE=true docker compose up -d --force-recreate database

                echo "Importing PostgreSQL database..."
                
                # Wait for database container to start before running pg_restore.
                until PGPASSWORD=$POSTGRES_PASSWORD psql -h $POSTGRES_HOST -U $POSTGRES_USER -d $POSTGRES_DB -c "SELECT 1" &> /dev/null; do
                        sleep 2
                done

                PGPASSWORD=$POSTGRES_PASSWORD psql -h $POSTGRES_HOST -U $POSTGRES_USER -d $POSTGRES_DB -q -f init.sql && \
                echo "PostgreSQL database is ready!"
        fi
