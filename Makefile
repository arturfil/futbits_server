DSN="host=localhost port=5432 user=postgres password=password dbname=chi_soccer sslmode=disable timezone=UTC connect_timeout=5"
BINARY_NAME=soccerApi

build:
	@echo "Building backend"
	go build -o ${BINARY_NAME} ./api
	@echo "Binary build!"

run: build
	@echo "Starting backend..."
	@env DSN=${DSN} ./${BINARY_NAME} &
	@echo "Backend started!"

dropdb:
	docker-compose exec postgres psql -U postgres -d postgres -c "DROP DATABASE chi_soccer"

createdb:
	docker exec -it backend_postgres_1 createdb --username=postgres --owner=postgres chi_soccer

migrateup:
	migrate -path db/migration -database "postgresql://postgres:password@localhost:5432/chi_soccer?sslmode=disable" -verbose up

migratedown:
	migrate -path db/migration -database "postgresql://postgres:password@localhost:5432/chi_soccer?sslmode=disable" -verbose down

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