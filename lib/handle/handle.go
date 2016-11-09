package handle

import (
    "fmt"
    "io/ioutil"
    "log"
    "path/filepath"
    libPath "path"
    "os"
    "strings"

    "github.com/mholt/archiver"

    "github.com/choldrim/logs-browser/lib/seafile"
)

type Handler struct {
    s *seafile.Seafile
    workingDir string
}

func New() (*Handler, error) {
    s, err := seafile.New()
    if err != nil {
        fmt.Errorf("errror creating init seafile: %v\n", err)
    }

    return &Handler{s: s}, nil
}

func extract(file string, folder string) error {
    err := archiver.TarGz.Open(file, folder)
    return err
}


func (h *Handler)HandleLog(file string) (string, error) {
    // create temp dir
    tempDir, err := ioutil.TempDir("", "logs_browser_")
    if err != nil {
        return "", fmt.Errorf("error creating temp dir: %v", err)
    }

    //defer os.RemoveAll(tempDir)
    h.workingDir = tempDir

    err = extract(file, h.workingDir)
    if err != nil {
        return "", err
    }

    logName := libPath.Base(file)
    nameList := strings.Split(logName, ".tar.gz")
    if len(nameList) > 0 {
        logName = nameList[0]
    }

    link, err := h.walkLogs(logName)
    if err != nil {
        return "", err
    }

    log.Printf("share link: %s\n", link)
    return link, nil
}


func (h *Handler)walkLogs(logName string) (string, error) {
    walkFunc := func (path string, info os.FileInfo, err error) error {
        if path == h.workingDir {  // skip self
            return nil
        }

        pathList := strings.Split(path, h.workingDir + "/")
        if len(pathList) == 0 {
            return fmt.Errorf("split path error with: %s", path)
        }

        path = pathList[1]
        logBase := "/" + logName
        if info.IsDir() {
            log.Printf("mkdir folder: %s/%s\n", logBase, path)
            err := h.s.NewFolder(path, logBase)
            if err != nil {
                return err
            }
        } else {
            truePath := h.workingDir + "/" + path
            fileName := libPath.Base(path)
            pathList := strings.Split(path, fileName)
            fileParent := ""
            if len(pathList) == 0 {
                return fmt.Errorf("split path error with: %s", path)
            } else {
                // for special path like: "deepin_music/~/.config/deepin-music-player/config"
                dir := strings.Join(pathList[:len(pathList)-1], fileName)
                if strings.HasSuffix(dir, "/") {
                    dir = dir[:len(dir) - 1]
                    dir = "/" + dir
                }
                fileParent =  logBase + dir
            }

            log.Printf("upload file: %s/%s\n", fileParent, fileName)
            err := h.s.NewFile(fileName, fileParent, truePath)
            if err != nil {
                return err
            }
        }

        return nil
    }

    // delete folder is exist
    err := h.s.DeleteIfExist("/" + logName)
    if err != nil {
        return "", err
    }

    // create log base dir
    err = h.s.NewFolder(logName, "/")
    if err != nil {
        return "", err
    }

    // upload log
    err = filepath.Walk(h.workingDir, walkFunc)
    if err != nil {
        return "", err
    }

    // Share Link
    link, err := h.s.ShareLink("/" + logName)
    if err != nil {
        return "", err
    }

    return link, nil
}
