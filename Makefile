build:
	cd client && npm run build && cd ..


run:	build
	go build && ./budgie

image:
	DOCKER_DEFAULT_PLATFORM=linux/amd64 docker build .
