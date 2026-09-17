package service

import (
    "context"
    "dodevops-api/api/db/redis/model"
    "dodevops-api/common/result"
    "encoding/base64"
    "fmt"
    "github.com/gin-gonic/gin"
    redisv8 "github.com/go-redis/redis/v8"
    "strings"
    "time"
    "unicode/utf8"
)

type IRedisService interface {
    TestConn(c *gin.Context, conn *model.Conn)
    Info(c *gin.Context, conn *model.Conn)
    Scan(c *gin.Context, req *model.ScanReq)
    ScanAllDb(c *gin.Context, req *model.ScanReq)
    ListDbs(c *gin.Context, conn *model.Conn)
    GetKey(c *gin.Context, conn *model.Conn, key string)
    DelKey(c *gin.Context, conn *model.Conn, key string)
    RunCmd(c *gin.Context, req *model.CmdReq)
    // Newer design-aligned APIs
    KeyInfo(c *gin.Context, conn *model.Conn, key string)
    TtlKey(c *gin.Context, conn *model.Conn, key string)
    MemoryUsage(c *gin.Context, conn *model.Conn, key string)
    ClusterInfo(c *gin.Context, conn *model.Conn)
}

type RedisServiceImpl struct{}

// RedisConn wraps standalone/cluster/sentinel
type RedisConn struct {
    Mode       string
    Cli        *redisv8.Client
    ClusterCli *redisv8.ClusterClient
}

func (r *RedisConn) Close() {
    if r.Cli != nil { _ = r.Cli.Close() }
    if r.ClusterCli != nil { _ = r.ClusterCli.Close() }
}

func (r *RedisConn) Cmdable() redisv8.Cmdable {
    if r.Mode == "cluster" && r.ClusterCli != nil { return r.ClusterCli }
    return r.Cli
}

// Universal returns a UniversalClient which supports Do.
func (r *RedisConn) Universal() redisv8.UniversalClient {
    if r.Mode == "cluster" && r.ClusterCli != nil { return r.ClusterCli }
    return r.Cli
}

func connect(conn *model.Conn) (*RedisConn, error) {
    mode := strings.ToLower(strings.TrimSpace(conn.Mode))
    if mode == "" { mode = "standalone" }
    if mode == "cluster" {
        addrs := strings.Split(conn.Host, ",")
        cli := redisv8.NewClusterClient(&redisv8.ClusterOptions{
            Addrs:       addrs,
            Password:    conn.Password,
            DialTimeout:  5 * time.Second,
            ReadTimeout:  5 * time.Second,
            WriteTimeout: 5 * time.Second,
        })
        if err := cli.Ping(context.Background()).Err(); err != nil { _ = cli.Close(); return nil, err }
        return &RedisConn{Mode: mode, ClusterCli: cli}, nil
    }
    if mode == "sentinel" {
        // sentinel: Host format masterName=host1:port,host2:port
        parts := strings.Split(conn.Host, "=")
        if len(parts) != 2 {
            return nil, fmt.Errorf("invalid sentinel host format: %s", conn.Host)
        }
        masterName, sentinelAddrs := parts[0], strings.Split(parts[1], ",")
        cli := redisv8.NewFailoverClient(&redisv8.FailoverOptions{
            MasterName:       masterName,
            SentinelAddrs:    sentinelAddrs,
            Password:         conn.RedisNodePassword, // node password
            SentinelPassword: conn.Password,          // sentinel password
            DB:               conn.DB,
            DialTimeout:      5 * time.Second,
            ReadTimeout:      5 * time.Second,
            WriteTimeout:     5 * time.Second,
        })
        if err := cli.Ping(context.Background()).Err(); err != nil { _ = cli.Close(); return nil, err }
        return &RedisConn{Mode: mode, Cli: cli}, nil
    }
    // standalone
    addr := conn.Host
    if !strings.Contains(addr, ":") { addr = fmt.Sprintf("%s:%d", conn.Host, conn.Port) }
    cli := redisv8.NewClient(&redisv8.Options{
        Addr:         addr,
        Password:     conn.Password,
        DB:           conn.DB,
        DialTimeout:  3 * time.Second,
        ReadTimeout:  5 * time.Second,
        WriteTimeout: 5 * time.Second,
    })
    if err := cli.Ping(context.Background()).Err(); err != nil { _ = cli.Close(); return nil, err }
    return &RedisConn{Mode: mode, Cli: cli}, nil
}

