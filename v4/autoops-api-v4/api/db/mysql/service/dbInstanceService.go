package service

import (
    "context"
    "dodevops-api/api/db/mysql/dao"
    "dodevops-api/api/db/mysql/model"
    "dodevops-api/common/result"
    "dodevops-api/common/util"
    "dodevops-api/pkg/dbm"
    "dodevops-api/pkg/dbm/dbi"
    "dodevops-api/pkg/redis"
    "encoding/json"
    "fmt"
    "github.com/gin-gonic/gin"
    "github.com/google/uuid"
    "net"
    "regexp"
    "strings"
    "time"
)

// IDbInstanceService 数据库实例服务接口
type IDbInstanceService interface {
	GetDbInstanceList(c *gin.Context, query *model.DbInstanceQuery)
	GetDbInstanceById(c *gin.Context, id uint)
	CreateDbInstance(c *gin.Context, dto *model.DbInstanceDto)
	UpdateDbInstance(c *gin.Context, dto *model.DbInstanceDto)
	DeleteDbInstance(c *gin.Context, id uint)
	TestConnection(c *gin.Context, dto *model.TestConnDto)
	GetDatabases(c *gin.Context, id uint)
	GetDbConn(ctx context.Context, instanceId uint, database string) (*dbi.DbConn, error)
	CloseDbConn(instanceId uint, database string)
}

type DbInstanceServiceImpl struct{}

// GetDbInstanceList 获取实例列表
func (s DbInstanceServiceImpl) GetDbInstanceList(c *gin.Context, query *model.DbInstanceQuery) {
	instances, total, err := dao.GetDbInstanceList(query)
	if err != nil {
		result.Failed(c, int(result.ApiCode.FAILED), err.Error())
		return
	}

	// 脱敏密码
	for i := range instances {
		instances[i].Password = "******"
	}

	data := map[string]interface{}{
		"list":     instances,
		"total":    total,
		"pageNum":  query.PageNum,
		"pageSize": query.PageSize,
	}
	result.Success(c, data)
}

// GetDbInstanceById 根据ID获取实例
func (s DbInstanceServiceImpl) GetDbInstanceById(c *gin.Context, id uint) {
	instance, err := dao.GetDbInstanceById(id)
	if err != nil {
		result.Failed(c, int(result.ApiCode.FAILED), err.Error())
		return
	}

	// 脱敏密码
	instance.Password = "******"

	result.Success(c, instance)
}

// CreateDbInstance 创建实例
func (s DbInstanceServiceImpl) CreateDbInstance(c *gin.Context, dto *model.DbInstanceDto) {
	// 检查编码是否已存在
	if dto.Code != "" && dao.CheckInstanceExists(dto.Code, 0) {
		result.Failed(c, int(result.ApiCode.FAILED), "实例编码已存在")
		return
	}

	// 生成编码
	if dto.Code == "" {
		dto.Code = uuid.New().String()
	}

	// 默认网络类型
	if dto.Network == "" {
		dto.Network = "tcp"
	}

	now := time.Now()
	instance := &model.DbInstance{
		Code:               dto.Code,
		Name:               dto.Name,
		Type:               dto.Type,
		Host:               dto.Host,
		Port:               dto.Port,
		Network:            dto.Network,
		Params:             dto.Params,
		Username:           dto.Username,
		Password:           dto.Password, // TODO: 加密存储
		Remark:             dto.Remark,
		SshTunnelMachineId: dto.SshTunnelMachineId,
		Status:             1,
		CreateTime:         util.HTime{Time: now},
		UpdateTime:         util.HTime{Time: now},
		Creator:            c.GetString("username"),
		CreatorId:          uint(c.GetInt("userId")),
	}

	if err := dao.CreateDbInstance(instance); err != nil {
		result.Failed(c, int(result.ApiCode.FAILED), err.Error())
		return
	}

	result.Success(c, instance.ID)
}

