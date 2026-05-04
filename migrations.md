I'll set up a proper migration system that works with sqlc. Let me first check the current sqlc configuration and then create a migration structure.

I'll set up a proper migration workflow using golang-migrate format with sqlc. This involves creating migration files, updating sqlc.yaml to use migrations as the schema source, and updating docker-compose to run migrations on startup.



Now I'll add a migration service to docker-compose and update the Makefile with migration commands:





Let me verify the final configuration:

Done! Here's the complete migration setup with sqlc:

## Changes Made

### 1. Created Migration Files
- [migrations/000001_init_schema.up.sql](cci:7://file:///Users/kennethsimpson/go/src/github.com/grpc_sqlc/migrations/000001_init_schema.up.sql:0:0-0:0) - Creates the `users` table
- [migrations/000001_init_schema.down.sql](cci:7://file:///Users/kennethsimpson/go/src/github.com/grpc_sqlc/migrations/000001_init_schema.down.sql:0:0-0:0) - Drops the table for rollback

### 2. Updated [sqlc.yaml](cci:7://file:///Users/kennethsimpson/go/src/github.com/grpc_sqlc/sqlc.yaml:0:0-0:0)
Changed `schema` from [schema.sql](cci:7://file:///Users/kennethsimpson/go/src/github.com/grpc_sqlc/schema.sql:0:0-0:0) to [migrations/](cci:9://file:///Users/kennethsimpson/go/src/github.com/grpc_sqlc/migrations:0:0-0:0) - sqlc now reads all migration files to compute the database schema for code generation.

### 3. Updated [docker-compose.yml](cci:7://file:///Users/kennethsimpson/go/src/github.com/grpc_sqlc/docker-compose.yml:0:0-0:0)
- **Removed** the direct [schema.sql](cci:7://file:///Users/kennethsimpson/go/src/github.com/grpc_sqlc/schema.sql:0:0-0:0) mount to postgres (was only for init)
- **Added** a `migrate` service using `migrate/migrate:v4.17.0` that runs migrations automatically
- **Added** healthcheck to postgres so migrations wait for DB readiness
- **Updated** server `depends_on` to wait for migrations to complete

### 4. Updated [Makefile](cci:7://file:///Users/kennethsimpson/go/src/github.com/grpc_sqlc/Makefile:0:0-0:0) with migration commands
- `make migrate-up` - Apply all pending migrations
- `make migrate-down` - Rollback last migration  
- `make migrate-version` - Check current migration version
- `make migrate-create name=add_table` - Create new migration files
- `make sqlc-generate` - Regenerate Go code from SQL

## Workflow

**To add a new database change:**
1. `make migrate-create name=add_new_table` - Creates `000002_add_new_table.up.sql` and `.down.sql`
2. Edit the migration files with your schema changes
3. `make sqlc-generate` - Regenerates Go code with new schema
4. `make up` - Docker compose will auto-run migrations before starting server