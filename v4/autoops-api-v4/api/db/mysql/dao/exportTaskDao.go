package dao

import (
	"dodevops-api/api/db/mysql/model"
	. "dodevops-api/pkg/db"
)

// CreateExportTask 创建导出任务
func CreateExportTask(task *model.ExportTask) error {
	return Db.Create(task).Error
}

// GetExportTaskByTaskId 根据任务ID获取导出任务
func GetExportTaskByTaskId(taskId string) (*model.ExportTask, error) {
	var task model.ExportTask
	err := Db.Where("task_id = ?", taskId).First(&task).Error
	return &task, err
}

// UpdateExportTask 更新导出任务
func UpdateExportTask(task *model.ExportTask) error {
	return Db.Save(task).Error
}

// DeleteExportTask 删除导出任务
func DeleteExportTask(id uint) error {
	return Db.Delete(&model.ExportTask{}, id).Error
}

// DeleteExportTaskByTaskId 根据任务ID删除导出任务
func DeleteExportTaskByTaskId(taskId string) error {
	return Db.Where("task_id = ?", taskId).Delete(&model.ExportTask{}).Error
}