func (s RedisServiceImpl) TestConn(c *gin.Context, connInfo *model.Conn) {
    rc, err := connect(connInfo)
    if err != nil {
        result.Failed(c, int(result.ApiCode.FAILED), fmt.Sprintf("连接失败: %v", err))
        return
    }
    defer rc.Close()
    result.Success(c, gin.H{"message": "连接成功"})
}

func (s RedisServiceImpl) Info(c *gin.Context, connInfo *model.Conn) {
    rc, err := connect(connInfo)
    if err != nil { result.Failed(c, int(result.ApiCode.FAILED), err.Error()); return }
    defer rc.Close()

    // Optional: format=map 返回解析后的 section map
    format := c.Query("format")
    out, err := rc.Cmdable().Info(c.Request.Context()).Result()
    if err != nil { result.Failed(c, int(result.ApiCode.FAILED), err.Error()); return }
    if strings.ToLower(format) != "map" {
        result.Success(c, gin.H{"info": out})
        return
    }
    // parse to section map similar to test/api
    lines := strings.Split(out, "\r\n")
    parsed := make(map[string]map[string]string)
    i := 0
    for i < len(lines) {
        if strings.HasPrefix(lines[i], "#") {
            sec := strings.TrimSpace(strings.TrimPrefix(lines[i], "#"))
            i++
            m := map[string]string{}
            for i < len(lines) && strings.Contains(lines[i], ":") {
                kv := strings.SplitN(lines[i], ":", 2)
                if len(kv) == 2 { m[kv[0]] = kv[1] }
                i++
            }
            parsed[sec] = m
        } else {
            i++
        }
    }
    result.Success(c, parsed)
}

func (s RedisServiceImpl) Scan(c *gin.Context, req *model.ScanReq) {
    rc, err := connect(&req.Conn)
    if err != nil { result.Failed(c, int(result.ApiCode.FAILED), err.Error()); return }
    defer rc.Close()
    if req.Count <= 0 { req.Count = 100 }
    match := req.Match
    if match == "" { match = req.Pattern }
    ctx := c.Request.Context()

    // Standalone/sentinel: single cursor
    if rc.Mode != "cluster" {
        keys, cursor, err := rc.Cmdable().Scan(ctx, req.Cursor, match, int64(req.Count)).Result()
        if err != nil { result.Failed(c, int(result.ApiCode.FAILED), err.Error()); return }
        result.Success(c, gin.H{"cursor": cursor, "keys": keys})
        return
    }

    // Cluster: scan each master, track cursor per node
    cursors := map[string]uint64{}
    if req.CursorMap != nil { for k, v := range req.CursorMap { cursors[k] = v } }
    keys := make([]string, 0, req.Count)
    // iterate masters
    _ = rc.ClusterCli.ForEachMaster(ctx, func(ctx context.Context, client *redisv8.Client) error {
        addr := client.Options().Addr
        cur := cursors[addr]
        for {
            ks, next, err := client.Scan(ctx, cur, match, int64(req.Count)).Result()
            if err != nil { break }
            if len(ks) > 0 { keys = append(keys, ks...) }
            cursors[addr] = next
            cur = next
            if int(len(keys)) >= req.Count || next == 0 { break }
        }
        return nil
    })
    result.Success(c, gin.H{"cursor": cursors, "keys": keys})
}

