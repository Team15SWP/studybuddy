current_dir := $(dir $(mkfile_path))

.PHONY: test
test:
	@echo "Testing app..."
	@go test ./...

.PHONY: lint
lint:
ifeq (,$(shell which golangci-lint))
	@$(call WARN, "golangci-lint не найден. Устанавливаю...")
	@go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	@$(call INFO, "Установка линтера завершена.")
	@echo "-------------------------------------"
endif
	@$(call INFO, "Запуск линтера...")
	@golangci-lint run
	@$(call INFO, "Линтер закончил работу")
	@echo "-------------------------------------"

MIGRATION_FOLDER=$(CURDIR)/migrations

ifeq ($(POSTGRES_SETUP),)
	POSTGRES_SETUP := user=test password=test dbname=db_study_buddy host=localhost port=5432 sslmode=disable
endif

.PHONY: migration-create
migration-create:
	goose -dir "$(MIGRATION_FOLDER)" create "$(name)" sql

.PHONY: migration-up
migration-up:
	goose -dir "$(MIGRATION_FOLDER)" postgres "$(POSTGRES_SETUP)" up

.PHONY: migration-down
migration-down:
	goose -dir "$(MIGRATION_FOLDER)" postgres "$(POSTGRES_SETUP)" down

.PHONY: migration-down-to
migration-down-to:
	goose -dir "$(MIGRATION_FOLDER)" postgres "$(POSTGRES_SETUP)" down-to "$(version)"

.PHONY: run
run:
	docker compose up -d
	docker compose exec study_buddy_db sh -c 'while ! pg_isready -U $(POSTGRES_USER) -d $(POSTGRES_DB); do sleep 5; done'
	goose -dir "$(MIGRATION_FOLDER)" postgres "$(POSTGRES_SETUP)" up

.PHONY: stop
stop:
	docker compose down

