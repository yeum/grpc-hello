export default {
    server: {
        port: 5173,
        proxy: {
            "/api": {target: "http://localhost:8080", changeOrigin: true},
            "/ws": {target: "http://localhost:8080", ws: true, changeOrigin: true}
        }
    },
    build: {outDir: "dist"}
}