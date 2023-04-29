build:
	cd client && npm run build && cd ..

run:
	cd client && npm run build && cd .. && go build && ./budgie

budgie:
	go build && ./budgie

image:
	DOCKER_DEFAULT_PLATFORM=linux/amd64 docker build .

lint:
	golangci-lint run