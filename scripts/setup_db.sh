#!/bin/bash

# Load environment variables from .env file
source .env

# Create database and user with proper permissions
echo "Creating database $DB_NAME and user $DB_USER..."

# Run SQL script as r user (with superuser privileges)
echo "Running SQL script with superuser privileges..."
psql -U r -d postgres -c "DO \$\$ 
BEGIN
    IF NOT EXISTS (SELECT FROM pg_catalog.pg_roles WHERE rolname = '$DB_USER') THEN
        CREATE USER \"$DB_USER\" WITH PASSWORD '$DB_PASSWORD';
    END IF;
END
\$\$;"

# Create database if not exists
psql -U r -d postgres -c "SELECT 'CREATE DATABASE \"$DB_NAME\"' WHERE NOT EXISTS (SELECT FROM pg_database WHERE datname = '$DB_NAME')\gexec"

# Grant privileges
psql -U r -d postgres -c "GRANT ALL PRIVILEGES ON DATABASE \"$DB_NAME\" TO \"$DB_USER\";"
psql -U r -d "$DB_NAME" -c "GRANT ALL PRIVILEGES ON SCHEMA public TO \"$DB_USER\";"
psql -U r -d "$DB_NAME" -c "ALTER USER \"$DB_USER\" WITH SUPERUSER;"

echo "Database setup completed successfully!" 
