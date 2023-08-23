SHELL := /bin/bash

make: gen_sql gen_mocks mermaid

gen_sql:
	sqlc generate -f ./postgres-driver/sqlc/sqlc.yaml
gen_mocks:
	mockery --name=Driver --filename=mock_driver.go --recursive --inpackage
gen_mermaid:
	@./update-mermaid.sh > /dev/null 2>&1 || true

mermaid: test_env_up gen_mermaid test_env_down

test: test_unit test_env_up run_driver_tests test_env_down
test_e2e: test_env_up run_driver_tests test_env_down

test_unit:
	go test ./... -count=1 -short || true
run_driver_tests:
	go test ./... -run Test_RunPGDriverSuite -count=1 || true
run_all_tests:
	go test ./... -count=1 || true

test_unit_ci:
	go test ./... -count=1 -short
run_driver_tests_ci:
	go test ./... -run Test_RunPGDriverSuite -count=1

test_env_up:
	@echo "🧪 Starting up Portal DB test database ..."
	@docker-compose -f ./testdata/docker-compose.test.yml up -d --remove-orphans --build
	@echo "⏳ Waiting for test DB to be ready ..."
	@attempts=0; while ! pg_isready -h localhost -p 5432 -U postgres -d postgres >/dev/null && [[ $$attempts -lt 5 ]]; do sleep 1; attempts=$$(($$attempts + 1)); done
	@[[ $$attempts -lt 5 ]] && echo "🐘 Test Portal DB is up ..." || (echo "❌ Test Portal DB failed to start" && make test_env_down >/dev/null && exit 1)
	@echo "🚀 Test environment is up ..."
test_env_down:
	@echo "🧪 Shutting down Pocket HTTP DB test environment ..."
	@docker-compose -f ./testdata/docker-compose.test.yml down --remove-orphans >/dev/null
	@echo "✅ Test environment is down."
reset_test_db:
	@echo "🧨 Starting the database reset operation..."
	@docker exec -it test-database psql -U postgres -d postgres -f /scripts/reset_test_db.sql >/dev/null || (echo "❌ Database reset operation failed:" && docker exec -it test-database psql -U postgres -d postgres -f /docker-entrypoint-initdb.d/reset_test_db.sql)
	@echo "✅ Database reset operation completed successfully!"
test_dev: reset_test_db run_driver_tests

init:
	wget https://github.com/pre-commit/pre-commit/releases/download/v2.20.0/pre-commit-2.20.0.pyz;
	python3 pre-commit-2.20.0.pyz install;
	python3 pre-commit-2.20.0.pyz autoupdate;
	go install golang.org/x/tools/cmd/goimports@v0.6.0;
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@v1.51.0;
	go install -v github.com/go-critic/go-critic/cmd/gocritic@v0.6.5;
	go install github.com/KarnerTh/mermerd@v0.9.0;
	go install github.com/sqlc-dev/sqlc/cmd/sqlc@v1.20.0;
	go install github.com/vektra/mockery/v2@v2.33.0;
	trap 'rm -f pre-commit-2.20.0.pyz*;' EXIT; \
	python3 pre-commit-2.20.0.pyz run --all-files;
