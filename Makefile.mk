# project service name
NAME := courses

# local database connection settings
TEST_PGDATABASE ?= test-coursesdb

PGDATABASE ?= coursesdb
PGHOST ?= localhost
PGPORT ?= 5432
PGUSER ?= mikhail
PGPASSWORD ?= postgres

# add -race to GOFLAGS if RACE=1
RACE=0
