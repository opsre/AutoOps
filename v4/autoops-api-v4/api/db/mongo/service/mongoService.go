package service

import (
	"context"
	"dodevops-api/api/db/mongo/dao"
	"dodevops-api/api/db/mongo/model"
	"dodevops-api/common/result"
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type IMongoService interface {
	TestConn(c *gin.Context, dto *model.MongoInstanceDto)
	GetDatabases(c *gin.Context, id uint)
	GetCollections(c *gin.Context, id uint, database string)
	RunCommand(c *gin.Context, id uint, form *model.MongoRunCommandForm)
	Find(c *gin.Context, id uint, form *model.MongoFindForm)
	UpdateById(c *gin.Context, id uint, form *model.MongoUpdateByIdForm)
	DeleteById(c *gin.Context, id uint, form *model.MongoDeleteByIdForm)
	InsertOne(c *gin.Context, id uint, form *model.MongoInsertForm)
}

type MongoServiceImpl struct{}

func MongoService() IMongoService { return MongoServiceImpl{} }

func (MongoServiceImpl) TestConn(c *gin.Context, dto *model.MongoInstanceDto) {
	// 构建 MongoDB URI
	var uri string
	if dto.Uri != "" {
		uri = dto.Uri
	} else {
		// 从独立字段构建 URI
		if dto.Host == "" {
			result.Failed(c, int(result.ApiCode.FAILED), "主机地址不能为空")
			return
		}
		if dto.Port == 0 {
			dto.Port = 27017
		}

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

	client, err := getMongoClient(uri)
	if err != nil {
		result.Failed(c, int(result.ApiCode.FAILED), "连接失败: "+err.Error())
		return
	}
	defer client.Disconnect(context.Background())

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// 获取服务器信息
	var serverStatus bson.M
	err = client.Database("admin").RunCommand(ctx, bson.D{{Key: "serverStatus", Value: 1}}).Decode(&serverStatus)
	if err != nil {
		result.Failed(c, int(result.ApiCode.FAILED), "连接测试失败: "+err.Error())
		return
	}

	result.Success(c, map[string]interface{}{
		"status":  "connected",
		"version": serverStatus["version"],
		"host":    serverStatus["host"],
	})
}

func (MongoServiceImpl) GetDatabases(c *gin.Context, id uint) {
	inst, err := dao.GetMongoInstanceById(id)
	if err != nil {
		result.Failed(c, int(result.ApiCode.FAILED), "实例不存在")
		return
	}

	client, err := getMongoClient(inst.Uri)
	if err != nil {
		result.Failed(c, int(result.ApiCode.FAILED), "连接失败: "+err.Error())
		return
	}
	defer client.Disconnect(context.Background())

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	databases, err := client.ListDatabases(ctx, bson.D{})
	if err != nil {
		result.Failed(c, int(result.ApiCode.FAILED), "获取数据库列表失败: "+err.Error())
		return
	}

	result.Success(c, databases)
}

func (MongoServiceImpl) GetCollections(c *gin.Context, id uint, database string) {
	if database == "" {
		result.Failed(c, int(result.ApiCode.FAILED), "database 参数不能为空")
		return
	}

	inst, err := dao.GetMongoInstanceById(id)
	if err != nil {
		result.Failed(c, int(result.ApiCode.FAILED), "实例不存在")
		return
	}

	client, err := getMongoClient(inst.Uri)
	if err != nil {
		result.Failed(c, int(result.ApiCode.FAILED), "连接失败: "+err.Error())
		return
	}
	defer client.Disconnect(context.Background())

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	collections, err := client.Database(database).ListCollectionNames(ctx, bson.D{})
	if err != nil {
		result.Failed(c, int(result.ApiCode.FAILED), "获取集合列表失败: "+err.Error())
		return
	}

	result.Success(c, collections)
}

func (MongoServiceImpl) RunCommand(c *gin.Context, id uint, form *model.MongoRunCommandForm) {
	inst, err := dao.GetMongoInstanceById(id)
	if err != nil {
		result.Failed(c, int(result.ApiCode.FAILED), "实例不存在")
		return
	}

	client, err := getMongoClient(inst.Uri)
	if err != nil {
		result.Failed(c, int(result.ApiCode.FAILED), "连接失败: "+err.Error())
		return
	}
	defer client.Disconnect(context.Background())

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// 构建命令
	commands := bson.D{}
	for _, cmd := range form.Command {
		for k, v := range cmd {
			commands = append(commands, bson.E{Key: k, Value: v})
		}
	}

	var commandResult bson.M
	err = client.Database(form.Database).RunCommand(ctx, commands).Decode(&commandResult)
	if err != nil {
		result.Failed(c, int(result.ApiCode.FAILED), "命令执行失败: "+err.Error())
		return
	}

	result.Success(c, commandResult)
}

func (MongoServiceImpl) Find(c *gin.Context, id uint, form *model.MongoFindForm) {
	if err := validateLimit(form.Limit); err != nil {
		result.Failed(c, int(result.ApiCode.FAILED), err.Error())
		return
	}

	inst, err := dao.GetMongoInstanceById(id)
	if err != nil {
		result.Failed(c, int(result.ApiCode.FAILED), "实例不存在")
		return
	}

	client, err := getMongoClient(inst.Uri)
	if err != nil {
		result.Failed(c, int(result.ApiCode.FAILED), "连接失败: "+err.Error())
		return
	}
	defer client.Disconnect(context.Background())

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// 处理 _id 查询字段
	filter := form.Filter
	if filter == nil {
		filter = bson.M{}
	}
	if idStr, ok := filter["_id"].(string); ok && idStr != "" {
		if objId, err := parseDocId(idStr); err == nil {
			filter["_id"] = objId
		}
	}

	// 设置查询选项
	findOptions := options.Find()
	if form.Sort != nil {
		findOptions.SetSort(form.Sort)
	}
	if form.Skip > 0 {
		findOptions.SetSkip(form.Skip)
	}
	if form.Limit > 0 {
		findOptions.SetLimit(form.Limit)
	} else {
		findOptions.SetLimit(20) // 默认限制 20 条
	}

	cursor, err := client.Database(form.Database).Collection(form.Collection).Find(ctx, filter, findOptions)
	if err != nil {
		result.Failed(c, int(result.ApiCode.FAILED), "查询失败: "+err.Error())
		return
	}
	defer cursor.Close(ctx)

	var results []bson.M
	if err = cursor.All(ctx, &results); err != nil {
		result.Failed(c, int(result.ApiCode.FAILED), "读取结果失败: "+err.Error())
		return
	}

	result.Success(c, results)
}

func (MongoServiceImpl) UpdateById(c *gin.Context, id uint, form *model.MongoUpdateByIdForm) {
	inst, err := dao.GetMongoInstanceById(id)
	if err != nil {
		result.Failed(c, int(result.ApiCode.FAILED), "实例不存在")
		return
	}

	client, err := getMongoClient(inst.Uri)
	if err != nil {
		result.Failed(c, int(result.ApiCode.FAILED), "连接失败: "+err.Error())
		return
	}
	defer client.Disconnect(context.Background())

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// 获取文档 ID，优先级：filter._id > _id > docId
	var docId any
	if form.Filter != nil && form.Filter["_id"] != nil {
		docId = form.Filter["_id"]
	} else if form.ID != nil {
		docId = form.ID
	} else if form.DocId != nil {
		docId = form.DocId
	} else {
		result.Failed(c, int(result.ApiCode.FAILED), "缺少文档 ID（filter._id、_id 或 docId）")
		return
	}

	// 解析 docId
	parsedDocId, err := parseDocId(docId)
	if err != nil {
		result.Failed(c, int(result.ApiCode.FAILED), "DocId 格式错误: "+err.Error())
		return
	}

	// 处理 update 字段：如果不包含 $ 开头的操作符，自动包装为 $set
	updateDoc := form.Update
	hasOperator := false
	for key := range form.Update {
		if len(key) > 0 && key[0] == '$' {
			hasOperator = true
			break
		}
	}
	if !hasOperator {
		// 移除 _id 字段（不能更新 _id）
		updateData := make(map[string]any)
		for key, value := range form.Update {
			if key != "_id" {
				updateData[key] = value
			}
		}
		updateDoc = bson.M{"$set": updateData}
	}

	updateResult, err := client.Database(form.Database).Collection(form.Collection).UpdateByID(ctx, parsedDocId, updateDoc)
	if err != nil {
		result.Failed(c, int(result.ApiCode.FAILED), "更新失败: "+err.Error())
		return
	}

	result.Success(c, updateResult)
}

func (MongoServiceImpl) DeleteById(c *gin.Context, id uint, form *model.MongoDeleteByIdForm) {
	inst, err := dao.GetMongoInstanceById(id)
	if err != nil {
		result.Failed(c, int(result.ApiCode.FAILED), "实例不存在")
		return
	}

	client, err := getMongoClient(inst.Uri)
	if err != nil {
		result.Failed(c, int(result.ApiCode.FAILED), "连接失败: "+err.Error())
		return
	}
	defer client.Disconnect(context.Background())

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// 获取文档 ID，优先级：filter._id > _id > docId
	var docId any
	if form.Filter != nil && form.Filter["_id"] != nil {
		docId = form.Filter["_id"]
	} else if form.ID != nil {
		docId = form.ID
	} else if form.DocId != nil {
		docId = form.DocId
	} else {
		result.Failed(c, int(result.ApiCode.FAILED), "缺少文档 ID（filter._id、_id 或 docId）")
		return
	}

	// 解析 docId
	parsedDocId, err := parseDocId(docId)
	if err != nil {
		result.Failed(c, int(result.ApiCode.FAILED), "DocId 格式错误: "+err.Error())
		return
	}

	deleteResult, err := client.Database(form.Database).Collection(form.Collection).DeleteOne(ctx, bson.D{{Key: "_id", Value: parsedDocId}})
	if err != nil {
		result.Failed(c, int(result.ApiCode.FAILED), "删除失败: "+err.Error())
		return
	}

	result.Success(c, deleteResult)
}

func (MongoServiceImpl) InsertOne(c *gin.Context, id uint, form *model.MongoInsertForm) {
	inst, err := dao.GetMongoInstanceById(id)
	if err != nil {
		result.Failed(c, int(result.ApiCode.FAILED), "实例不存在")
		return
	}

	client, err := getMongoClient(inst.Uri)
	if err != nil {
		result.Failed(c, int(result.ApiCode.FAILED), "连接失败: "+err.Error())
		return
	}
	defer client.Disconnect(context.Background())

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// 获取文档数据，优先级：document > doc
	var doc map[string]any
	if form.Document != nil && len(form.Document) > 0 {
		doc = form.Document
	} else if form.Doc != nil && len(form.Doc) > 0 {
		doc = form.Doc
	} else {
		result.Failed(c, int(result.ApiCode.FAILED), "缺少文档数据（document 或 doc）")
		return
	}

	insertResult, err := client.Database(form.Database).Collection(form.Collection).InsertOne(ctx, doc)
	if err != nil {
		result.Failed(c, int(result.ApiCode.FAILED), "插入失败: "+err.Error())
		return
	}

	result.Success(c, insertResult)
}
