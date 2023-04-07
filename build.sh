make dependencies
make build_postgres
make create_db
make migrate_up
go install .
ca-back-end