// ScanAllDb 扫描全部DB（仅 standalone / sentinel）
// 返回：
// - cursor: map[string]uint64 键名为 "db:<n>"，表示对应DB的下一游标
// - keys: []string 仅键名列表（兼容旧前端）
// - keysWithDb: []map{"db":int, "key":string} 包含DB信息的新结构
func (s RedisServiceImpl) ScanAllDb(c *gin.Context, req *model.ScanReq) {
    rc, err := connect(&req.Conn)
    if err != nil { result.Failed(c, int(result.ApiCode.FAILED), err.Error()); return }
    defer rc.Close()

    if rc.Mode == "cluster" {
        result.Failed(c, int(result.ApiCode.FAILED), "集群模式不支持多DB，默认为DB0")
        return
    }

    if req.Count <= 0 { req.Count = 100 }
    match := req.Match
    if match == "" { match = req.Pattern }
    ctx := c.Request.Context()

    // 通过 INFO keyspace 解析出实际存在的 db 列表
    info, err := rc.Cmdable().Info(ctx, "keyspace").Result()
    if err != nil { result.Failed(c, int(result.ApiCode.FAILED), err.Error()); return }
    dbList := make([]int, 0, 16)
    for _, line := range strings.Split(info, "\n") {
        line = strings.TrimSpace(line)
        if !strings.HasPrefix(line, "db") { continue }
        // line: db0:keys=1,expires=0,avg_ttl=0
        var n int
        _, _ = fmt.Sscanf(strings.Split(line, ":")[0], "db%d", &n)
        if n >= 0 { dbList = append(dbList, n) }
    }
    // 若未能从 keyspace 获取到，降级使用 0..15 常见范围
    if len(dbList) == 0 {
        for i := 0; i < 16; i++ { dbList = append(dbList, i) }
    }

    // 聚合各DB的扫描
    totalKeys := make([]string, 0, req.Count)
    keysWithDb := make([]map[string]interface{}, 0, req.Count)
    cursors := map[string]uint64{}
    if req.CursorMap != nil {
        for k, v := range req.CursorMap { cursors[k] = v }
    }

    for _, db := range dbList {
        if len(totalKeys) >= req.Count { break }
        // SELECT DB
        if err := rc.Universal().Do(ctx, "SELECT", db).Err(); err != nil { continue }
        key := fmt.Sprintf("db:%d", db)
        cur := cursors[key]
        ks, next, err := rc.Cmdable().Scan(ctx, cur, match, int64(req.Count-len(totalKeys))).Result()
        if err != nil { continue }
        cursors[key] = next
        if len(ks) > 0 {
            totalKeys = append(totalKeys, ks...)
            for _, k := range ks { keysWithDb = append(keysWithDb, map[string]interface{}{"db": db, "key": k}) }
        }
    }

    result.Success(c, gin.H{"cursor": cursors, "keys": totalKeys, "keysWithDb": keysWithDb})
}

func (s RedisServiceImpl) GetKey(c *gin.Context, connInfo *model.Conn, key string) {
    rc, err := connect(connInfo)
    if err != nil { result.Failed(c, int(result.ApiCode.FAILED), err.Error()); return }
    defer rc.Close()
    ctx := context.Background()
    typ, err := rc.Cmdable().Type(ctx, key).Result()
    if err != nil { result.Failed(c, int(result.ApiCode.FAILED), err.Error()); return }
    var val interface{}
    switch typ {
    case "string":
        val, err = rc.Cmdable().Get(ctx, key).Result()
    case "hash":
        val, err = rc.Cmdable().HGetAll(ctx, key).Result()
    case "list":
        val, err = rc.Cmdable().LRange(ctx, key, 0, -1).Result()
    case "set":
        val, err = rc.Cmdable().SMembers(ctx, key).Result()
    case "zset":
        val, err = rc.Cmdable().ZRangeWithScores(ctx, key, 0, -1).Result()
    default:
        val = nil
    }
    if err != nil { result.Failed(c, int(result.ApiCode.FAILED), err.Error()); return }
    ttl, _ := rc.Cmdable().TTL(ctx, key).Result()
    result.Success(c, gin.H{"type": typ, "value": val, "ttl": ttl.Seconds()})
}

func (s RedisServiceImpl) DelKey(c *gin.Context, connInfo *model.Conn, key string) {
    rc, err := connect(connInfo)
    if err != nil { result.Failed(c, int(result.ApiCode.FAILED), err.Error()); return }
    defer rc.Close()
    n, err := rc.Cmdable().Del(c.Request.Context(), key).Result()
    if err != nil { result.Failed(c, int(result.ApiCode.FAILED), err.Error()); return }
    result.Success(c, gin.H{"deleted": n})
}

