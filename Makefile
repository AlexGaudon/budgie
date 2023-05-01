build:
	go build -o bin/budgie
	cd client && npm run build && cd ..


run: build
	./bin/budgie

image:
	DOCKER_DEFAULT_PLATFORM=linux/amd64 docker build .

lint:
	golangci-lint run