// UpdateDbInstance 更新实例
func (s DbInstanceServiceImpl) UpdateDbInstance(c *gin.Context, dto *model.DbInstanceDto) {
	// 检查实例是否存在
	instance, err := dao.GetDbInstanceById(dto.ID)
	if err != nil {
		result.Failed(c, int(result.ApiCode.FAILED), "实例不存在")
		return
	}

	// 检查编码是否重复
	if dto.Code != "" && dto.Code != instance.Code && dao.CheckInstanceExists(dto.Code, dto.ID) {
		result.Failed(c, int(result.ApiCode.FAILED), "实例编码已存在")
		return
	}

	// 更新字段
	instance.Name = dto.Name
	instance.Type = dto.Type
	instance.Host = dto.Host
	instance.Port = dto.Port
	if dto.Network != "" {
		instance.Network = dto.Network
	}
	instance.Params = dto.Params
	instance.Username = dto.Username
	if dto.Password != "" && dto.Password != "******" {
		instance.Password = dto.Password // TODO: 加密存储
	}
	instance.Remark = dto.Remark
	instance.SshTunnelMachineId = dto.SshTunnelMachineId
	instance.UpdateTime = util.HTime{Time: time.Now()}
	instance.Modifier = c.GetString("username")
	instance.ModifierId = uint(c.GetInt("userId"))

	if err := dao.UpdateDbInstance(&instance); err != nil {
		result.Failed(c, int(result.ApiCode.FAILED), err.Error())
		return
	}

	// 清理实例缓存
	if redis.RedisDb != nil {
		_ = redis.RedisDb.Del(context.Background(), fmt.Sprintf("db:instance:%d", dto.ID)).Err()
	}

	// 关闭旧连接
	s.CloseDbConn(dto.ID, "")

	result.Success(c, nil)
}

// DeleteDbInstance 删除实例
func (s DbInstanceServiceImpl) DeleteDbInstance(c *gin.Context, id uint) {
	// 检查是否有关联的数据库
	dbs, _ := dao.GetDbsByInstanceId(id)
	if len(dbs) > 0 {
		result.Failed(c, int(result.ApiCode.FAILED), "该实例下还有数据库，无法删除")
		return
	}

	if err := dao.DeleteDbInstance(id); err != nil {
		result.Failed(c, int(result.ApiCode.FAILED), err.Error())
		return
	}

	// 清理实例缓存
	if redis.RedisDb != nil {
		_ = redis.RedisDb.Del(context.Background(), fmt.Sprintf("db:instance:%d", id)).Err()
	}

	// 关闭连接
	s.CloseDbConn(id, "")

	result.Success(c, nil)
}

// TestConnection 测试连接
func (s DbInstanceServiceImpl) TestConnection(c *gin.Context, dto *model.TestConnDto) {
    // 若提供 instanceId，则优先使用保存的实例信息，避免前端填错或加密不一致
    var dbInfo *dbi.DbInfo
    var serverHost string
    var serverPort int
    if dto.InstanceId > 0 {
        // 读取实例
        inst, err := dao.GetDbInstanceById(dto.InstanceId)
        if err != nil { result.Failed(c, int(result.ApiCode.FAILED), "实例不存在"); return }
        if dto.Network == "" { dto.Network = inst.Network }
        // 使用保存的凭据，db 使用传入 database（可为空仅测连通）
        dbInfo = &dbi.DbInfo{
            InstanceId: uint64(inst.ID),
            Type:       dbi.DbType(inst.Type),
            Host:       inst.Host,
            Port:       inst.Port,
            Username:   inst.Username,
            Password:   inst.Password,
            Database:   dto.Database,
            Network:    firstNonEmpty(dto.Network, inst.Network, "tcp"),
            Params:     firstNonEmpty(dto.Params, inst.Params, ""),
        }
        serverHost = inst.Host
        serverPort = inst.Port
    } else {
        // 设置默认网络类型
        network := dto.Network
        if network == "" { network = "tcp" }
        // 先进行基础可达性检测（快速失败）
        ctx, cancel := context.WithTimeout(c.Request.Context(), 3*time.Second)
        defer cancel()
        d := net.Dialer{}
        if connTest, err := d.DialContext(ctx, "tcp", fmt.Sprintf("%s:%d", dto.Host, dto.Port)); err != nil {
            result.Failed(c, int(result.ApiCode.FAILED), fmt.Sprintf("无法连接到 %s:%d (TCP): %v", dto.Host, dto.Port, err))
            return
        } else { _ = connTest.Close() }
        dbInfo = &dbi.DbInfo{
            Type:     dbi.DbType(dto.Type),
            Host:     dto.Host,
            Port:     dto.Port,
            Username: dto.Username,
            Password: dto.Password,
            Database: dto.Database,
            Network:  network,
            Params:   dto.Params,
        }
        serverHost = dto.Host
        serverPort = dto.Port
    }

    // 测试连接（带超时）
    conn, err := dbm.Conn(c.Request.Context(), dbInfo)
    if err != nil {
        // 识别常见错误并给出指引
        msg := err.Error()
        if strings.Contains(msg, "Error 1045") || strings.Contains(strings.ToLower(msg), "access denied") {
            srcIP := extractDeniedSourceIP(msg)
            if srcIP == "" { srcIP = c.ClientIP() }
            msg = fmt.Sprintf("认证失败(1045): 目标 %s:%d | 客户端IP(以MySQL为准): %s. 请检查账号/密码，或为该IP授权", serverHost, serverPort, srcIP)
        } else if strings.Contains(msg, "Error 1049") || strings.Contains(strings.ToLower(msg), "unknown database") {
            msg = fmt.Sprintf("数据库不存在(1049): 目标 %s:%d | 请确认 database 参数是否正确", serverHost, serverPort)
        } else {
            msg = fmt.Sprintf("目标 %s:%d | %s", serverHost, serverPort, msg)
        }
        result.Failed(c, int(result.ApiCode.FAILED), fmt.Sprintf("连接失败: %s", msg))
        return
    }
    defer conn.Close()

    // 获取数据库版本
    metadata := conn.GetMetadata()
    server, err := metadata.GetDbServer()
	if err != nil {
		result.Failed(c, int(result.ApiCode.FAILED), fmt.Sprintf("获取数据库信息失败: %v", err))
		return
	}

	result.Success(c, map[string]interface{}{
		"version": server.Version,
		"message": "连接成功",
	})
}

