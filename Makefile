ENV=.env
include $(ENV)

MIGRATIONS_PATH=database/migrations/
DB_PATH=postgresql://$(PG_USER):$(PG_PASS)@$(PG_HOST):$(PG_PORT)/$(PG_NAME)?sslmode=disable

create_migration:
	migrate create -ext sql -dir $(MIGRATIONS_PATH) $(MIGRATION_NAME)

migration_up: 
	migrate -path $(MIGRATIONS_PATH) -database $(DB_PATH) -verbose up

migration_down: 
	migrate -path $(MIGRATIONS_PATH) -database $(DB_PATH) -verbose down;

migration_fix: 
	migrate -path $(MIGRATIONS_PATH) -database $(DB_PATH) force VERSION;