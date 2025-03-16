#!/bin/bash

# TODO use ansible to install postgresql

# Configuration
DB_NAME="top_golang_arch"
DB_USER="postgres"
DB_PASSWORD="postgres"

echo "Starting PostgreSQL 16 installation and configuration..."

# Check if Homebrew is installed
if ! command -v brew &> /dev/null; then
    echo "Homebrew is not installed. Please install it first."
    echo "Run this command: /bin/bash -c \"\$(curl -fsSL https://raw.githubusercontent.com/Homebrew/install/HEAD/install.sh)\""
    exit 1
fi

# Check if PostgreSQL 16 is already installed
if brew list postgresql@16 &> /dev/null; then
    echo "PostgreSQL 16 is already installed."
else
    # Install PostgreSQL 16
    echo "Installing PostgreSQL 16..."
    brew install postgresql@16 || { echo "Failed to install PostgreSQL 16"; exit 1; }
fi

# Add PostgreSQL environment variables to .zshrc if not already present
if ! grep -q "PostgreSQL 16 environment variables" ~/.zshrc; then
    echo "Adding PostgreSQL environment variables..."
    cat << EOF >> ~/.zshrc

# PostgreSQL 16 environment variables
export PATH="/opt/homebrew/opt/postgresql@16/bin:\$PATH"
export LDFLAGS="-L/opt/homebrew/opt/postgresql@16/lib"
export CPPFLAGS="-I/opt/homebrew/opt/postgresql@16/include"
export PKG_CONFIG_PATH="/opt/homebrew/opt/postgresql@16/lib/pkgconfig"
export PGDATA="/opt/homebrew/var/postgresql@16"
export PGHOST="localhost"
export PGPORT="5432"
export PGUSER="$DB_USER"
export PGPASSWORD="$DB_PASSWORD"
EOF
fi

# Reload .zshrc if it exists
if [ -f ~/.zshrc ]; then
    echo "Reloading .zshrc..."
    source ~/.zshrc
fi

# Start PostgreSQL service
echo "Starting PostgreSQL service..."
brew services start postgresql@16 || { echo "Failed to start PostgreSQL service"; exit 1; }

# Wait for PostgreSQL to start
sleep 5

# Configure postgres user
echo "Configuring postgres user..."
/opt/homebrew/opt/postgresql@16/bin/psql postgres -c "ALTER USER postgres WITH LOGIN;" || { echo "Failed to configure postgres user"; exit 1; }
/opt/homebrew/opt/postgresql@16/bin/psql postgres -c "ALTER USER postgres WITH PASSWORD '$DB_PASSWORD';" || { echo "Failed to set password for postgres user"; exit 1; }
/opt/homebrew/opt/postgresql@16/bin/psql postgres -c "ALTER USER postgres CREATEDB;" || { echo "Failed to grant CREATEDB to postgres user"; exit 1; }

# Create database
echo "Creating database: $DB_NAME"
PGPASSWORD=$DB_PASSWORD /opt/homebrew/opt/postgresql@16/bin/createdb -U $DB_USER -h localhost $DB_NAME || { echo "Failed to create database $DB_NAME"; exit 1; }

echo "PostgreSQL 16 installation and configuration completed!"
echo "Database details:"
echo "  Database name: $DB_NAME"
echo "  User: $DB_USER"
echo "  Password: $DB_PASSWORD"
echo "  Port: 5432"
echo "  Host: localhost"

# Verify connection
echo "Verifying database connection..."
PGPASSWORD=$DB_PASSWORD /opt/homebrew/opt/postgresql@16/bin/psql -U $DB_USER -h localhost -d $DB_NAME -c "\conninfo" || { echo "Failed to verify database connection"; exit 1; }
