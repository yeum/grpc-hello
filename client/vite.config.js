export default {
    server: {
        port: 5173,
        proxy: {
            "/grpc": {target: "http://localhost:8080", changeOrigin: true},
        }
    },
    build: {outDir: "dist"}
}