ensure-reflex:
	@if ! command -v reflex >/dev/null 2>&1; then \
		echo "🔍 Reflex not found. Installing..."; \
		go install github.com/cespare/reflex@latest; \
		echo "✅ Reflex installed."; \
	fi

ensure-goose:
	@if ! command -v goose >/dev/null 2>&1; then \
		echo "🔍 Goose not found. Installing..."; \
		go install github.com/pressly/goose/v3/cmd/goose@latest; \
		echo "✅ Goose installed."; \
	fi

ensure-swagger:
	@if ! command -v swag >/dev/null 2>&1; then \
		echo "🔍 Swag not found. Installing..."; \
		go install github.com/swaggo/swag/cmd/swag@latest; \
		echo "✅ Swag installed."; \
	fi

ensure-mockgen:
	@if ! command -v mockgen >/dev/null 2>&1; then \
		echo "🔍 Mockgen not found. Installing..."; \
		go install github.com/golang/mock/mockgen@v1.6.0; \
		echo "✅ Mockgen installed."; \
	fi

# mockgen dependency
ensure-gomock: 
	@go get -d github.com/golang/mock/gomock@v1.6.0 >/dev/null 2>&1 || true

tidy:
	go mod tidy

mock: ensure-mockgen ensure-gomock
	go generate ./...

unit-test:
	go test -cover $$(go list ./... | grep -v -E '/(mocks|vendor|docs)')

migrate-up: ensure-goose
	goose -dir ./migrations -table migration_version postgres "postgres://doitpayuser:doitpaypassword@doitpay-postgresql:5432/transfer_db?sslmode=disable" up

migrate-down: ensure-goose
	goose -dir ./migrations -table migration_version postgres "postgres://doitpayuser:doitpaypassword@doitpay-postgresql:5432/transfer_db?sslmode=disable" down

migrate-create: ensure-goose
	@if [ -z "$(filter-out $@,$(MAKECMDGOALS))" ]; then \
		echo "❌ Please provide migration name. example: make migrate-create create_table_user"; \
		exit 1; \
	fi
	@goose -dir ./migrations create $(filter-out $@,$(MAKECMDGOALS)) sql

test: mock unit-test

pre-run: tidy mock unit-test

api-docs: ensure-swagger
	swag init -g server/http_router.go --parseDependency true --parseInternal true

run-http: migrate-up ensure-reflex
	reflex -r '\.go' -s -- sh -c "go run main.go serve-http"

run-worker: ensure-reflex
	reflex -r '\.go' -s -- sh -c "go run main.go serve-worker"