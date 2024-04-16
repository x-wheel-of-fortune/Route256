ifeq ($(POSTGRES_SETUP_TEST),)
	POSTGRES_SETUP_TEST := user=test password=test dbname=test host=localhost port=5432 sslmode=disable
endif

INTERNAL_PKG_PATH=$(CURDIR)/internal/pkg
MOCKGEN_TAG=1.2.0
MIGRATION_FOLDER=$(INTERNAL_PKG_PATH)/db/migrations

.PHONY: start-test-env ## Запуск тестового окружения при помощи docker-compose
start-test-env:
	docker-compose up -d

.PHONY: stop-test-env
stop-test-env:
	docker-compose down

.PHONY: migration-create
migration-create:
	goose -dir "$(MIGRATION_FOLDER)" create "$(name)" sql

.PHONY: test-migration-up ## Запуск скрипта миграций
test-migration-up:
	goose -dir "$(MIGRATION_FOLDER)" postgres "$(POSTGRES_SETUP_TEST)" up

.PHONY: test-migration-down ## Очистка БД от тестовых данных
test-migration-down:
	goose -dir "$(MIGRATION_FOLDER)" postgres "$(POSTGRES_SETUP_TEST)" down

.PHONY: .generate-mockgen-deps
.generate-mockgen-deps:
ifeq ($(wildcard $(MOCKGEN_BIN)),)
	@GOBIN=$(LOCAL_BIN) go install github.com/golang/mock/mockgen@$(MOCKGEN_TAG)
endif

.PHONY: .generate-mockgen
.generate-mockgen:
	PATH="$(LOCAL_BIN):$$PATH" go generate -x -run=mockgen ./...

.PHONY: gofmt
gofmt:
	goimports -l -w $(CURDIR)
.test:
	$(info Running tests...)
	go test ./...

test: .test ## запуск всех юнит тестов

.itest:
	$(info Running integration tests...)
	go test -tags=integration ./tests/...

itest: .itest ## запуск всех интеграционных тестов


