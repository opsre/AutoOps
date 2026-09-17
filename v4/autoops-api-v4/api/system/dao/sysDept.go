package dao

import (
	"dodevops-api/api/system/model"
	"dodevops-api/common/util"
	. "dodevops-api/pkg/db"
	"time"
)

// 查询部门列表
func GetSysDeptList(DeptName string, DeptStatus string) (sysDept []model.SysDept) {
	curDb := Db.Table("sys_dept")
	if DeptName != "" {
		curDb = curDb.Where("dept_name = ?", DeptName)
	}
	if DeptStatus != "" {
		curDb = curDb.Where("dept_status = ?", DeptStatus)
	}
	curDb.Find(&sysDept)
	return sysDept
}

// 根据部门名称查询
func GetSysDeptByName(deptName string) (sysDept model.SysDept) {
	Db.Where("dept_name = ?", deptName).First(&sysDept)
	return sysDept
}

// 新增部门 之前先查询部门是否存在。
func CreateSysDept(sysDept model.SysDept) bool {
	sysDeptByName := GetSysDeptByName(sysDept.DeptName)
	if sysDeptByName.ID > 0 {
		return false
	}
	if sysDept.DeptType == 1 {
		sysDept := model.SysDept{
			DeptStatus: sysDept.DeptStatus,
			ParentId:   0,
			DeptType:   sysDept.DeptType,
			DeptName:   sysDept.DeptName,
			CreateTime: util.HTime{Time: time.Now()},
		}
		Db.Create(&sysDept)
		return true
	} else {
		sysDept := model.SysDept{
			DeptStatus: sysDept.DeptStatus,
			ParentId:   sysDept.ParentId,
			DeptType:   sysDept.DeptType,
			DeptName:   sysDept.DeptName,
			CreateTime: util.HTime{Time: time.Now()},
		}
		Db.Create(&sysDept)
		return true
	}
}

// 根据id查询部门
func GetSysDeptById(Id int) (sysDept model.SysDept) {
	Db.First(&sysDept, Id)
	return sysDept
}

// 修改部门
func UpdateSysDept(dept model.SysDept) (sysDept model.SysDept) {
	Db.First(&sysDept, dept.ID)
	sysDept.ParentId = dept.ParentId
	sysDept.DeptType = dept.DeptType
	sysDept.DeptName = dept.DeptName
	sysDept.DeptStatus = dept.DeptStatus
	Db.Save(&sysDept)
	return sysDept
}

// 查询部门是否有人
func GetSysAdminDept(id int) (sysAdmin model.SysAdmin) {
	Db.Where("dept_id = ?", id).First(&sysAdmin)
	return sysAdmin
}

// 删除部门前先确认部门及其子部门是否有人
func DeleteSysDeptById(dto model.SysDeptIdDto) bool {
	// 获取当前部门及所有子部门的ID
	deptIds := GetAllChildDeptIds(dto.Id)

	// 检查这些部门中是否有用户
	var count int64
	Db.Model(&model.SysAdmin{}).Where("dept_id IN ?", deptIds).Count(&count)
	if count > 0 {
		return false
	}

	// 递归删除所有子部门
	Db.Where("parent_id = ?", dto.Id).Delete(&model.SysDept{})
	// 删除当前部门
	Db.Delete(&model.SysDept{}, dto.Id)
	return true
}

// 部门下拉列表
func QuerySysDeptVoList() (sysDeptVo []model.SysDeptVo) {
	Db.Table("sys_dept").Select("id, dept_name AS label, parent_id").Scan(&sysDeptVo)
	return sysDeptVo
}

// 递归获取所有子部门ID
func GetAllChildDeptIds(deptId int) []int {
	deptIds := []int{deptId}

	var childDepts []model.SysDept
	Db.Where("parent_id = ?", deptId).Find(&childDepts)

	for _, child := range childDepts {
		childIds := GetAllChildDeptIds(int(child.ID))
		deptIds = append(deptIds, childIds...)
	}

	return deptIds
}

// 获取某部门下的所有用户（包含子部门）
func GetUsersByDeptId(deptId int) []model.SysAdminVo {
	// 递归获取当前部门及所有子部门的ID
	deptIds := GetAllChildDeptIds(deptId)

	var users []model.SysAdminVo
	Db.Table("sys_admin").
		Select(`
			sys_admin.*,
			sys_admin_role.role_id AS role_id,
			sys_dept.dept_name AS dept_name,
			sys_post.post_name AS post_name,
			sys_role.role_name AS role_name
		`).
		Joins("LEFT JOIN sys_admin_role ON sys_admin.id = sys_admin_role.admin_id").
		Joins("LEFT JOIN sys_dept ON sys_admin.dept_id = sys_dept.id").
		Joins("LEFT JOIN sys_post ON sys_admin.post_id = sys_post.id").
		Joins("LEFT JOIN sys_role ON sys_admin_role.role_id = sys_role.id").
		Where("sys_admin.dept_id IN ?", deptIds).
		Scan(&users)
	return users
}
