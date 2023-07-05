build:
	go build -o bin/budgie cmd/server/main.go
	cd client && npm run build && cd ..

dev:
	go build -o bin/budgie cmd/server/main.go
	./bin/budgie

web:
	cd client && npm run dev && cd ..

run: build
	./bin/budgie

image:
	DOCKER_DEFAULT_PLATFORM=linux/amd64 docker build .

lint:
	golangci-lint run
