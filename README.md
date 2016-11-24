# logs browser
## a logs handling server for deepin feedback logs

## Depends
- an available seafile server

## Usage (shell)
### config
```shell
cp config-example.ini config.ini

# fill config.ini basic on your environment
vim config.ini
```

### run
```shell
go run main.go
```

## Usage (Docker)
### build image
```shell
docker build -t logs-browser-server .
```

### config
```shell
cp config-example.ini config.ini

# fill config.ini basic on your environment
vim config.ini
```

### run
```shell
docker run -t $PWD/config.ini:/go/src/app/config.ini -p 8090:8090 --restart=always logs-browser-server
```

## Test
```shell
# for example
curl http://localhost:8090/api/v1/log -XPOST -F "file=@deepin-feedback-cli-Deepin-15.3-all-20161124-151557.tar.gz"
```
