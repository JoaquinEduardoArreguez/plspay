BINARY_NAME = main
POSTGRES_DSN = "host=localhost user=postgres password=admin dbname=plspay port=5432 sslmode=disable"

build:
	go build -o ./builds/${BINARY_NAME}.out ./src/api/*

run: build
	echo "\n\n\n\n\n"
	clear
	./builds/${BINARY_NAME}.out -postgresDsn=${POSTGRES_DSN}

run-docker:
	docker-compose -f docker.compose.yml up -d

stop-docker:
	docker-compose -f docker.compose.yml down

debug-build:
	go build -o ./builds/${BINARY_NAME}.out -gcflags all=-N -l ./src/api/*

debug:
	/Users/joaquinarreguez/go/bin/dlv dap --listen=127.0.0.1:54415 --log-dest=3