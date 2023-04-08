# Comment Anywhere Back End

This repository hosts the Comment Anywhere back-end. It consists of two components, the Database Server and the HTTP Server. This repository can build both to Docker containers with make commands. Go code is generated from database queries using sqlc and used with API endpoints to handle.

## Example .env file

When you clone this repo for the first time, you will need to add the .env file. It is not uploaded to GitHub for security reasons. Here is a sample .env file. You should alter things like DB_USER, DB_PASSWORD, DB_NAME, and JWT_KEY for security. It wouldn't hurt to change the ports either. The .env file lives in the root of the project.

```env
# -------------- General Options ---------------------------

# The .env file provides important constants. It should not be included in public repos. 
CA_TESTING_MODE=true

# If true, server will implement a middleware that will log some information about every incoming HTTP Request to its console.. 
SERVER_LOG_ALL=true


# -------------- Database Configuration ---------------------------

# DB_IMAGE is the docker image that the VM containing the database will be built from.
DB_IMAGE=postgres:14.5-alpine
# DB_CONTAINER_NAME is the name of the container in docker.
DB_CONTAINER_NAME=923postgres
# DB_CONTAINER_PORT is the port the container will be listening on in its environment. It will be mapped to the value of DB_HOST_PORT on the host device. These values can be the same but if you already have postgres installed on your computer like I do, you may want to map to a different port.
DB_CONTAINER_PORT=5432
# This may have to be changed at deployment time. 
DB_HOST=localhost
# DB_HOST_PORT is the port the database is actually served on. It's the port that Go will listen to database with(for port mapping)
DB_HOST_PORT=5433
DB_USER=root
DB_PASSWORD=dbsuperuser991

# -------------- Server Configuration ---------------------------

# SERVER_PORT is the port that the HTTP server is to be served on.
SERVER_PORT=3000

# JWT_KEY is used for encrypting cookie/tokens which keep users logged in
JWT_KEY=some random sentence for encrypting and decrypting cookie tokens

# The name of the cookie to store in a user's browser. 
JWT_COOKIE_NAME=ca-auth-tok

# -------------- Production Environment -------------------------
DB_DATABASE_NAME=comm-anything

# -------------- Development Environment -------------------------
TEST_DB_DATABASE_NAME = comm-anything-tests

```

## Running

To start the comment anywhere back-end for the first time, you have to do the following actions:
1. Pull the needed dependencies
    1. `make dependencies`
2. Create the docker postgres container and start it
    1. `make build_postgres`
    2. `make create_db`
3. Load the schema into the postgres database
    1. `make migrate_up`
4. Start the server
    1. `go run .`

## Deploying

CommentAnywhere deploys as two docker containers connected by a bridge network. One container is the server, one container is the postgres database.

1. Copy the source files to the server 
 - run the following `scp` command:
 - `scp -r * user@111.11.111.11:/COMANY/gosrc`
2. Copy the env file to the server.
 - `scp .env user@111.11.111.11:/COMANY/gosrc/.env`
3. SSH into the server
 - `ssh user@111.11.111.11`
4. CD into `/COMANY/gosrc`
5. Run appropriate `make` commands to build the server and postgres containers
   - `make dependencies`
   - `make build_postgres`
   - `make create_db`
   - `make migrate_up`
   - `make build_server_image`
   - `make create_net`
   - `make net_con_db`
   - `make create_server`
   - `make net_con_server`
   - `make start_server`

The server should be running!