func (s RedisServiceImpl) RunCmd(c *gin.Context, req *model.CmdReq) {
    rc, err := connect(&req.Conn)
    if err != nil { result.Failed(c, int(result.ApiCode.FAILED), err.Error()); return }
    defer rc.Close()
    args := make([]interface{}, 0, len(req.Args))
    for _, a := range req.Args { args = append(args, a) }
    // go-redis v8 的 Do(ctx, args...) 第一个参数应为命令名
    doArgs := make([]interface{}, 0, 1+len(args))
    doArgs = append(doArgs, req.Command)
    doArgs = append(doArgs, args...)
    res := rc.Universal().Do(c.Request.Context(), doArgs...)
    if err := res.Err(); err != nil { result.Failed(c, int(result.ApiCode.FAILED), err.Error()); return }
    v, _ := res.Result()

    // 处理二进制数据：转换为 base64 编码
    if processed := processBinaryResult(v); processed != nil {
        result.Success(c, gin.H{"result": processed})
        return
    }

    result.Success(c, gin.H{"result": v})
}

// processBinaryResult 处理 Redis 命令返回的二进制数据
// 将二进制字符串/字节数组转换为 base64 编码，避免 JSON 序列化时出现乱码
func processBinaryResult(v interface{}) interface{} {
    switch val := v.(type) {
    case string:
        // 检查字符串是否包含非 UTF-8 字符（可能是二进制数据）
        if !utf8.ValidString(val) {
            return map[string]string{
                "_type":   "binary",
                "_base64": base64.StdEncoding.EncodeToString([]byte(val)),
            }
        }
    case []byte:
        // 字节数组直接转 base64
        return map[string]string{
            "_type":   "binary",
            "_base64": base64.StdEncoding.EncodeToString(val),
        }
    case []interface{}:
        // 递归处理数组元素
        result := make([]interface{}, len(val))
        for i, item := range val {
            if processed := processBinaryResult(item); processed != nil {
                result[i] = processed
            } else {
                result[i] = item
            }
        }
        return result
    }
    return nil
}

func RedisService() IRedisService { return RedisServiceImpl{} }

// Newer design-aligned implementations
func (s RedisServiceImpl) KeyInfo(c *gin.Context, connInfo *model.Conn, key string) {
    rc, err := connect(connInfo)
    if err != nil { result.Failed(c, int(result.ApiCode.FAILED), err.Error()); return }
    defer rc.Close()
    ctx := context.Background()
    typ, err := rc.Cmdable().Type(ctx, key).Result()
    if err != nil { result.Failed(c, int(result.ApiCode.FAILED), err.Error()); return }
    ttl, err := rc.Cmdable().TTL(ctx, key).Result()
    if err != nil { result.Failed(c, int(result.ApiCode.FAILED), err.Error()); return }
    ttlSec := -1.0
    if ttl > 0 { ttlSec = ttl.Seconds() }
    result.Success(c, gin.H{"key": key, "type": typ, "ttl": ttlSec})
}

func (s RedisServiceImpl) TtlKey(c *gin.Context, connInfo *model.Conn, key string) {
    rc, err := connect(connInfo)
    if err != nil { result.Failed(c, int(result.ApiCode.FAILED), err.Error()); return }
    defer rc.Close()
    ttl, err := rc.Cmdable().TTL(context.Background(), key).Result()
    if err != nil { result.Failed(c, int(result.ApiCode.FAILED), err.Error()); return }
    if ttl == -1 { result.Success(c, -1); return }
    result.Success(c, ttl.Seconds())
}

func (s RedisServiceImpl) MemoryUsage(c *gin.Context, connInfo *model.Conn, key string) {
    rc, err := connect(connInfo)
    if err != nil { result.Failed(c, int(result.ApiCode.FAILED), err.Error()); return }
    defer rc.Close()
    v := rc.Cmdable().MemoryUsage(context.Background(), key).Val()
    result.Success(c, v)
}

