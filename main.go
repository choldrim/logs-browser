package main


import (
    "log"
    "net/http"

    "gopkg.in/ini.v1"

    "github.com/choldrim/logs-browser/lib/server"
)

const (
    confFile = "./config.ini"
)

func getListenAddr() (string, error){
    cfg, err := ini.Load(confFile)
    if err != nil {
        return "", err
    }

    p, err := cfg.Section("server").GetKey("port")
    if err != nil {
        return "", err
    }

    port := p.String()
    addr := "0.0.0.0:" + port

    return addr, nil
}

func main() {
    addr, err := getListenAddr()
    if err != nil {
        log.Fatalln(err)
    }

    router := server.NewRouter()
    log.Fatalln(http.ListenAndServe(addr, router))
}
