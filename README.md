# Real-Time 2D RPG Sample

> Note: This project was initially a 1v1 battle simulator using Go and gRPC, 
> and has now been expanded into a **2D RPG** with a Phaser frontend.

A backend + frontend prototype for simulating real-time multiplayer
gameplay. It showcases the design of concurrent systems, matchmaking
queues, and integration between a Go backend and a web-based client.

## 🔧 Features

-   **2D RPG Gameplay (Phaser)**: Tilemap-based map rendering with player
    movement and camera follow.
-   **Battle Logic**: Previously turn-based 1v1 simulation, now being
    adapted into real-time RPG interactions.
-   **gRPC API**:
    -   `JoinQueue`: Enter matchmaking queue (legacy)
    -   `GetBattleResult`: Retrieve result by match ID (legacy)
-   **Concurrency Management**: Channels, mutexes, and context handling
-   **Dockerized**: Ready for containerized deployment

## 🧱 Tech Stack

-   Go 1.24+
-   gRPC / Protocol Buffers
-   Phaser 3 (frontend, via Vite)
-   Docker / Docker Compose
-   Context, Goroutines, Mutex, Channels

## 📁 Project Structure

    grpc-hello/
    ├── proto/              # match.proto definitions
    ├── server/             # gRPC backend server
    │   ├── main.go         
    │   └── Dockerfile
    ├── client/             # Phaser 2D RPG client (Vite project)
    │   ├── src/
    │   ├── public/assets/
    │   ├── package.json
    │   └── vite.config.js
    ├── go.mod / go.sum
    ├── docker-compose.yaml
    └── README.md

## 🚀 Getting Started

### Option A - Local (Server in Docker, Client on host)
#### Build & Run (server)

``` bash
docker build -t grpc-hello-server -f server/Dockerfile .
docker run -d -p 50051:50051 --name grpc-hello grpc-hello-server
```

### Run Client (locally with Vite)

``` bash
cd client
npm install
npm run dev
```

### Option B - Docker Compose (Server + Client)

``` bash
docker compose up --build -d
```

## 📃 License

MIT (feel free to fork or use)

## 🙋‍♂️ Author

Hyerin Yeum\
[GitHub Profile](https://github.com/yeum)

---
