FROM golang:1.24

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY proto proto
COPY server/server.go .

RUN CGO_ENABLED=0 GOOS=linux go build -0 /server ./server

RUN go build -o server server.go

#Run time
FROM gcr.io/distroless/base-debian12
COPY --from=build /server /server
EXPOSE 50051
USER nonroot:nonroot
ENTRYPOINT [ "/server" ]