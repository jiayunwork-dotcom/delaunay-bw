# delaunay-bw：平面点集 Bowyer–Watson 三角剖分与 Voronoi 对偶核算工具

Go 同进程托管静态前端（web/）与 JSON API；启动后访问 http://localhost:8080/，可加载 example/grid-jitter.json 再计算。

## 构建 / 运行 / 测试

```text
go build ./...
go run . -http :8080                 # 打开 http://localhost:8080，加载 example 再算
go test ./...
```

## 评测镜像

本目录评测专用文件（勿覆盖项目自带 Dockerfile/README）：

- `benzhi.Dockerfile`
- `build_benzhi_docker.sh`
- `BENZHI_README.md`（本文件）

两种架构都要构建并进容器验证：

```bash
chmod +x build_benzhi_docker.sh
./build_benzhi_docker.sh <image-name> linux/arm64
./build_benzhi_docker.sh <image-name> linux/amd64
docker run -it <image-name>:latest
```
