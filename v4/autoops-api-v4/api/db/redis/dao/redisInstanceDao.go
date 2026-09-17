package dao

import (
    "dodevops-api/api/db/redis/model"
    . "dodevops-api/pkg/db"
    "fmt"
)

func CreateRedisInstance(m *model.RedisInstance) error { return Db.Create(m).Error }

func UpdateRedisInstance(m *model.RedisInstance) error {
    return Db.Model(&model.RedisInstance{}).Where("id = ?", m.ID).Updates(m).Error
}

func DeleteRedisInstance(id uint) error { return Db.Delete(&model.RedisInstance{}, id).Error }

func GetRedisInstanceById(id uint) (m model.RedisInstance, err error) {
    err = Db.Where("id = ?", id).First(&m).Error
    return
}

func CheckRedisCodeExists(code string, excludeId uint) bool {
    var count int64
    db := Db.Model(&model.RedisInstance{}).Where("code = ?", code)
    if excludeId > 0 { db = db.Where("id != ?", excludeId) }
    db.Count(&count)
    return count > 0
}

func GetRedisInstanceList(query *model.RedisInstanceQuery) (list []model.RedisInstance, total int64, err error) {
    db := Db.Model(&model.RedisInstance{})
    if query.Name != "" { db = db.Where("name LIKE ?", "%"+query.Name+"%") }
    if query.Host != "" { db = db.Where("host LIKE ?", "%"+query.Host+"%") }
    if query.Mode != "" { db = db.Where("mode = ?", query.Mode) }
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

func BatchDeleteRedisInstances(ids []uint) error {
    if len(ids) == 0 { return fmt.Errorf("ids cannot be empty") }
    return Db.Where("id IN ?", ids).Delete(&model.RedisInstance{}).Error
}

