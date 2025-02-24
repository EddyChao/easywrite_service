```shell
 cd .. && docker buildx build --platform linux/amd64 -t easy-write-service -f docker/Dockerfile
```

```shell
 cd .. && docker buildx build --platform linux/arm64 -t easy-write-service -f docker/Dockerfile
```