# Comment Anywhere Back End

This repository hosts the Comment Anywhere back-end. It consists of two components, the Database Server and the HTTP Server. This repository can build both to Docker containers with make commands. Go code is generated from database queries using sqlc and used with API endpoints to handle.


## Example .env file

```env
# The .env file provides important constants. It should not be included in public repos. 

CA_TESTING_MODE = true

# DB_IMAGE is the docker image that the database will be built from.
DB_IMAGE=postgres:14.5-alpine
# DB_CONTAINER_NAME is the name of the container in docker.
DB_CONTAINER_NAME=923postgres
# DB_CONTAINER_PORT is the port the container will be listening on in its environment. It will be mapped to the value of DB_HOST_PORT on the host device. These values can be the same but if you already have postgres installed on your computer like I do, you may want to map to a different port.
DB_CONTAINER_PORT=5432

DB_HOST=localhost


# -------------- Production Environment -------------------------
# The port that Go will listen to database with(for port mapping)
DB_HOST_PORT=5433
DB_USER=root
DB_PASSWORD=dbsuperuser991
DB_DATABASE_NAME=comm-anything

# SERVER_PORT is the port that Go will be served on.
SERVER_PORT=3000
SERVER_LOG_ALL=true

# -------------- Development Environment -------------------------

TEST_DB_HOST_PORT = 5434
TEST_DB_USER = root
TEST_DB_PASSWORD = dbsuperuser991
TEST_DB_DATABASE_NAME = comm-anything-tests

```

