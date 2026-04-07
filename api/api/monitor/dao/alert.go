package dao

import (
        "dodevops-api/api/monitor/model"
        "dodevops-api/common"

        "gorm.io/gorm"
)

type AlertDao interface {
        // PrometheusAlertDB 模板相关
        CreateTemplate(tpl *model.PrometheusAlertDB) error
        DeleteTemplate(id int) error
        UpdateTemplate(tpl *model.PrometheusAlertDB) error
        GetTemplateList() ([]*model.PrometheusAlertDB, error)
        GetTemplateById(id int) (*model.PrometheusAlertDB, error)

        GetAllAlertRouter() ([]*model.AlertRouter, error)
        GetRecordExist(alertname, level, labels, instance, startAt, endAt, summary, description, status string) bool
        AddAlertRecord(record *model.AlertRecord) error
        GetAlertConfig(key string) string
        GetAllAlertConfig() map[string]string
}

type alertDao struct {
        db *gorm.DB
}

func NewAlertDao() AlertDao {
        return &alertDao{
                db: common.GetDB(),
        }
}

func (d *alertDao) CreateTemplate(tpl *model.PrometheusAlertDB) error {
        return d.db.Create(tpl).Error
}

func (d *alertDao) DeleteTemplate(id int) error {
        return d.db.Delete(&model.PrometheusAlertDB{}, id).Error
}

func (d *alertDao) UpdateTemplate(tpl *model.PrometheusAlertDB) error {
        return d.db.Save(tpl).Error
}

func (d *alertDao) GetTemplateList() ([]*model.PrometheusAlertDB, error) {
        var list []*model.PrometheusAlertDB
        err := d.db.Find(&list).Error
        return list, err
}

func (d *alertDao) GetTemplateById(id int) (*model.PrometheusAlertDB, error) {
        var tpl model.PrometheusAlertDB
        err := d.db.First(&tpl, id).Error
        return &tpl, err
}

func (d *alertDao) GetAllAlertRouter() ([]*model.AlertRouter, error) {
        var list []*model.AlertRouter
        err := d.db.Preload("Tpl").Find(&list).Error
        return list, err
}

func (d *alertDao) GetRecordExist(alertname, level, labels, instance, startAt, endAt, summary, description, status string) bool {
        var count int64
        d.db.Model(&model.AlertRecord{}).Where("alertname = ? AND alert_level = ? AND labels = ? AND instance = ? AND starts_at = ? AND ends_at = ? AND summary = ? AND description = ? AND alert_status = ?",
                alertname, level, labels, instance, startAt, endAt, summary, description, status).Count(&count)
        return count > 0
}

func (d *alertDao) AddAlertRecord(record *model.AlertRecord) error {
        return d.db.Create(record).Error
}

func (d *alertDao) GetAlertConfig(key string) string {
        var cfg model.AlertConfig
        err := d.db.Where("conf_key = ?", key).First(&cfg).Error
        if err != nil {
                return ""
        }
        return cfg.ConfValue
}

func (d *alertDao) GetAllAlertConfig() map[string]string {
        var list []*model.AlertConfig
        configs := make(map[string]string)
        err := d.db.Find(&list).Error
        if err != nil {
                return configs
        }
        for _, cfg := range list {
                configs[cfg.ConfKey] = cfg.ConfValue
        }
        return configs
}
