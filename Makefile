migrate:
	@echo "migrating...."
	goose -dir ./internal/repository/postgres/migrations postgres "user=postgres password=postgres dbname=users sslmode=disable" up
generate-mocks:
	mockery --case=underscore --dir ./internal/repository --output ./internal/mocks/repository --all --keeptree
#TODO: exclude filter factory
generate:
	cd api/proto; buf generate

run:
	go run cmd/main/main.go