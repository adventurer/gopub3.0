# gopub3.0
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o gopub main.go

daemon.json
{
  "registry-mirrors": ["https://syq2pc87.mirror.aliyuncs.com"]
}