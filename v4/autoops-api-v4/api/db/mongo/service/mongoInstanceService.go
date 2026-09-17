package service

import (
	"context"
	"dodevops-api/api/db/mongo/dao"
	"dodevops-api/api/db/mongo/model"
	"dodevops-api/common/result"
	"dodevops-api/common/util"
	"errors"
	"fmt"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type IMongoInstanceService interface {
	GetMongoInstanceList(c *gin.Context, query *model.MongoInstanceQuery)
	GetMongoInstanceById(c *gin.Context, id uint)
	SaveMongoInstance(c *gin.Context, dto *model.MongoInstanceDto)
	DeleteMongoInstance(c *gin.Context, id uint)
	BatchDeleteMongoInstances(c *gin.Context, ids []uint)
}

type MongoInstanceServiceImpl struct{}

func MongoInstanceService() IMongoInstanceService { return MongoInstanceServiceImpl{} }

func (MongoInstanceServiceImpl) GetMongoInstanceList(c *gin.Context, query *model.MongoInstanceQuery) {
	list, total, err := dao.GetMongoInstanceList(query)
	if err != nil {
		result.Failed(c, int(result.ApiCode.FAILED), err.Error())
		return
	}

	// 转换为 VO 并解析 URI
	voList := make([]model.MongoInstanceVO, len(list))
	for i, inst := range list {
		voList[i] = parseMongoInstanceToVO(&inst)
	}

	result.Success(c, map[string]interface{}{
		"list":     voList,
		"total":    total,
		"pageNum":  query.PageNum,
		"pageSize": query.PageSize,
	})
}

func (MongoInstanceServiceImpl) GetMongoInstanceById(c *gin.Context, id uint) {
	inst, err := dao.GetMongoInstanceById(id)
	if err != nil {
		result.Failed(c, int(result.ApiCode.FAILED), "实例不存在")
		return
	}
	vo := parseMongoInstanceToVO(&inst)
	result.Success(c, vo)
}

func (MongoInstanceServiceImpl) SaveMongoInstance(c *gin.Context, dto *model.MongoInstanceDto) {
	// 验证必填字段
	if dto.Name == "" {
		result.Failed(c, int(result.ApiCode.FAILED), "实例名称不能为空")
		return
	}

	// 构建 MongoDB URI
	var uri string
	if dto.Uri != "" {
		// 如果直接提供了 URI，使用它
		uri = dto.Uri
	} else {
		// 从独立字段构建 URI
		if dto.Host == "" {
			result.Failed(c, int(result.ApiCode.FAILED), "主机地址不能为空")
			return
		}
		if dto.Port == 0 {
			dto.Port = 27017 // 默认端口
		}

		// 构建 URI: mongodb://[username:password@]host:port[?authSource=xxx]
		if dto.Username != "" && dto.Password != "" {
			if dto.AuthDatabase != "" {
				// 使用 authSource 参数指定认证数据库
				uri = fmt.Sprintf("mongodb://%s:%s@%s:%d/?authSource=%s",
					dto.Username, dto.Password, dto.Host, dto.Port, dto.AuthDatabase)
			} else {
				uri = fmt.Sprintf("mongodb://%s:%s@%s:%d",
					dto.Username, dto.Password, dto.Host, dto.Port)
			}
		} else {
			uri = fmt.Sprintf("mongodb://%s:%d", dto.Host, dto.Port)
		}
	}

	now := time.Now()
	inst := model.MongoInstance{
		ID:                 dto.ID,
		Code:               dto.Code,
		Name:               dto.Name,
		Type:               dto.Type,
		Uri:                uri,
		SshTunnelMachineId: dto.SshTunnelMachineId,
		Remark:             dto.Remark,
		Status:             1,
		CreateTime:         util.HTime{Time: now},
		UpdateTime:         util.HTime{Time: now},
	}

	// 如果没有指定类型，默认为 mongodb
	if inst.Type == "" {
		inst.Type = "mongodb"
	}

	var err error
	if inst.ID == 0 {
		// 生成唯一 Code
		inst.Code = generateUniqueCode()
		err = dao.CreateMongoInstance(&inst)
		if err != nil {
			result.Failed(c, int(result.ApiCode.FAILED), "创建失败: "+err.Error())
			return
		}
		result.Success(c, inst)
	} else {
		// 更新不修改 Code
		inst.Code = ""
		err = dao.UpdateMongoInstance(&inst)
		if err != nil {
			result.Failed(c, int(result.ApiCode.FAILED), "更新失败: "+err.Error())
			return
		}
		result.Success(c, "更新成功")
	}
}

func (MongoInstanceServiceImpl) DeleteMongoInstance(c *gin.Context, id uint) {
	err := dao.DeleteMongoInstance(id)
	if err != nil {
		result.Failed(c, int(result.ApiCode.FAILED), err.Error())
		return
	}
	result.Success(c, "删除成功")
}

func (MongoInstanceServiceImpl) BatchDeleteMongoInstances(c *gin.Context, ids []uint) {
	err := dao.BatchDeleteMongoInstances(ids)
	if err != nil {
		result.Failed(c, int(result.ApiCode.FAILED), err.Error())
		return
	}
	result.Success(c, fmt.Sprintf("批量删除成功，共 %d 条", len(ids)))
}

// maskMongoUri 脱敏 MongoDB URI 中的密码
func maskMongoUri(uri string) string {
	// 简单的密码脱敏处理
	// mongodb://user:password@host:port/db -> mongodb://user:******@host:port/db
	// 使用正则表达式会更准确，这里简化处理
	if len(uri) > 10 && uri[:10] == "mongodb://" {
		return "mongodb://******"
	}
	if len(uri) > 14 && uri[:14] == "mongodb+srv://" {
		return "mongodb+srv://******"
	}
	return "******"
}

// generateUniqueCode 生成唯一编码
func generateUniqueCode() string {
	return fmt.Sprintf("mongo_%d", time.Now().UnixNano())
}

// 获取 MongoDB 客户端连接
func getMongoClient(uri string) (*mongo.Client, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	clientOptions := options.Client().ApplyURI(uri)
	clientOptions.SetServerSelectionTimeout(10 * time.Second)
	clientOptions.SetConnectTimeout(10 * time.Second)

	// 注意：不要强制设置 TLS 配置
	// 某些 MongoDB 服务器不支持 TLS，强制设置会导致连接失败
	// 如果需要 TLS，应该在 URI 中指定：mongodb://...?tls=true

	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		return nil, err
	}

	// 测试连接
	err = client.Ping(ctx, nil)
	if err != nil {
		client.Disconnect(ctx)
		return nil, err
	}

	return client, nil
}

