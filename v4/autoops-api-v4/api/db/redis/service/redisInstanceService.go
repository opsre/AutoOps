package service

import (
    "dodevops-api/api/db/redis/dao"
    "dodevops-api/api/db/redis/model"
    "dodevops-api/common/result"
    "dodevops-api/common/util"
    "github.com/gin-gonic/gin"
    "github.com/google/uuid"
    "time"
)

type IRedisInstanceService interface {
    GetRedisInstanceList(c *gin.Context, query *model.RedisInstanceQuery)
    GetRedisInstanceById(c *gin.Context, id uint)
    SaveRedisInstance(c *gin.Context, dto *model.RedisInstanceDto)
    DeleteRedisInstance(c *gin.Context, id uint)
}

type RedisInstanceServiceImpl struct{}

func RedisInstanceService() IRedisInstanceService { return RedisInstanceServiceImpl{} }

func (s RedisInstanceServiceImpl) GetRedisInstanceList(c *gin.Context, query *model.RedisInstanceQuery) {
    if query.PageNum <= 0 { query.PageNum = 1 }
    if query.PageSize <= 0 { query.PageSize = 10 }
    list, total, err := dao.GetRedisInstanceList(query)
    if err != nil { result.Failed(c, int(result.ApiCode.FAILED), err.Error()); return }
    for i := range list { list[i].Password = "******"; if list[i].RedisNodePassword != "" { list[i].RedisNodePassword = "******" } }
    result.SuccessWithPage(c, list, total, query.PageNum, query.PageSize)
}

func (s RedisInstanceServiceImpl) GetRedisInstanceById(c *gin.Context, id uint) {
    inst, err := dao.GetRedisInstanceById(id)
    if err != nil { result.Failed(c, int(result.ApiCode.FAILED), err.Error()); return }
    inst.Password = "******"; if inst.RedisNodePassword != "" { inst.RedisNodePassword = "******" }
    result.Success(c, inst)
}

func (s RedisInstanceServiceImpl) SaveRedisInstance(c *gin.Context, dto *model.RedisInstanceDto) {
    now := time.Now()
    if dto.ID == 0 {
        // create
        code := dto.Code
        if code == "" { code = uuid.New().String() }
        inst := &model.RedisInstance{
            Code:               code,
            Name:               dto.Name,
            Mode:               dto.Mode,
            Host:               dto.Host,
            Port:               dto.Port,
            DB:                 dto.DB,
            Username:           dto.Username,
            Password:           dto.Password, // TODO encrypt
            RedisNodePassword:  dto.RedisNodePassword,
            Remark:             dto.Remark,
            SshTunnelMachineId: dto.SshTunnelMachineId,
            Status:             1,
            CreateTime:         util.HTime{Time: now},
            UpdateTime:         util.HTime{Time: now},
            Creator:            c.GetString("username"),
            CreatorId:          uint(c.GetInt("userId")),
        }
        if err := dao.CreateRedisInstance(inst); err != nil { result.Failed(c, int(result.ApiCode.FAILED), err.Error()); return }
        result.Success(c, inst.ID)
        return
    }
    // update
    inst, err := dao.GetRedisInstanceById(dto.ID)
    if err != nil { result.Failed(c, int(result.ApiCode.FAILED), "实例不存在"); return }
    inst.Name = dto.Name
    inst.Mode = dto.Mode
    inst.Host = dto.Host
    inst.Port = dto.Port
    inst.DB = dto.DB
    inst.Username = dto.Username
    if dto.Password != "" && dto.Password != "******" { inst.Password = dto.Password }
    if dto.RedisNodePassword != "" && dto.RedisNodePassword != "******" { inst.RedisNodePassword = dto.RedisNodePassword }
    inst.Remark = dto.Remark
    inst.SshTunnelMachineId = dto.SshTunnelMachineId
    inst.UpdateTime = util.HTime{Time: now}
    inst.Modifier = c.GetString("username")
    inst.ModifierId = uint(c.GetInt("userId"))
    if err := dao.UpdateRedisInstance(&inst); err != nil { result.Failed(c, int(result.ApiCode.FAILED), err.Error()); return }
    result.Success(c, nil)
}

func (s RedisInstanceServiceImpl) DeleteRedisInstance(c *gin.Context, id uint) {
    if err := dao.DeleteRedisInstance(id); err != nil { result.Failed(c, int(result.ApiCode.FAILED), err.Error()); return }
    result.Success(c, nil)
}

