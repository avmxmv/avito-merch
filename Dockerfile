FROM golang:1.23.4

RUN apt-get update && apt-get install -y protobuf-compiler

RUN CGO_ENABLED=0 go install github.com/pressly/goose/v3/cmd/goose@latest

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -mod=mod -o avito-merch ./cmd/server

EXPOSE 8080

CMD ["./avito-merch"]