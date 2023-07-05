# Build stage
FROM golang:1.20-alpine AS build
WORKDIR /app
COPY . .
RUN ls && go build -o bin/app cmd/server/main.go

# Deploy stage
FROM alpine:3.13
WORKDIR /app
COPY --from=build /app/bin/app .
COPY .env .env
COPY ./migrations ./migrations
COPY ./client ./client
CMD ["./app"]
