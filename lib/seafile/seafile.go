package seafile

import (
    "fmt"
    "log"
    "os"
    "strings"

    "github.com/mozillazg/request"
    "gopkg.in/ini.v1"
)

const (
    confFile = "./config.ini"
)

var (
    serverAPIV2 string
    serverAPIV2_1 string
    repoID string
    userName string
    password string
)


type Seafile struct {
    token string
    uploadLink string
}

func init() {
    cfg, err := ini.Load(confFile)
    if err != nil {
        log.Fatalln(err)
    }

    apiv2, err := cfg.Section("seafile").GetKey("apiv2")
    if err != nil {
        log.Fatalln(err)
    }

    serverAPIV2 = apiv2.String()

    apiv2_1, err := cfg.Section("seafile").GetKey("apiv2_1")
    if err != nil {
        log.Fatalln(err)
    }

    serverAPIV2_1 = apiv2_1.String()

    id, err := cfg.Section("seafile").GetKey("repo_id")
    if err != nil {
        log.Fatalln(err)
    }

    repoID = id.String()

    u, err := cfg.Section("seafile").GetKey("username")
    if err != nil {
        log.Fatalln(err)
    }

    userName = u.String()

    p, err := cfg.Section("seafile").GetKey("password")
    if err != nil {
        log.Fatalln(err)
    }

    password = p.String()
}


func New() (*Seafile, error) {
    token, err := GetToken(userName, password)
    if err != nil {
        return nil, fmt.Errorf("failed to get token (%s)", err)
    }

    link, err := GetUploadLink(token)
    if err != nil {
        return nil, fmt.Errorf("failed to get upload link (%s)", err)
    }

    s := &Seafile{token: token, uploadLink: link}
    return s, nil
}


func GetToken(userName string, password string) (string, error) {
    url := serverAPIV2 + "/auth-token/"
    req := request.NewRequest(nil)
    req.Data = map[string]string {
        "username": userName,
        "password": password,
    }

    resp, err := req.Post(url)
    if err != nil {
        return "", err
    }

    if ! resp.Ok() {
        return "", fmt.Errorf(resp.Text())
    }

    j, err := resp.Json()
    if err != nil {
        return "", err
    }

    tokenJson := j.Get("token")
    token, err := tokenJson.String()
    if err != nil {
        return "", err
    }

    return token, nil
}


func (s *Seafile)NewFolder(name string, parent string) (error) {
    req := request.NewRequest(nil)
    req.Params = map[string]string{
        "p": parent + "/" + name,
    }

    req.Data = map[string]string{
        "operation": "mkdir",
    }

    req.Headers = map[string]string{
        "authorization": "Token " + s.token,
    }

    url := serverAPIV2 + "/repos/" + repoID + "/dir/"
    resp, err := req.Post(url)
    if err != nil {
        return err
    }

    if ! resp.Ok() {
        return fmt.Errorf(resp.Text())
    }

    return nil
}


func GetUploadLink(token string) (string, error) {
    req := request.NewRequest(nil)
    req.Params = map[string]string{
        "p": "/",
    }

    req.Headers = map[string]string{
        "authorization": "Token " + token,
    }

    url := serverAPIV2 + "/repos/" + repoID + "/upload-link/"
    resp, err := req.Get(url)
    if err != nil {
        return "", err
    }

    if ! resp.Ok() {
        return "", fmt.Errorf(resp.Text())
    }

    link, err := resp.Text()
    if err != nil {
        return "", err
    }

    link = strings.Replace(link, "\"", "", -1)

    return link, nil
}


func (s *Seafile)NewFile(name string, parent string, truePath string) (error) {
    req := request.NewRequest(nil)
    form := map[string]string{
        "parent_dir": parent,
    }

    req.Headers = map[string]string{
        "authorization": "Token " + s.token,
    }

    f, err := os.Open(truePath)
    if err != nil {
        return fmt.Errorf("open file error: %v", err)
    }

    defer f.Close()

    req.Files = []request.FileField{
        request.FileField{"file", name, f},
    }

    resp, err := req.PostForm(s.uploadLink, form)
    if err != nil {
        return fmt.Errorf("request error on creating file: %v", err)
    }

    if ! resp.Ok() {
        return fmt.Errorf(resp.Text())
    }

    return nil
}


func (s *Seafile)DeleteIfExist(path string) error {
    req := request.NewRequest(nil)
    url := serverAPIV2 + "/repos/" + repoID + "/dir/"
    req.Params = map[string]string{
        "p": path,
    }

    req.Headers = map[string]string{
        "authorization": "Token " + s.token,
    }

    resp, err := req.Delete(url)
    if err != nil {
        return fmt.Errorf("request error on delete dir: %v", err)
    }

    if ! resp.Ok() {
        return fmt.Errorf(resp.Text())
    }

    return nil
}


func (s *Seafile)ShareLink(path string) (string, error) {
    req := request.NewRequest(nil)
    url := serverAPIV2_1 + "/share-links/"

    req.Headers = map[string]string{
        "authorization": "Token " + s.token,
    }

    req.Data = map[string]string{
        "repo_id": repoID,
        "path": path,
    }

    resp, err := req.Post(url)
    if err != nil {
        return "", fmt.Errorf("error sharing log: ", err)
    }

    if ! resp.Ok() {
        return "", fmt.Errorf(resp.Text())
    }

    data, err := resp.Json()
    if err != nil {
        return "", err
    }

    linkData := data.Get("link")
    link, err := linkData.String()
    if err != nil {
        return "", err
    }

    return link, nil
}
