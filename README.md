# Real-Time Matchmaking & Battle Simulator

A backend prototype for simulating real-time 1v1 player battles using Go
and gRPC. This project showcases the design of concurrent systems with
blocking queues, context-based cancellation, and battle simulation
logic.

## 🔧 Features

-   **1v1 Auto Matchmaking**: Players are automatically matched in pairs
    through a blocking queue.
-   **Battle Simulation**: Randomized turn-based combat with HP, damage,
    and win condition logic.
-   **gRPC API**:
    -   `JoinQueue`: Enter matchmaking queue
    -   `GetBattleResult`: Retrieve result by match ID
-   **Concurrency Management**: Channels, mutexes, and context handling
-   **Dockerized**: Ready for containerized deployment

## 🧱 Tech Stack

-   Go 1.24+
-   gRPC / Protocol Buffers
-   Docker
-   Context, Goroutines, Mutex, Channels

## 📁 Project Structure

    grpc-hello/
    ├── proto/              # match.proto definitions
    ├── server/             # gRPC backend server
    │   ├── main.go         
    │   └── Dockerfile
    ├── client/             # Simple CLI gRPC client
    │   └── main.go
    ├── go.mod / go.sum

## 🚀 Getting Started

### Build & Run (server)

``` bash
docker build -t grpc-hello-server -f server/Dockerfile .
docker run -d -p 50051:50051 --name grpc-hello grpc-hello-server
```

### Run Client (locally)

``` bash
go run client/main.go UserA
```

> Note: Run twice (with different user IDs) to simulate a full match.

------------------------------------------------------------------------

## 📌 TODO (WIP)

-   [ ] Add Dockerfile for client
-   [ ] Add `docker-compose.yml`
-   [ ] Migrate deployment to Kubernetes
-   [ ] Add Prometheus + Grafana metrics
-   [ ] Add CI/CD pipeline with GitHub Actions

## 📃 License

MIT (feel free to fork or use)

## 🙋‍♂️ Author

Hyerin Yeum\
[GitHub Profile](https://github.com/yeum)

------------------------------------------------------------------------

> This project was created as part of a portfolio to demonstrate backend
> system design with real-time constraints and DevOps readiness.
