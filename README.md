# 简单记账后端

## 使用 Docker Compose 部署
1. 构建 Docker 镜像，在项目根目录下执行`docker buildx build --platform linux/amd64 -t easywrite-service -f docker/Dockerfile .`，如果你是 MacOS 苹果硅芯片，请将平台替换为`linux/arm64`
2. 拷贝`docker/release/docker-compose.yml`文件并根据你的实际情况修改它，然后`docker-compose up -d`或者`docker compose up -d`启动它
3. 默认端口为`10010`, 可根据需要修改，推荐使用 `Caddy` 等此类Web网关反向代理到 `10010` 端口

## 使用 GITHUB ACTION 构建 Docker 镜像并上传到镜像仓库，使用 [.github/workflows/deploy.yml](.github/workflows/deploy.yml)
## 编译
```shell
 $env:GOOS="linux" ; $env:GOARCH="arm64" ; go build -o easywrite_service cmd/main.go
```
```shell
 $env:GOOS="linux" ; $env:GOARCH="amd64" ; go build -o easywrite_service  ./main.go
```

## 配置文件名模板在 [service_config_tmplate.json5](docker%2Fdev%2Fservice_config_tmplate.json5)

## 主函数在 main.go
```shell
./easywrite_service -g #生成模板配置文件
```
```shell
./easywrite_service -c /path/to/config.json #指定配置文件，默认路径为 ./service_config.json5 , 使用默认路径可以不指定 -c
```

## 生成文档
```shell
swag init --parseDependency --parseInternal
```
