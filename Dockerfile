# Build stage
FROM golang:1.20 AS build

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY *.go ./
COPY server/* ./server/
COPY storage/* ./storage/
COPY utils/* ./utils/
COPY config/* ./config/
COPY models/* ./models/

RUN CGO_ENABLED=0 go build -o budgie

# Final stage
FROM scratch

COPY --from=build /app/budgie /app/
COPY migrations/* ./migrations/
COPY client/* ./client/

EXPOSE 3000

CMD ["/app/budgie"]
