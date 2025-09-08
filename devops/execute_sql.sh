#!/bin/bash

# SQL Execution Tracking Script
# This script executes SQL files in sorted order and tracks which ones have been executed

set -e  # Exit on any error

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
SQL_DIR="$SCRIPT_DIR"
EXECUTED_FILE="$SCRIPT_DIR/.executed_sql_files.txt"
DB_FILE="$SCRIPT_DIR/../contatos.db"

# Create executed files tracking file if it doesn't exist
if [ ! -f "$EXECUTED_FILE" ]; then
    touch "$EXECUTED_FILE"
    echo "Created tracking file: $EXECUTED_FILE"
fi

# Function to check if a file has been executed
is_executed() {
    local filename="$1"
    grep -q "^$filename$" "$EXECUTED_FILE" 2>/dev/null
}

# Function to mark file as executed
mark_executed() {
    local filename="$1"
    echo "$filename" >> "$EXECUTED_FILE"
    echo "Marked $filename as executed"
}

# Get all SQL files sorted by name
echo "Looking for SQL files in: $SQL_DIR"
sql_files=($(find "$SQL_DIR" -name "*.sql" -type f | sort))

if [ ${#sql_files[@]} -eq 0 ]; then
    echo "No SQL files found in $SQL_DIR"
    exit 0
fi

echo "Found ${#sql_files[@]} SQL file(s):"
for file in "${sql_files[@]}"; do
    basename "$file"
done
echo

# Check if sqlite3 is available
if ! command -v sqlite3 &> /dev/null; then
    echo "Error: sqlite3 is not installed or not in PATH"
    exit 1
fi

# Execute SQL files that haven't been executed yet
executed_count=0
skipped_count=0

for sql_file in "${sql_files[@]}"; do
    filename=$(basename "$sql_file")

    if is_executed "$filename"; then
        echo "SKIPPED: $filename (already executed)"
        ((skipped_count++))
    else
        echo "EXECUTING: $filename"

        # Show file contents for review
        echo "--- Content of $filename ---"
        cat "$sql_file"
        echo "--- End of $filename ---"
        echo

        # Execute the SQL file
        if sqlite3 "$DB_FILE" < "$sql_file"; then
            mark_executed "$filename"
            echo "SUCCESS: $filename executed successfully"
            ((executed_count++))
        else
            echo "ERROR: Failed to execute $filename"
            exit 1
        fi
        echo
    fi
done

echo "=== SUMMARY ==="
echo "Total SQL files found: ${#sql_files[@]}"
echo "Files executed: $executed_count"
echo "Files skipped: $skipped_count"
echo "Database file: $DB_FILE"
echo "Tracking file: $EXECUTED_FILE"

if [ $executed_count -gt 0 ]; then
    echo
    echo "Database tables created/updated successfully!"

    # Show current database schema
    echo
    echo "=== CURRENT DATABASE SCHEMA ==="
    sqlite3 "$DB_FILE" ".schema"
fi
