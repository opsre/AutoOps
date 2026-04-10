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
        GetTemplateList(query *model.TemplateQuery) ([]*model.PrometheusAlertDB, int64, error)
        GetTemplateById(id int) (*model.PrometheusAlertDB, error)

        GetAllAlertRouter(query *model.RouterQuery) ([]*model.AlertRouter, int64, error)
	CreateAlertRouter(router *model.AlertRouter) error
	DeleteAlertRouter(id int) error
	UpdateAlertRouter(router *model.AlertRouter) error
	GetAlertRouterById(id int) (*model.AlertRouter, error)
        GetRecordExist(alertname, level, labels, instance, startAt, endAt, summary, description, status string) bool
        GetAlertRecords(query *model.RecordQuery) ([]*model.AlertRecord, int64, error)
	CleanAlertRecords() error
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

func (d *alertDao) GetTemplateList(query *model.TemplateQuery) ([]*model.PrometheusAlertDB, int64, error) {
var list []*model.PrometheusAlertDB
var total int64
db := d.db.Model(&model.PrometheusAlertDB{})
if query.Tplname != "" {
db = db.Where("tplname LIKE ?", "%"+query.Tplname+"%")
}
if query.Tpltype != "" {
db = db.Where("tpltype LIKE ?", "%"+query.Tpltype+"%")
}
if query.Tpluse != "" {
db = db.Where("tpluse LIKE ?", "%"+query.Tpluse+"%")
}
if err := db.Count(&total).Error; err != nil {
return nil, 0, err
}
if query.Page > 0 && query.PageSize > 0 {
db = db.Offset((query.Page - 1) * query.PageSize).Limit(query.PageSize)
}
err := db.Order("id desc").Find(&list).Error
return list, total, err
}

func (d *alertDao) GetTemplateById(id int) (*model.PrometheusAlertDB, error) {
        var tpl model.PrometheusAlertDB
        err := d.db.First(&tpl, id).Error
        return &tpl, err
}

func (d *alertDao) GetAllAlertRouter(query *model.RouterQuery) ([]*model.AlertRouter, int64, error) {
var list []*model.AlertRouter
var total int64
db := d.db.Model(&model.AlertRouter{})
if query.Name != "" {
db = db.Where("name LIKE ?", "%"+query.Name+"%")
}
if query.UrlOrPhone != "" {
db = db.Where("url_or_phone LIKE ?", "%"+query.UrlOrPhone+"%")
}
if query.TplId > 0 {
db = db.Where("tpl_id = ?", query.TplId)
}
if err := db.Count(&total).Error; err != nil {
return nil, 0, err
}
if query.Page > 0 && query.PageSize > 0 {
db = db.Offset((query.Page - 1) * query.PageSize).Limit(query.PageSize)
}
err := db.Preload("Tpl").Order("id desc").Find(&list).Error
return list, total, err
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

func (d *alertDao) CreateAlertRouter(router *model.AlertRouter) error {
return d.db.Create(router).Error
}

func (d *alertDao) DeleteAlertRouter(id int) error {
return d.db.Delete(&model.AlertRouter{}, id).Error
}

func (d *alertDao) UpdateAlertRouter(router *model.AlertRouter) error {
return d.db.Model(router).Select("*").Omit("Created").Updates(router).Error
}

func (d *alertDao) GetAlertRouterById(id int) (*model.AlertRouter, error) {
var router model.AlertRouter
err := d.db.Preload("Tpl").First(&router, id).Error
return &router, err
}

func (d *alertDao) GetAlertRecords(query *model.RecordQuery) ([]*model.AlertRecord, int64, error) {
var list []*model.AlertRecord
var total int64
db := d.db.Model(&model.AlertRecord{})
if query.Alertname != "" {
db = db.Where("alertname LIKE ?", "%"+query.Alertname+"%")
}
if query.AlertLevel != "" {
db = db.Where("alert_level = ?", query.AlertLevel)
}
if query.Instance != "" {
db = db.Where("instance LIKE ?", "%"+query.Instance+"%")
}
if query.AlertStatus != "" {
db = db.Where("alert_status = ?", query.AlertStatus)
}
if err := db.Count(&total).Error; err != nil {
return nil, 0, err
}
if query.Page > 0 && query.PageSize > 0 {
db = db.Offset((query.Page - 1) * query.PageSize).Limit(query.PageSize)
}
err := db.Order("created_time desc").Find(&list).Error
return list, total, err
}

func (d *alertDao) CleanAlertRecords() error {
return d.db.Where("1 = 1").Delete(&model.AlertRecord{}).Error
}
