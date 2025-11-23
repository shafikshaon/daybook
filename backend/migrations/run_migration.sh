#!/bin/bash

# Migration Runner Script for Daybook Backend
# This script runs SQL migrations on the PostgreSQL database

set -e  # Exit on error

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

echo -e "${GREEN}========================================${NC}"
echo -e "${GREEN}Daybook Database Migration Runner${NC}"
echo -e "${GREEN}========================================${NC}"

# Check if DATABASE_URL is set
if [ -z "$DATABASE_URL" ]; then
    echo -e "${YELLOW}DATABASE_URL not set. Using default...${NC}"
    DATABASE_URL="postgresql://postgres:postgres@localhost:5432/daybook_db"
fi

echo -e "\n${YELLOW}Database: ${DATABASE_URL}${NC}"

# Function to run a migration file
run_migration() {
    local migration_file=$1
    local migration_name=$(basename "$migration_file")

    echo -e "\n${YELLOW}Running migration: ${migration_name}${NC}"

    if psql "${DATABASE_URL}" -f "$migration_file"; then
        echo -e "${GREEN}✓ Migration ${migration_name} completed successfully${NC}"
        return 0
    else
        echo -e "${RED}✗ Migration ${migration_name} failed${NC}"
        return 1
    fi
}

# Get the directory where this script is located
SCRIPT_DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"

echo -e "\n${YELLOW}Migration files directory: ${SCRIPT_DIR}${NC}"

# Check if migration file exists
if [ ! -f "${SCRIPT_DIR}/001_convert_to_timestamptz.sql" ]; then
    echo -e "${RED}Migration file not found: 001_convert_to_timestamptz.sql${NC}"
    exit 1
fi

# Backup reminder
echo -e "\n${RED}WARNING: This migration will modify your database schema.${NC}"
echo -e "${YELLOW}It is HIGHLY recommended to backup your database first!${NC}"
echo -e "\n${YELLOW}To backup your database, run:${NC}"
echo -e "pg_dump ${DATABASE_URL} > daybook_backup_\$(date +%Y%m%d_%H%M%S).sql"
echo -e "\n"

read -p "Have you backed up your database? (yes/no): " -n 3 -r
echo
if [[ ! $REPLY =~ ^[Yy][Ee][Ss]$ ]]; then
    echo -e "${RED}Migration aborted. Please backup your database first.${NC}"
    exit 1
fi

# Run the migration
echo -e "\n${GREEN}Starting migration...${NC}"

if run_migration "${SCRIPT_DIR}/001_convert_to_timestamptz.sql"; then
    echo -e "\n${GREEN}========================================${NC}"
    echo -e "${GREEN}All migrations completed successfully!${NC}"
    echo -e "${GREEN}========================================${NC}"

    # Verification
    echo -e "\n${YELLOW}Verifying timestamp columns...${NC}"
    psql "${DATABASE_URL}" -c "
    SELECT table_name, column_name, data_type
    FROM information_schema.columns
    WHERE table_schema = 'public'
      AND (column_name LIKE '%date%' OR column_name LIKE '%time%' OR data_type LIKE '%timestamp%')
    ORDER BY table_name, column_name;
    "
else
    echo -e "\n${RED}========================================${NC}"
    echo -e "${RED}Migration failed!${NC}"
    echo -e "${RED}========================================${NC}"
    echo -e "\n${YELLOW}To rollback, you can restore from your backup:${NC}"
    echo -e "psql ${DATABASE_URL} < your_backup_file.sql"
    exit 1
fi

echo -e "\n${GREEN}Migration process completed.${NC}"
echo -e "${YELLOW}Please restart your backend server to apply changes.${NC}"
