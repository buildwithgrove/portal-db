make: gen_sql gen_mocks

gen_sql:
	sqlc generate -f ./postgres-driver/sqlc/sqlc.yaml
gen_mocks:
	mockery --name=Driver --filename=mock_driver.go --recursive --inpackage

test: test_unit test_env_up run_driver_tests test_env_down
test_driver: test_env_up run_driver_tests test_env_down
test_unit:
	go test ./... -count=1 -short;

test_env_up:
	docker-compose -f ./testdata/docker-compose.test.yml up -d --remove-orphans --build
	@echo "⏳ Waiting for test DB to be ready ..."
	until pg_isready -h localhost -p 5432 -U postgres -d postgres >/dev/null 2>&1; do sleep 0.01; done
	@echo "🚀 Test environment is up ..."
test_env_down:
	docker-compose -f ./testdata/docker-compose.test.yml down --remove-orphans -v
	@echo "✅ Test environment is down."
run_driver_tests:
	-go test ./... -run Test_RunPGDriverSuite -count=1;
run_all_tests:
	-go test ./... -count=1;

init-pre-commit:
	wget https://github.com/pre-commit/pre-commit/releases/download/v2.20.0/pre-commit-2.20.0.pyz;
	python3 pre-commit-2.20.0.pyz install;
	python3 pre-commit-2.20.0.pyz autoupdate;
	go install golang.org/x/tools/cmd/goimports@v0.6.0;
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@v1.51.0;
	go install -v github.com/go-critic/go-critic/cmd/gocritic@v0.6.5;
	python3 pre-commit-2.20.0.pyz run --all-files;
	rm pre-commit-2.20.0.pyz;
