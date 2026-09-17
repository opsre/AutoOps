package service

import (
    "dodevops-api/api/db/elasticsearch/dao"
    "dodevops-api/api/db/elasticsearch/model"
    "dodevops-api/common/result"
    "dodevops-api/common/util"
    "github.com/gin-gonic/gin"
    "github.com/google/uuid"
    "time"
)

type IEsInstanceService interface {
    GetEsInstanceList(c *gin.Context, query *model.EsInstanceQuery)
    GetEsInstanceById(c *gin.Context, id uint)
    SaveEsInstance(c *gin.Context, dto *model.EsInstanceDto)
    DeleteEsInstance(c *gin.Context, id uint)
}

type EsInstanceServiceImpl struct{}

func EsInstanceService() IEsInstanceService { return EsInstanceServiceImpl{} }

func (s EsInstanceServiceImpl) GetEsInstanceList(c *gin.Context, query *model.EsInstanceQuery) {
    if query.PageNum <= 0 { query.PageNum = 1 }
    if query.PageSize <= 0 { query.PageSize = 10 }
    list, total, err := dao.GetEsInstanceList(query)
    if err != nil { result.Failed(c, int(result.ApiCode.FAILED), err.Error()); return }
    for i := range list { if list[i].Password != "" { list[i].Password = "******" } }
    result.SuccessWithPage(c, list, total, query.PageNum, query.PageSize)
}

func (s EsInstanceServiceImpl) GetEsInstanceById(c *gin.Context, id uint) {
    inst, err := dao.GetEsInstanceById(id)
    if err != nil { result.Failed(c, int(result.ApiCode.FAILED), err.Error()); return }
    if inst.Password != "" { inst.Password = "******" }
    result.Success(c, inst)
}

func (s EsInstanceServiceImpl) SaveEsInstance(c *gin.Context, dto *model.EsInstanceDto) {
    now := time.Now()
    if dto.ID == 0 {
        code := dto.Code
        if code == "" { code = uuid.New().String() }
        inst := &model.EsInstance{
            Code:               code,
            Name:               dto.Name,
            Protocol:           dto.Protocol,
            Host:               dto.Host,
            Port:               dto.Port,
            Username:           dto.Username,
            Password:           dto.Password,
            Remark:             dto.Remark,
            SshTunnelMachineId: dto.SshTunnelMachineId,
            Status:             1,
            CreateTime:         util.HTime{Time: now},
            UpdateTime:         util.HTime{Time: now},
            Creator:            c.GetString("username"),
            CreatorId:          uint(c.GetInt("userId")),
        }
        if inst.Protocol == "" { inst.Protocol = "http" }
        if inst.Port == 0 { inst.Port = 9200 }
        if err := dao.CreateEsInstance(inst); err != nil { result.Failed(c, int(result.ApiCode.FAILED), err.Error()); return }
        result.Success(c, inst.ID)
        return
    }
    inst, err := dao.GetEsInstanceById(dto.ID)
    if err != nil { result.Failed(c, int(result.ApiCode.FAILED), "实例不存在"); return }
    inst.Name = dto.Name
    if dto.Protocol != "" { inst.Protocol = dto.Protocol }
    inst.Host = dto.Host
    if dto.Port > 0 { inst.Port = dto.Port }
    inst.Username = dto.Username
    if dto.Password != "" && dto.Password != "******" { inst.Password = dto.Password }
    inst.Remark = dto.Remark
    inst.SshTunnelMachineId = dto.SshTunnelMachineId
    inst.UpdateTime = util.HTime{Time: now}
    inst.Modifier = c.GetString("username")
    inst.ModifierId = uint(c.GetInt("userId"))
    if err := dao.UpdateEsInstance(&inst); err != nil { result.Failed(c, int(result.ApiCode.FAILED), err.Error()); return }
    result.Success(c, nil)
}

func (s EsInstanceServiceImpl) DeleteEsInstance(c *gin.Context, id uint) {
    if err := dao.DeleteEsInstance(id); err != nil { result.Failed(c, int(result.ApiCode.FAILED), err.Error()); return }
    result.Success(c, nil)
}

