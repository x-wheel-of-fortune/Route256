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
