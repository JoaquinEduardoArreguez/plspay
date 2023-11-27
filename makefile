BINARY_NAME = main

build:
	go build -o ./builds/${BINARY_NAME}.out ./src/api/*

run: build
	./builds/${BINARY_NAME}.out

run-docker:
	docker-compose -f docker.compose.yml up -d

stop-docker:
	docker-compose -f docker.compose.yml down