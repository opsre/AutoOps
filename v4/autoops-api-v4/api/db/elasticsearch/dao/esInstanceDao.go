package dao

import (
    "dodevops-api/api/db/elasticsearch/model"
    . "dodevops-api/pkg/db"
    "fmt"
)

func CreateEsInstance(m *model.EsInstance) error { return Db.Create(m).Error }

func UpdateEsInstance(m *model.EsInstance) error {
    return Db.Model(&model.EsInstance{}).Where("id = ?", m.ID).Updates(m).Error
}

func DeleteEsInstance(id uint) error { return Db.Delete(&model.EsInstance{}, id).Error }

func GetEsInstanceById(id uint) (m model.EsInstance, err error) {
    err = Db.Where("id = ?", id).First(&m).Error
    return
}

func CheckEsCodeExists(code string, excludeId uint) bool {
    var count int64
    db := Db.Model(&model.EsInstance{}).Where("code = ?", code)
    if excludeId > 0 { db = db.Where("id != ?", excludeId) }
    db.Count(&count)
    return count > 0
}

func GetEsInstanceList(query *model.EsInstanceQuery) (list []model.EsInstance, total int64, err error) {
    db := Db.Model(&model.EsInstance{})
    if query.Name != "" { db = db.Where("name LIKE ?", "%"+query.Name+"%") }
    if query.Host != "" { db = db.Where("host LIKE ?", "%"+query.Host+"%") }
    if query.Status > 0 { db = db.Where("status = ?", query.Status) }
    err = db.Count(&total).Error
    if err != nil { return }
    if query.PageNum > 0 && query.PageSize > 0 {
        offset := (query.PageNum - 1) * query.PageSize
        db = db.Offset(offset).Limit(query.PageSize)
    }
    err = db.Order("create_time DESC").Find(&list).Error
    return
}

func BatchDeleteEsInstances(ids []uint) error {
    if len(ids) == 0 { return fmt.Errorf("ids cannot be empty") }
    return Db.Where("id IN ?", ids).Delete(&model.EsInstance{}).Error
}

