include .env

postgres:
	docker run --name ${DB_DOCKER_CONTAINER} -p 5432:5432 -e POSTGRES_USER=${DB_USER}-e POSTGRES_PASSWORD=${DB_PASSWORD} -d postgres:12-alpine
# creates the db withing the postgres container
createdb:
	docker exec -it ${DB_DOCKER_CONTAINER} createdb --username=${DB_USER} --owner=${DB_USER} ${DB_NAME}

create-migrations-init:
	sqlx migrate add -r init

migrate-up:
	sqlx migrate run --database-url "postgresql://${DB_USER}:${DB_PASSWORD}@${DB_HOST}:${DB_PORT}/${DB_NAME}?sslmode=disable" 

migrate-down:
	sqlx migrate revert --database-url "postgresql://${DB_USER}:${DB_PASSWORD}@${DB_HOST}:${DB_PORT}/${DB_NAME}?sslmode=disable" 

create-migrations-seed:
	sqlx migrate add -r seed 

seed_data:
	sqlx migrate run --database-url "postgresql://${DB_USER}:${DB_PASSWORD}@${DB_HOST}:${DB_PORT}/${DB_NAME}?sslmode=disable" 

force_flag_false:
	migrate -path migrations -database "postgresql://${DB_USER}:${DB_PASSWORD}@${DB_HOST}:${DB_PORT}/${DB_NAME}?sslmode=disable" force 1

sign_up:
	curl -X POST http://localhost:8080/api/v1/auth/signup \
	-H 'Content-Type: application/json' \
	-d '{ \
		"first_name": "Arturo", \
		"last_name": "Filio", \
		"email": "arturo@test.com", \
		"password": "Password123" \
	}' \

build:
	@echo "Building backend"
	go build -o ${BINARY_NAME} cmd/server/*.go
	@echo "Binary build!"

buildbackend:
	env DSN=${DSN}
	env GOOS=linux GOARCH=amd64 go build -o futbitsProd cmd/server/*.go

stop_containers:
	@echo "Stoping all docker containers..."
	if [ $$(docker ps -q) ]; then \
		echo "found and stopped containers..."; \
		docker stop $$(docker ps -q); \
	else \
		echo "no active containers found..."; \
	fi

docker-run:
	@echo "Running docker images"
	docker-compose up --build -d

docker-stop:
	@echo "\nStopping all images\n"
	docker-compose stop

start-docker:
	docker start ${DB_DOCKER_CONTAINER}

run: build stop_containers start-docker
	@echo "Starting db docker container"
	docker start ${DB_DOCKER_CONTAINER}
	@echo "Starting backend..."
	@env PORT=${PORT} DSN=${DSN} ./${BINARY_NAME}  &
	@echo "Backend started!"

run-prod: 	
	@echo "Starting backend..."
	@env PORT=${PORT} DSN=${DSN} ./${BINARY_NAME}  &
	@echo "Backend started!"


dirtflagfalse:
	docker exec -it backend_postgres_1 update schema_migrations set dirty = false

dropdb:
	docker exec -it ${DB_DOCKER_CONTAINER} psql -U root -d postgres -c "DROP DATABASE chi_soccerdb"

clean:
	@echo "Cleaning..."
	@go clean
	@rm ${BINARY_NAME}
	@echo "Cleaned!"

start: run

stop:
	@echo "Stopping backend"
	@-pkill -SIGTERM -f "./${BINARY_NAME}"
	@ECHO "Stopped backend"

restart: stop start
