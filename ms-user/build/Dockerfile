FROM golang:1.22-alpine

WORKDIR /app

COPY go.* /app/

RUN go mod download

COPY . /app/

RUN go build -o main cmd/main.go

CMD ["/app/main"]