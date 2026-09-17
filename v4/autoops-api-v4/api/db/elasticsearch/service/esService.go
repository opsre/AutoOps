package service

import (
    "bytes"
    "crypto/tls"
    "dodevops-api/api/db/elasticsearch/dao"
    "dodevops-api/api/db/elasticsearch/model"
    "dodevops-api/common/result"
    "encoding/base64"
    "io"
    "strconv"
    "net/http"
    "strings"
    "time"

    "github.com/gin-gonic/gin"
)

type IEsService interface {
    TestConn(c *gin.Context, inst *model.EsInstance)
    Proxy(c *gin.Context, inst *model.EsInstance, path string)
    ProxyById(c *gin.Context, id uint, path string)
}

type EsServiceImpl struct{}

func EsService() IEsService { return EsServiceImpl{} }

func buildBaseURL(inst *model.EsInstance) string {
    proto := inst.Protocol
    if proto == "" { proto = "http" }
    port := inst.Port
    if port == 0 { port = 9200 }
    return proto + "://" + inst.Host + ":" + strconv.Itoa(port)
}

func buildHttpClient(inst *model.EsInstance) *http.Client {
    if strings.ToLower(inst.Protocol) == "https" {
        tr := &http.Transport{TLSClientConfig: &tls.Config{InsecureSkipVerify: true}}
        return &http.Client{Transport: tr, Timeout: 10 * time.Second}
    }
    return &http.Client{Timeout: 10 * time.Second}
}

func (EsServiceImpl) TestConn(c *gin.Context, inst *model.EsInstance) {
    client := buildHttpClient(inst)
    req, _ := http.NewRequest(http.MethodGet, buildBaseURL(inst)+"/", nil)
    if inst.Username != "" && inst.Password != "" {
        auth := base64.StdEncoding.EncodeToString([]byte(inst.Username + ":" + inst.Password))
        req.Header.Set("Authorization", "Basic "+auth)
    }
    resp, err := client.Do(req)
    if err != nil { result.Failed(c, int(result.ApiCode.FAILED), err.Error()); return }
    defer resp.Body.Close()
    body, _ := io.ReadAll(resp.Body)
    if resp.StatusCode >= 200 && resp.StatusCode < 300 {
        c.Data(resp.StatusCode, "application/json", body)
        return
    }
    result.Failed(c, int(result.ApiCode.FAILED), string(body))
}

func (EsServiceImpl) Proxy(c *gin.Context, inst *model.EsInstance, p string) {
    // cleanup query: remove id/path params if present
    r := c.Request
    // clone query
    q := r.URL.Query()
    q.Del("id"); q.Del("path")
    // construct target URL
    base := buildBaseURL(inst)
    // ensure path leading slash handling
    if !strings.HasPrefix(p, "/") { p = "/" + p }
    target := base + p
    // rebuild raw query
    if len(q) > 0 {
        target += "?" + q.Encode()
    }

    // read body for forwarding
    var body io.Reader
    if r.Body != nil {
        data, _ := io.ReadAll(r.Body)
        r.Body.Close()
        body = bytes.NewReader(data)
    }

    req, _ := http.NewRequest(r.Method, target, body)
    // copy content-type and others minimal
    for k, vs := range r.Header {
        if strings.ToLower(k) == "host" { continue }
        for _, v := range vs { req.Header.Add(k, v) }
    }
    if inst.Username != "" && inst.Password != "" {
        auth := base64.StdEncoding.EncodeToString([]byte(inst.Username + ":" + inst.Password))
        req.Header.Set("Authorization", "Basic "+auth)
    }

    client := buildHttpClient(inst)
    resp, err := client.Do(req)
    if err != nil { result.Failed(c, int(result.ApiCode.FAILED), err.Error()); return }
    defer resp.Body.Close()
    // stream back response - copy headers first, then set status, then write body
    for k, vs := range resp.Header { for _, v := range vs { c.Writer.Header().Add(k, v) } }
    c.Writer.WriteHeader(resp.StatusCode)
    io.Copy(c.Writer, resp.Body)
}

func (EsServiceImpl) ProxyById(c *gin.Context, id uint, p string) {
    inst, err := dao.GetEsInstanceById(id)
    if err != nil { result.Failed(c, int(result.ApiCode.FAILED), "实例不存在"); return }
    if inst.Password == "******" { inst.Password = "" }
    EsServiceImpl{}.Proxy(c, &inst, p)
}
