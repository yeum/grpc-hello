FROM golang:1.24

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY proto proto
COPY server/server.go .

# RUN protoc --go_out=. --go-grpc_out=. proto/*.proto

RUN go build -o server server.go

CMD ["./server"]