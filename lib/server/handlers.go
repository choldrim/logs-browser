package server

import (
    "fmt"
    "os"
    "io"
    "net/http"
    "log"

    "github.com/choldrim/logs-browser/lib/handle"
)

const (
    uploadDir = "./upload_files"
)

func init() {
    os.Mkdir(uploadDir, 0755)
}

func Index(w http.ResponseWriter, r *http.Request) {
    fmt.Fprintln(w, "hello world!")
}

func ReceiveLog(w http.ResponseWriter, r *http.Request) {
    fin, fHeader, err := r.FormFile("file")
    if err != nil {
        http.Error(w, "" + err.Error(), http.StatusBadRequest)
        log.Println(err)
    }

    filePath := uploadDir + "/" + fHeader.Filename
    fmt.Println("uploaded log: ", filePath)
    fout, err := os.Create(filePath)
    if err != nil {
        log.Printf("error in creating upload file: %v", err)
        http.Error(w, "" + err.Error(), http.StatusBadRequest)
        return
    }

    defer fout.Close()

    _, err = io.Copy(fout, fin)
    if err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        log.Println(err)
    }

    link, err := handleLogs(filePath)
    if err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        log.Println(err)
    }

    fmt.Fprintln(w, link)
}

func handleLogs(path string) (string, error) {
    h, err := handle.New()
    if err != nil {
        return "", fmt.Errorf("error in initing Handler Object: %v", err)
    }

    link, err := h.HandleLog(path)
    if err != nil {
        return "", fmt.Errorf("error in handleLog: %v", err)
    }

    return link, nil
}
