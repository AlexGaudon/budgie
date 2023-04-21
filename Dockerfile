FROM golang:1.20

WORKDIR /app

ADD go.mod .
ADD go.sum .

RUN go mod download

COPY *.go ./
COPY server/* ./server/
COPY storage/* ./storage/
COPY utils/* ./utils/
COPY config/* ./config/
COPY client/dist/* ./client/dist

RUN CGO_ENABLED=0 GOOS=linux go build -o /budgie

EXPOSE 3000

CMD ["/budgie"]