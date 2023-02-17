make: gen_sql gen_mocks

gen_sql:
	sqlc generate -f ./postgres-driver/sqlc/sqlc.yaml
gen_mocks:
	mockery --name=Driver --filename=mock_driver.go --recursive --inpackage

test: test_env_up run_driver_tests test_env_down
test_env_up:
	docker-compose -f ./docker-compose.test.yml up -d --remove-orphans --build;
	sleep 2;
test_env_down:
	docker-compose -f ./docker-compose.test.yml down --remove-orphans -v
run_driver_tests:
	-go test ./... -run Test_RunPGDriverSuite -count=1 -v;

init-pre-commit:
	wget https://github.com/pre-commit/pre-commit/releases/download/v2.20.0/pre-commit-2.20.0.pyz;
	python3 pre-commit-2.20.0.pyz install;
	python3 pre-commit-2.20.0.pyz autoupdate;
	go install golang.org/x/tools/cmd/goimports@v0.6.0;
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@v1.51.0;
	go install -v github.com/go-critic/go-critic/cmd/gocritic@v0.6.5;
	python3 pre-commit-2.20.0.pyz run --all-files;
	rm pre-commit-2.20.0.pyz;
