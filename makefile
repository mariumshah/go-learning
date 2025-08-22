SRC_DIR=.
DST_DIR=.
MODEL_DIR=pkg/models

# Command to force-remove files/directories
CMD_RM = rm -rf

MIGRATE_CFG = resources/migrations/app/dbconfig.yml
SCHEMA_CFG = resources/migrations/schema/dbconfig.yml
# runs sql migrations
CMD_MIGRATE = sql-migrate
# generates ORM models
CMD_BOILER = sqlboiler -c resources/migrations/sqlboiler.toml -o $(MODEL_DIR)
# purpose????
# LDFLAGS = -ldflags "-X main.Version=$(VERSION) -X main.Build=$(BUILD) -X main.GitHash=$(GIT_HASH)"
LDFLAGS = -ldflags "-X github.com/safepay/userapi/pkg/utils.Version=$(VERSION)"

CMD_GO = go

.PHONY: db-all
# full reset and migrate
db-all: db-migrate-down
	make db-migrate-up n=0
	make db-generate-models

.PHONY: db-migrate-status
# shows current migration status for schema and app migrations
db-migrate-status:
	${CMD_MIGRATE} status -config=$(SCHEMA_CFG)
	${CMD_MIGRATE} status -config=$(MIGRATE_CFG)

.PHONY: db-migrate-schema-up
db-migrate-schema-up:
	${CMD_MIGRATE} up -config=$(SCHEMA_CFG)

.PHONY: db-migrate-app-up
db-migrate-app-up:
	${CMD_MIGRATE} up -config=$(MIGRATE_CFG)

.PHONY: db-migrate-down
# Roll back the last applied app migration
db-migrate-down:
	$(CMD_MIGRATE) down -config=$(MIGRATE_CFG)

.PHONY: db-migrate-schema-down
# Roll back the last applied app migration
db-migrate-schema-down:
		$(CMD_MIGRATE) down -config=$(SCHEMA_CFG)

.PHONY: db-clean-models
db-clean-models:
	${CMD_BOILER} mysql --wipe
	${CMD_RM} $(MODEL_DIR)

.PHONY: db-generate-models
# Regenerate models using SQLBoiler (skipping hooks & tests)
db-generate-models: db-clean-models
	$(CMD_BOILER) mysql --no-hooks --no-tests

# check
PHONY: run-%
# Pattern rule: `make run-foo` will `go run cmd/foo/main.go`
run-%:
	$(CMD_GO) run -v $(LDFLAGS) cmd/$*/main.go
