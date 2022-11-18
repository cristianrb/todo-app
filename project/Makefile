BROKER_BINARY=cmd/brokerService
RECORD_BINARY=cmd/recordService
COLLECTOR_BINARY=cmd/collectorService

## up: starts all containers in the background without forcing build
up:
	echo "Starting Docker images..."
	docker-compose up -d
	echo "Docker images started!"

## up_build: stops docker-compose (if running), builds all projects and starts docker compose
up_build: build_broker build_record build_collector
	echo "Stopping docker images (if running...)"
	docker-compose down
	echo "Building (when required) and starting docker images..."
	docker-compose up --build -d
	echo "Docker images built and started!"

## down: stop docker compose
down:
	echo "Stopping docker compose..."
	docker-compose down
	echo "Done!"

## build_broker: builds the broker binary as a linux executable
build_broker:
	echo "Building broker binary..."
	cd ../broker-service && env GOOS=linux CGO_ENABLED=0 go build -o ${BROKER_BINARY} ./cmd
	echo "Done!"

## build_record: builds the record binary as a linux executable
build_record:
	echo "Building record binary..."
	cd ../record-service && env GOOS=linux CGO_ENABLED=0 go build -o ${RECORD_BINARY} ./cmd
	echo "Done!"

## build_collector: builds the collector binary as a linux executable
build_collector:
	echo "Building collector binary..."
	cd ../collector-service && env GOOS=linux CGO_ENABLED=0 go build -o ${COLLECTOR_BINARY} ./cmd
	echo "Done!"