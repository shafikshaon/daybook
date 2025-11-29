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
echo -e "\n${YELLOW}Usage: ./run_migration.sh [migration_number|all]${NC}"
echo -e "${YELLOW}Examples:${NC}"
echo -e "  ./run_migration.sh          # Run all migrations"
echo -e "  ./run_migration.sh all      # Run all migrations"
echo -e "  ./run_migration.sh 001      # Run migration 001"
echo -e "  ./run_migration.sh 000001   # Run migration 000001"

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

# Migration mode: 'all' or specific migration number
MIGRATION_MODE="${1:-all}"

# List of migrations in order
MIGRATIONS=(
    "000001_convert_transaction_datetime_to_date.up.sql"
)

# Function to check if migration file exists
check_migration_exists() {
    local migration_file=$1
    if [ ! -f "${SCRIPT_DIR}/${migration_file}" ]; then
        echo -e "${YELLOW}Warning: Migration file not found: ${migration_file}${NC}"
        return 1
    fi
    return 0
}

# Determine which migrations to run
MIGRATIONS_TO_RUN=()

if [ "$MIGRATION_MODE" = "all" ]; then
    echo -e "\n${YELLOW}Running all available migrations...${NC}"
    for migration in "${MIGRATIONS[@]}"; do
        if check_migration_exists "$migration"; then
            MIGRATIONS_TO_RUN+=("$migration")
        fi
    done
else
    # Run specific migration by number or name
    MIGRATION_FOUND=false
    for migration in "${MIGRATIONS[@]}"; do
        if [[ "$migration" == *"$MIGRATION_MODE"* ]]; then
            if check_migration_exists "$migration"; then
                MIGRATIONS_TO_RUN+=("$migration")
                MIGRATION_FOUND=true
            fi
        fi
    done

    if [ "$MIGRATION_FOUND" = false ]; then
        echo -e "${RED}Migration not found: ${MIGRATION_MODE}${NC}"
        echo -e "${YELLOW}Available migrations:${NC}"
        for migration in "${MIGRATIONS[@]}"; do
            echo -e "  - ${migration}"
        done
        exit 1
    fi
fi

if [ ${#MIGRATIONS_TO_RUN[@]} -eq 0 ]; then
    echo -e "${RED}No migrations to run!${NC}"
    exit 1
fi

echo -e "\n${YELLOW}Migrations to run:${NC}"
for migration in "${MIGRATIONS_TO_RUN[@]}"; do
    echo -e "  - ${migration}"
done

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

# Run the migrations
echo -e "\n${GREEN}Starting migrations...${NC}"

FAILED=false
for migration in "${MIGRATIONS_TO_RUN[@]}"; do
    if ! run_migration "${SCRIPT_DIR}/${migration}"; then
        FAILED=true
        break
    fi
done

if [ "$FAILED" = false ]; then
    echo -e "\n${GREEN}========================================${NC}"
    echo -e "${GREEN}All migrations completed successfully!${NC}"
    echo -e "${GREEN}========================================${NC}"

    # Verification
    echo -e "\n${YELLOW}Verifying timestamp and date columns...${NC}"
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
