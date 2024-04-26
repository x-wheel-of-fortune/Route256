ifeq ($(POSTGRES_SETUP_TEST),)
	POSTGRES_SETUP_TEST := user=test password=test dbname=test host=localhost port=5431 sslmode=disable
endif

INTERNAL_PKG_PATH=$(CURDIR)/internal/pkg
MIGRATION_FOLDER=$(INTERNAL_PKG_PATH)/db/migrations

.PHONY: up

generate:
	rm -rf internal/pkg/pb
	mkdir -p internal/pkg/pb

	protoc \
		--proto_path=api/ \
		--go_out=internal/pkg/pb \
		--go-grpc_out=internal/pkg/pb \
		api/*.proto

up:
	docker-compose up

cleanup:
	docker-compose down \
        --rmi local \
		--volumes \
		--remove-orphans \
		--timeout 60; \
  	docker-compose rm -f

.PHONY: test-migration-up ## Запуск скрипта миграций
test-migration-up:
	goose -dir "$(MIGRATION_FOLDER)" postgres "$(POSTGRES_SETUP_TEST)" up

.PHONY: test-migration-down ## Очистка БД от тестовых данных
test-migration-down:
	goose -dir "$(MIGRATION_FOLDER)" postgres "$(POSTGRES_SETUP_TEST)" down