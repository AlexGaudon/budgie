build:
	cd client && npm run build && cd ..
	go build


run:
	go build && ./budgie

image:
	DOCKER_DEFAULT_PLATFORM=linux/amd64 docker build .