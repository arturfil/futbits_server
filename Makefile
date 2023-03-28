DSN="host=localhost port=5432 user=root password=secret dbname=chi_soccerdb sslmode=disable timezone=UTC connect_timeout=5" BINARY_NAME=soccerApi
PORT=8080
DB_DOCKER_CONTAINER=chi_soccer
BINARY_NAME=soccerApi
SECRET_KEY=asdkjq234-081234j-lkasdf82314-32jlkjadsf0-891234ljasdf0-143jlaksdf

# creates container with postgres software
postgres:
	docker run --name ${DB_DOCKER_CONTAINER} -p 5432:5432 -e POSTGRES_USER=root -e POSTGRES_PASSWORD=secret -d postgres:12-alpine
# creates the db withing the postgres container
createdb:
	docker exec -it ${DB_DOCKER_CONTAINER} createdb --username=root --owner=root chi_soccerdb

migrateup:
	migrate -path migrations -database "postgresql://root:secret@localhost:5432/chi_soccerdb?sslmode=disable" -verbose up

migratedown:
	migrate -path migrations -database "postgresql://root:secret@localhost:5432/chi_soccerdb?sslmode=disable" -verbose down

seed_data:
	migrate -path migrations -database "postgresql://root:secret@localhost:5432/chi_soccerdb?sslmode=disable" -verbose up 2

force_flag_false:
	migrate -path migrations -database "postgresql://root:secret@localhost:5432/chi_soccerdb?sslmode=disable" force 1

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

run: build
	@echo "Starting db docker container"
	docker start ${DB_DOCKER_CONTAINER}
	@echo "Starting backend..."
	@env PORT=${PORT} DSN=${DSN} ./${BINARY_NAME}  &
	@echo "Backend started!"

buildbackend:
	set GOOS=linux&& set GOARCH=amd64&& set CGO_ENABLED=0 && go build -o ${BINARY_NAME} ./api

dirtflagfalse:
	docker exec -it backend_postgres_1 update schema_migrations set dirty = false

dropdb:
	docker-compose exec postgres psql -U postgres -d postgres -c "DROP DATABASE chi_soccer"

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
