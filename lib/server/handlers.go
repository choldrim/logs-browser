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
        panic(err)
    }

    filePath := uploadDir + "/" + fHeader.Filename
    fmt.Println("uploaded log: ", filePath)
    fout, err := os.Create(filePath)
    if err != nil {
        log.Printf("error in creating upload file: %v", err)
        panic(err)
    }

    defer fout.Close()

    _, err = io.Copy(fout, fin)
    if err != nil {
        http.Error(w, "error saving file: " + err.Error(), http.StatusBadRequest)
        panic(err)
    }

    link, err := handleLogs(filePath)
    if err != nil {
        http.Error(w, "error handling log: " + err.Error(), http.StatusBadRequest)
        panic(err)
    }

    fmt.Fprintln(w, link)
}

func handleLogs(path string) (string, error) {
    h, err := handle.New()
    if err != nil {
        return "", fmt.Errorf("error in initing handle obj: %v", err)
    }

    link, err := h.HandleLog(path)
    if err != nil {
        return "", fmt.Errorf("error in handling log: %v", err)
    }

    return link, nil
}
