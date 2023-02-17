include .env

dependencies:
	docker pull ${DB_IMAGE}
	docker pull kjconroy/sqlc
	go get .
	go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest

build_postgres:
	docker run --name ${DB_CONTAINER_NAME} -p ${DB_HOST_PORT}:${DB_CONTAINER_PORT} -e POSTGRES_USER=${DB_USER} -e POSTGRES_PASSWORD=${DB_PASSWORD} -d ${DB_IMAGE}

start_postgres:
	docker start ${DB_CONTAINER_NAME}

stop_postgres:
	docker stop ${DB_CONTAINER_NAME}

remove_postgres:
	docker stop ${DB_CONTAINER_NAME}
	docker rm ${DB_CONTAINER_NAME}

create_db:
	docker exec -it ${DB_CONTAINER_NAME} createdb --username=${DB_USER} --owner=${DB_USER} ${DB_DATABASE_NAME}

create_test_db:
	docker exec -it ${DB_CONTAINER_NAME} createdb --username=${DB_USER} --owner=${DB_USER} ${TEST_DB_DATABASE_NAME}

migrate_up:
	migrate -path database/migration -database "postgresql://${DB_USER}:${DB_PASSWORD}@${DB_HOST}:${DB_HOST_PORT}/${DB_DATABASE_NAME}?sslmode=disable" -verbose up

migrate_test_up:
	migrate -path database/migration -database "postgresql://${DB_USER}:${DB_PASSWORD}@${DB_HOST}:${DB_HOST_PORT}/${TEST_DB_DATABASE_NAME}?sslmode=disable" -verbose up

migrate_down:
	migrate -path database/migration -database "postgresql://${DB_USER}:${DB_PASSWORD}@${DB_HOST}:${DB_HOST_PORT}/${DB_DATABASE_NAME}?sslmode=disable" -verbose down

migrate_test_down:
	migrate -path database/migration -database "postgresql://${DB_USER}:${DB_PASSWORD}@${DB_HOST}:${DB_HOST_PORT}/${TEST_DB_DATABASE_NAME}?sslmode=disable" -verbose down

drop_db: 
	docker exec -it ${DB_CONTAINER_NAME} dropdb ${DB_DATABASE_NAME}

drop_test_db:
	docker exec -it ${DB_CONTAINER_NAME} dropdb ${TEST_DB_DATABASE_NAME}

psql:
	docker exec -it ${DB_CONTAINER_NAME} psql -U ${DB_USER} ${DB_DATABASE_NAME}

psql_testdb:
	docker exec -it ${DB_CONTAINER_NAME} psql -U ${DB_USER} ${TEST_DB_DATABASE_NAME}

container_shell:
	docker exec -it ${DB_CONTAINER_NAME} /bin/sh

sqlc:
	docker run --rm -v "$(CURDIR):/src" -w /src kjconroy/sqlc generate

test-server:
	go test ./server -v -cover

test-config:
	go test ./config -v -cover

test-database:
	go test ./database -v -cover

test-generated:
	go test ./database/generated -v -cover
test-util:
	go test ./util -v -cover

test-verbose:
	go test ./... -v -cover

test-detailed-report:
	go test ./... -coverprofile=tests_rep.tmp
	go tool cover -html tests_rep.tmp

