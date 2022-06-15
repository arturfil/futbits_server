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