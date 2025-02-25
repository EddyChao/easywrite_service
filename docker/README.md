```shell
 cd .. && docker buildx build --platform linux/amd64 -t easywrite-service -f docker/Dockerfile .
```

```shell
 cd .. && docker buildx build --platform linux/arm64 -t easywrite-service -f docker/Dockerfile .
```