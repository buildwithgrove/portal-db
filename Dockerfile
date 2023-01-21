FROM postgres:14.3

COPY ./postgres-driver/sqlc/schema.sql /docker-entrypoint-initdb.d/
COPY ./testdata/seed_test_db.sql /docker-entrypoint-initdb.d/
