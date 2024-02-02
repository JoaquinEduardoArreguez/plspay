BINARY_NAME = main
POSTGRES_DSN = "host=localhost user=postgres password=admin dbname=plspay port=5432 sslmode=disable"

build:
	go build -o ./builds/${BINARY_NAME}.out ./src/api/*

run: build
	./builds/${BINARY_NAME}.out -postgresDsn=${POSTGRES_DSN}

run-docker:
	docker-compose -f docker.compose.yml up -d

run-prod:
	git pull
	sleep 2
	make run-docker
	sleep 2
	nohup make run &

stop-docker:
	docker-compose -f docker.compose.yml down

debug-build:
	go build -gcflags=all="-N -l" -o ./builds/${BINARY_NAME}.debug.out ./src/api/*

debug: debug-build
	clear
	./builds/${BINARY_NAME}.debug.out -postgresDsn=${POSTGRES_DSN}

test:
	go test ./...