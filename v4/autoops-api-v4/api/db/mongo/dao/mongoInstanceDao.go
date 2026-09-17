package dao

import (
	"dodevops-api/api/db/mongo/model"
	. "dodevops-api/pkg/db"
)

func CreateMongoInstance(m *model.MongoInstance) error { return Db.Create(m).Error }

func UpdateMongoInstance(m *model.MongoInstance) error {
	return Db.Model(&model.MongoInstance{}).Where("id = ?", m.ID).Updates(m).Error
}

func DeleteMongoInstance(id uint) error { return Db.Delete(&model.MongoInstance{}, id).Error }

func GetMongoInstanceById(id uint) (m model.MongoInstance, err error) {
	err = Db.Where("id = ?", id).First(&m).Error
	return
}

func CheckMongoCodeExists(code string, excludeId uint) bool {
	var count int64
	db := Db.Model(&model.MongoInstance{}).Where("code = ?", code)
	if excludeId > 0 {
		db = db.Where("id != ?", excludeId)
	}
	db.Count(&count)
	return count > 0
}

func GetMongoInstanceList(query *model.MongoInstanceQuery) (list []model.MongoInstance, total int64, err error) {
	db := Db.Model(&model.MongoInstance{})
	if query.Name != "" {
		db = db.Where("name LIKE ?", "%"+query.Name+"%")
	}
	if query.Status > 0 {
		db = db.Where("status = ?", query.Status)
	}
	err = db.Count(&total).Error
	if err != nil {
		return
	}
	if query.PageNum > 0 && query.PageSize > 0 {
		offset := (query.PageNum - 1) * query.PageSize
		db = db.Offset(offset).Limit(query.PageSize)
	}
	err = db.Order("create_time DESC").Find(&list).Error
	return
}

func BatchDeleteMongoInstances(ids []uint) error {
	if len(ids) == 0 {
		return nil
	}
	return Db.Where("id IN ?", ids).Delete(&model.MongoInstance{}).Error
}
