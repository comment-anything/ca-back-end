include .env

# ---------------------------------------------------------------
#						    Database, Docker
# ---------------------------------------------------------------

# Pulls dependencies needed to build the docker image, including go dependencies.
dependencies:
	docker pull ${DB_IMAGE}
	docker pull kjconroy/sqlc
	go get .
	go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest

# Builds and starts the container.
build_postgres:
	docker run --name ${DB_CONTAINER_NAME} -p ${DB_HOST_PORT}:${DB_CONTAINER_PORT} -e POSTGRES_USER=${DB_USER} -e POSTGRES_PASSWORD=${DB_PASSWORD} -d ${DB_IMAGE}

# Starts the container (if it was stopped)
start_postgres:
	docker start ${DB_CONTAINER_NAME}

# Stops the container (if it was stopped)
stop_postgres:
	docker stop ${DB_CONTAINER_NAME}

# Removes the container
remove_postgres:
	docker stop ${DB_CONTAINER_NAME}
	docker rm ${DB_CONTAINER_NAME} --force

# Creates the Comment Anywhere database
create_db:
	docker exec -it ${DB_CONTAINER_NAME} createdb --username=${DB_USER} --owner=${DB_USER} ${DB_DATABASE_NAME}

# Creates a parallel Comment Anywhere testing database
create_test_db:
	docker exec -it ${DB_CONTAINER_NAME} createdb --username=${DB_USER} --owner=${DB_USER} ${TEST_DB_DATABASE_NAME}

# Uses migrate to execute database migrations
migrate_up:
	migrate -path database/migration -database "postgresql://${DB_USER}:${DB_PASSWORD}@${DB_HOST}:${DB_HOST_PORT}/${DB_DATABASE_NAME}?sslmode=disable" -verbose up

migrate_test_up:
	migrate -path database/migration -database "postgresql://${DB_USER}:${DB_PASSWORD}@${DB_HOST}:${DB_HOST_PORT}/${TEST_DB_DATABASE_NAME}?sslmode=disable" -verbose up

migrate_down:
	migrate -path database/migration -database "postgresql://${DB_USER}:${DB_PASSWORD}@${DB_HOST}:${DB_HOST_PORT}/${DB_DATABASE_NAME}?sslmode=disable" -verbose down

migrate_test_down:
	migrate -path database/migration -database "postgresql://${DB_USER}:${DB_PASSWORD}@${DB_HOST}:${DB_HOST_PORT}/${TEST_DB_DATABASE_NAME}?sslmode=disable" -verbose down

# Drops the main database DANGEROUS
drop_db: 
	docker exec -it ${DB_CONTAINER_NAME} dropdb ${DB_DATABASE_NAME}

# Drops the test database 
drop_test_db:
	docker exec -it ${DB_CONTAINER_NAME} dropdb ${TEST_DB_DATABASE_NAME}

# Enters the psql command line where you can run arbitrary psql commands
psql:
	docker exec -it ${DB_CONTAINER_NAME} psql -U ${DB_USER} ${DB_DATABASE_NAME}

# Enters the psql command line where you can run arbitrary psql commands
psql_testdb:
	docker exec -it ${DB_CONTAINER_NAME} psql -U ${DB_USER} ${TEST_DB_DATABASE_NAME}

# ---------------------------------------------------------------
#						    Shells
# ---------------------------------------------------------------

# Enters the container shell
container_shell:
	docker exec -it ${DB_CONTAINER_NAME} /bin/sh

# Runs sqlc to generate ORM code from database schema and queries
sqlc:
	docker run --rm -v "$(CURDIR):/src" -w /src kjconroy/sqlc generate

# ---------------------------------------------------------------
#						    Tests
# ---------------------------------------------------------------

# Runs server tests
test-server:
	go test ./server -v -cover

# Runs config tests
test-config:
	go test ./config -v -cover

# Runs database tests
test-database:
	go test ./database -v -cover

# Runs tests on generated code
test-generated:
	go test ./database/generated -v -cover

# Runs tests on package util
test-util:
	go test ./util -v -cover

# Runs all tests with verbose output
test-verbose:
	go test ./... -v -cover

# Runs all tests and presents a detailed report
test-detailed-report:
	go test ./... -coverprofile=tests_rep.tmp
	go tool cover -html tests_rep.tmp


# ---------------------------------------------------------------
#						    Deploying
# ---------------------------------------------------------------

# -- A test container for confirming bridging
test_container:
	docker run -t -p 3000 --name ${SERVER_CONTAINER_NAME} -d alpine:3.17.3

# -- Build the custom server image
build_server_image:
	docker build -t ${SERVER_IMAGE} --rm . \
	--build-arg port=${SERVER_PORT} 
	
# -- Create the network, connect the containers

create_net:
	docker network create -d bridge ${SERVER_DB_NETWORK}

net_con_server:
	docker network connect ${SERVER_DB_NETWORK} ${SERVER_CONTAINER_NAME}

net_con_db:
	docker network connect ${SERVER_DB_NETWORK} ${DB_CONTAINER_NAME}

net_dc_server:
	docker network disconnect ${SERVER_DB_NETWORK} ${SERVER_CONTAINER_NAME}

net_dc_db:
	docker network disconnect ${SERVER_DB_NETWORK} ${DB_CONTAINER_NAME}

remove_net:
	docker network rm ${SERVER_DB_NETWORK}

# -- Create server


create_server:
	docker run -t --name ${SERVER_CONTAINER_NAME}  \
	-p ${SERVER_PORT}:${SERVER_PORT} -d ${SERVER_IMAGE}

start_server:
	docker start ${SERVER_CONTAINER_NAME}

stop_server:
	docker stop ${SERVER_CONTAINER_NAME}

sc_server:
	$(MAKE) create_server
	$(MAKE) net_con_server
	$(MAKE) start_server

rm_server:
	docker rm ${SERVER_CONTAINER_NAME} --force

rebuild_server:
	$(MAKE) build_server_image
	$(MAKE) sc_server






# ------------------------
# 		Deprecated
# ------------------------





# --- Older stuff for when trying to run go inside of remote server rather than container

buildgo:
	go build .

# runs ca-back-end as a headless process
headless:
	$(MAKE) buildgo
	systemd-run --unit=serve-comm-any --working-directory=/COMANY/gosrc ./ca-back-end --nocli=true

kill_headless:
	systemctl kill serve-comm-any.service

ca_cli:
	journalctl -xfu serve-comm-any.service


# really just for reference; this command has to be run manually and can't be done through make
send_source:
	scp -r * $(REMOTE_USER)@$(IP_DEPLOY):$(REMOTE_PATH)
