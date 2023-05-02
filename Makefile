build:
	go build -o bin/budgie
	cd client && npm run build && cd ..

dev:
	go build -o bin/budgie
	./bin/budgie

run: build
	./bin/budgie

image:
	DOCKER_DEFAULT_PLATFORM=linux/amd64 docker build .

lint:
	golangci-lint run