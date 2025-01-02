# 简单记账后端

## 编译
```shell
 $env:GOOS="linux" ; $env:GOARCH="arm64" ; go build -o easywrite_service cmd/main.go
```
```shell
 $env:GOOS="linux" ; $env:GOARCH="amd64" ; go build -o easywrite_service  ./main.go
```

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
