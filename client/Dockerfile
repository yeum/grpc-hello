FROM golang:1.24

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY proto proto
COPY client/client.go .

# RUN protoc --go_out=. --go-grpc_out=. proto/*.proto

RUN go build -o client client.go

CMD ["./client", "UserA"]