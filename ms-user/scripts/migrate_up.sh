export PG="postgres://postgres:postgres@localhost:5432/postgres?sslmode=disable" && migrate -database ${PG} -path ./database/migrations up