// 解析 DocId，支持 ObjectID 和普通字符串
func parseDocId(docId any) (any, error) {
	switch v := docId.(type) {
	case string:
		// 尝试转换为 ObjectID
		if objId, err := primitive.ObjectIDFromHex(v); err == nil {
			return objId, nil
		}
		return v, nil
	default:
		return docId, nil
	}
}

// 验证 Limit 参数
func validateLimit(limit int64) error {
	if limit > 0 && limit > 100 {
		return errors.New("limit 不能超过 100")
	}
	return nil
}

// parseMongoInstanceToVO 将 MongoInstance 转换为 VO，解析 URI
func parseMongoInstanceToVO(inst *model.MongoInstance) model.MongoInstanceVO {
	vo := model.MongoInstanceVO{
		ID:                 inst.ID,
		Code:               inst.Code,
		Name:               inst.Name,
		Type:               inst.Type,
		SshTunnelMachineId: inst.SshTunnelMachineId,
		Remark:             inst.Remark,
		Status:             inst.Status,
		CreateTime:         inst.CreateTime,
		UpdateTime:         inst.UpdateTime,
		Creator:            inst.Creator,
		CreatorId:          inst.CreatorId,
		Modifier:           inst.Modifier,
		ModifierId:         inst.ModifierId,
		Password:           "******", // 默认脱敏
	}

	// 解析 URI
	if inst.Uri != "" {
		parsedURL, err := url.Parse(inst.Uri)
		if err == nil {
			// 解析 host 和 port
			host := parsedURL.Hostname()
			port := parsedURL.Port()
			if port == "" {
				port = "27017" // 默认端口
			}
			vo.Host = host
			if p, err := strconv.Atoi(port); err == nil {
				vo.Port = p
			}

			// 解析用户名
			if parsedURL.User != nil {
				vo.Username = parsedURL.User.Username()
			}

			// 解析 authDatabase (从 query 参数 authSource 获取)
			queryParams := parsedURL.Query()
			if authSource := queryParams.Get("authSource"); authSource != "" {
				vo.AuthDatabase = authSource
			} else if parsedURL.Path != "" && parsedURL.Path != "/" {
				// 兼容旧格式：从 Path 获取
				vo.AuthDatabase = strings.TrimPrefix(parsedURL.Path, "/")
			}
		}
	}

	return vo
}
