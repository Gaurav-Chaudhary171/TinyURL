#!/bin/bash

# Exit on any error
set -e

# Function to check if a command exists
command_exists() {
    command -v "$1" >/dev/null 2>&1
}

# Check if MySQL is installed
if ! command_exists mysql; then
    echo "❌ MySQL is not installed. Please install MySQL first."
    exit 1
fi

# Stop any running instance
echo "Stopping any running MySQL instances..."
brew services stop mysql || true

# Create Data directories
MYSQL_DIR="$HOME/mysql_tinyurl"
CONFIG_FILE="$HOME/mysql_tinyurl.cnf"
SOCKET_FILE="/tmp/mysql_tinyurl.sock"

# Clean up existing directory and socket
echo "Cleaning up existing MySQL directory and socket..."
rm -rf "$MYSQL_DIR"
rm -f "$SOCKET_FILE"

# Create mysql Config
echo "Creating MySQL configuration..."
cat > "$CONFIG_FILE" << EOL
[mysqld]
port=3306
datadir=$MYSQL_DIR
socket=$SOCKET_FILE
pid-file=$MYSQL_DIR/mysql.pid
log-error=$MYSQL_DIR/error.log
user=$(whoami)
mysqlx=OFF
mysqlx_socket=/tmp/mysqlx_tinyurl.sock
EOL

# Create directory if it doesn't exist
echo "Creating MySQL data directory..."
mkdir -p "$MYSQL_DIR"

# Initialize MySQL data directories
echo "Initializing MySQL data directory..."
if ! mysqld --initialize-insecure --datadir="$MYSQL_DIR" 2>&1; then
    echo "❌ Failed to initialize MySQL data directory"
    exit 1
fi

# Start MySQL instance in background
echo "Starting MySQL server..."
mysqld --defaults-file="$CONFIG_FILE" &
MYSQL_PID=$!

# Wait for MySQL to start
echo "Waiting for MySQL to start..."
for i in {1..30}; do
    if [ -S "$SOCKET_FILE" ] && mysqladmin ping -S "$SOCKET_FILE" --silent 2>/dev/null; then
        echo "MySQL started successfully!"
        break
    fi
    if [ $i -eq 30 ]; then
        echo "❌ MySQL failed to start after 30 seconds"
        echo "Checking error log:"
        cat "$MYSQL_DIR/error.log"
        kill $MYSQL_PID
        exit 1
    fi
    echo "Waiting for MySQL to start... ($i/30)"
    sleep 1
done

# Create databases and users
echo "Creating database and user..."
mysql -S "$SOCKET_FILE" -u root -e "
CREATE DATABASE IF NOT EXISTS tinyurl;
CREATE USER IF NOT EXISTS 'tinyurl'@'localhost' IDENTIFIED BY 'tinyurl';
GRANT ALL PRIVILEGES ON tinyurl.* TO 'tinyurl'@'localhost';
FLUSH PRIVILEGES;"

if [ $? -eq 0 ]; then
    echo "✅ Setup completed successfully"
    echo "To connect to the database, use either:"
    echo "1. Via TCP:"
    echo "   mysql -u tinyurl -h localhost -P 3306 -p"
    echo "2. Via Socket:"
    echo "   mysql -u tinyurl -S $SOCKET_FILE -p"
    echo "Password: tinyurl"
else
    echo "❌ Database setup failed"
    kill $MYSQL_PID
    exit 1
fi