func (s RedisServiceImpl) ClusterInfo(c *gin.Context, connInfo *model.Conn) {
    rc, err := connect(connInfo)
    if err != nil { result.Failed(c, int(result.ApiCode.FAILED), err.Error()); return }
    defer rc.Close()
    if rc.Mode != "cluster" || rc.ClusterCli == nil {
        result.Failed(c, int(result.ApiCode.FAILED), "非集群模式")
        return
    }
    info, _ := rc.ClusterCli.ClusterInfo(context.Background()).Result()
    nodesStr, _ := rc.ClusterCli.ClusterNodes(context.Background()).Result()
    nodesRes := make([]map[string]string, 0)
    for _, line := range strings.Split(nodesStr, "\n") {
        if line == "" { continue }
        segs := strings.Split(line, " ")
        node := map[string]string{"nodeId": segs[0]}
        if len(segs) > 1 { node["ip"] = segs[1] }
        if len(segs) > 2 { node["flags"] = segs[2] }
        if len(segs) > 3 { node["masterSlaveRelation"] = segs[3] }
        if len(segs) > 4 { node["pingSent"] = segs[4] }
        if len(segs) > 5 { node["pongRecv"] = segs[5] }
        if len(segs) > 6 { node["configEpoch"] = segs[6] }
        if len(segs) > 7 { node["linkState"] = segs[7] }
        if len(segs) > 8 { node["slot"] = segs[8] }
        nodesRes = append(nodesRes, node)
    }
    result.Success(c, gin.H{"clusterInfo": info, "clusterNodes": nodesRes})
}

// ListDbs 返回可用DB列表。
// - cluster: 仅返回 [0]
// - standalone/sentinel: 解析 INFO keyspace 获取已有DB；若为空，则回退 0..15
// 返回结构：{"dbs": [0,1,...], "dbNames": ["db0","db1",...]}，并包含每个DB的keys统计（若可获取）。
func (s RedisServiceImpl) ListDbs(c *gin.Context, connInfo *model.Conn) {
    rc, err := connect(connInfo)
    if err != nil { result.Failed(c, int(result.ApiCode.FAILED), err.Error()); return }
    defer rc.Close()
    ctx := c.Request.Context()

    if rc.Mode == "cluster" {
        result.Success(c, gin.H{"dbs": []int{0}, "dbNames": []string{"db0"}})
        return
    }

    info, err := rc.Cmdable().Info(ctx, "keyspace").Result()
    dbList := make([]int, 0, 16)
    dbStats := map[int]map[string]int{}
    if err == nil {
        for _, line := range strings.Split(info, "\n") {
            line = strings.TrimSpace(line)
            if !strings.HasPrefix(line, "db") { continue }
            // eg: db0:keys=1,expires=0,avg_ttl=0
            parts := strings.SplitN(line, ":", 2)
            if len(parts) != 2 { continue }
            var n int
            _, _ = fmt.Sscanf(parts[0], "db%d", &n)
            if n < 0 { continue }
            dbList = append(dbList, n)
            stat := map[string]int{}
            for _, kv := range strings.Split(parts[1], ",") {
                kv = strings.TrimSpace(kv)
                seg := strings.SplitN(kv, "=", 2)
                if len(seg) != 2 { continue }
                var vi int
                _, _ = fmt.Sscanf(seg[1], "%d", &vi)
                stat[seg[0]] = vi
            }
            dbStats[n] = stat
        }
    }
    // 获取配置的数据库数量
    dbCount := 16
    if rc.Mode != "cluster" {
        if cfg, e := rc.Cmdable().ConfigGet(ctx, "databases").Result(); e == nil && len(cfg) == 2 {
            if s, ok := cfg[1].(string); ok {
                var n int
                if _, se := fmt.Sscanf(s, "%d", &n); se == nil && n > 0 { dbCount = n }
            }
        }
    }
    // 用配置数量覆盖列表范围（保证显示 db0..dbN-1），同时保留已有的统计
    seen := map[int]struct{}{}
    for _, d := range dbList { seen[d] = struct{}{} }
    full := make([]int, 0, dbCount)
    for i := 0; i < dbCount; i++ { full = append(full, i); if _, ok := seen[i]; !ok { /* no-op, stats为空 */ } }
    dbList = full
    names := make([]string, 0, len(dbList))
    stats := make([]map[string]interface{}, 0, len(dbList))
    for _, d := range dbList {
        names = append(names, fmt.Sprintf("db%d", d))
        s := map[string]interface{}{"db": d}
        if st, ok := dbStats[d]; ok { for k, v := range st { s[k] = v } }
        stats = append(stats, s)
    }
    result.Success(c, gin.H{"dbs": dbList, "dbNames": names, "stats": stats})
}