func firstNonEmpty(vals ...string) string {
    for _, v := range vals { if strings.TrimSpace(v) != "" { return v } }
    return ""
}

var accessDeniedRe = regexp.MustCompile(`Access denied for user '\S+'@'([^']+)'`)

func extractDeniedSourceIP(msg string) string {
    m := accessDeniedRe.FindStringSubmatch(msg)
    if len(m) >= 2 { return m[1] }
    return ""
}

// GetDatabases 获取数据库列表
func (s DbInstanceServiceImpl) GetDatabases(c *gin.Context, id uint) {
	instance, err := dao.GetDbInstanceById(id)
	if err != nil {
		result.Failed(c, int(result.ApiCode.FAILED), "实例不存在")
		return
	}

	// 构建数据库信息
	dbInfo := &dbi.DbInfo{
		Type:     dbi.DbType(instance.Type),
		Host:     instance.Host,
		Port:     instance.Port,
		Username: instance.Username,
		Password: instance.Password,
		Network:  instance.Network,
		Params:   instance.Params,
	}

	conn, err := dbm.Conn(context.Background(), dbInfo)
	if err != nil {
		result.Failed(c, int(result.ApiCode.FAILED), fmt.Sprintf("连接失败: %v", err))
		return
	}
	defer conn.Close()

	// 获取数据库列表
	metadata := conn.GetMetadata()
	databases, err := metadata.GetDbNames()
	if err != nil {
		result.Failed(c, int(result.ApiCode.FAILED), fmt.Sprintf("获取数据库列表失败: %v", err))
		return
	}

	result.Success(c, databases)
}

// GetDbConn 获取数据库连接
func (s DbInstanceServiceImpl) GetDbConn(ctx context.Context, instanceId uint, database string) (*dbi.DbConn, error) {
    // Redis 缓存实例信息，减少重复查库
    var instance model.DbInstance
    cacheKey := fmt.Sprintf("db:instance:%d", instanceId)
    if redis.RedisDb != nil {
        if bs, err := redis.RedisDb.Get(context.Background(), cacheKey).Bytes(); err == nil {
            _ = json.Unmarshal(bs, &instance)
        }
    }
    var err error
    if instance.ID == 0 {
        instance, err = dao.GetDbInstanceById(instanceId)
        if err != nil {
            return nil, fmt.Errorf("实例不存在")
        }
        if redis.RedisDb != nil {
            if bs, e := json.Marshal(&instance); e == nil {
                _ = redis.RedisDb.Set(context.Background(), cacheKey, bs, 5*time.Minute).Err()
            }
        }
    }

	// 构建数据库信息
	dbInfo := &dbi.DbInfo{
		InstanceId: uint64(instance.ID),
		Type:       dbi.DbType(instance.Type),
		Host:       instance.Host,
		Port:       instance.Port,
		Username:   instance.Username,
		Password:   instance.Password,
		Database:   database,
		Network:    instance.Network,
		Params:     instance.Params,
	}

	return dbm.GetDbConn(ctx, uint64(instanceId), database, func() (*dbi.DbInfo, error) {
		return dbInfo, nil
	})
}

// CloseDbConn 关闭数据库连接
func (s DbInstanceServiceImpl) CloseDbConn(instanceId uint, database string) {
	if database == "" {
		// 关闭实例下所有数据库连接
		dbm.CloseDb(uint64(instanceId), "")
	} else {
		dbm.CloseDb(uint64(instanceId), database)
	}
}

// DbInstanceService 服务实例
func DbInstanceService() IDbInstanceService {
	return DbInstanceServiceImpl{}
}
