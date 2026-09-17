package service

import (
	"errors"

	cmdbdao "dodevops-api/api/cmdb/dao"
	"dodevops-api/api/monitor/dao"
	"dodevops-api/api/monitor/model"
	"gorm.io/gorm"
)

type APIBusinessFunctionService interface {
	Create(dto *model.CreateAPIBusinessFunctionDto) (*model.APIBusinessFunction, error)
	Update(dto *model.UpdateAPIBusinessFunctionDto) (*model.APIBusinessFunction, error)
	Delete(id uint) error
	GetByID(id uint) (*model.APIBusinessFunction, error)
	GetByGroupIDAndName(groupID uint, name string) (*model.APIBusinessFunction, error)
	GetAll() ([]*model.APIBusinessFunction, error)
	GetByGroupID(groupID uint) ([]*model.APIBusinessFunction, error)
	GetList(groupID *uint, keyword string, page, pageSize int) ([]*model.APIBusinessFunction, int64, error)
	UpdateStatus(id uint, status int) error
	CreateIfNotExists(groupID uint, name string) (*model.APIBusinessFunction, error)
	CheckNameExists(groupID uint, name string, excludeID uint) (bool, error)
}

type APIBusinessFunctionServiceImpl struct {
	functionDao dao.APIBusinessFunctionDao
}

func NewAPIBusinessFunctionService() APIBusinessFunctionService {
	return &APIBusinessFunctionServiceImpl{
		functionDao: dao.NewAPIBusinessFunctionDao(),
	}
}

// Create 创建业务功能
func (s *APIBusinessFunctionServiceImpl) Create(dto *model.CreateAPIBusinessFunctionDto) (*model.APIBusinessFunction, error) {
	// 检查同一业务组内是否已存在同名功能
	exists, err := s.CheckNameExists(dto.GroupID, dto.Name, 0)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, errors.New("业务功能名称已存在: " + dto.Name)
	}

	// 检查业务组是否存在
	cmdbGroupDao := cmdbdao.NewCmdbGroupDao()
	cmdbGroup, err := cmdbGroupDao.GetCmdbGroupById(dto.GroupID)
	if err != nil || cmdbGroup.ID == 0 {
		return nil, errors.New("业务组不存在")
	}

	function := &model.APIBusinessFunction{
		Name:        dto.Name,
		Description: dto.Description,
		GroupID:     dto.GroupID,
		GroupName:   cmdbGroup.Name,
		Sort:        dto.Sort,
		Status:      1, // 默认启用
	}

	if err := s.functionDao.Create(function); err != nil {
		return nil, err
	}

	return function, nil
}

// Update 更新业务功能
func (s *APIBusinessFunctionServiceImpl) Update(dto *model.UpdateAPIBusinessFunctionDto) (*model.APIBusinessFunction, error) {
	// 检查同一业务组内是否已存在同名功能（排除当前记录）
	exists, err := s.CheckNameExists(dto.GroupID, dto.Name, dto.ID)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, errors.New("业务功能名称已存在: " + dto.Name)
	}

	function := &model.APIBusinessFunction{
		ID:          dto.ID,
		Name:        dto.Name,
		Description: dto.Description,
		GroupID:     dto.GroupID,
		Sort:        dto.Sort,
	}

	if err := s.functionDao.Update(function); err != nil {
		return nil, err
	}

	return function, nil
}

// Delete 删除业务功能
func (s *APIBusinessFunctionServiceImpl) Delete(id uint) error {
	return s.functionDao.Delete(id)
}

// GetByID 根据ID获取业务功能
func (s *APIBusinessFunctionServiceImpl) GetByID(id uint) (*model.APIBusinessFunction, error) {
	return s.functionDao.GetByID(id)
}

// GetByGroupIDAndName 根据业务组ID和名称获取业务功能
func (s *APIBusinessFunctionServiceImpl) GetByGroupIDAndName(groupID uint, name string) (*model.APIBusinessFunction, error) {
	return s.functionDao.GetByGroupIDAndName(groupID, name)
}

// GetAll 获取所有业务功能
func (s *APIBusinessFunctionServiceImpl) GetAll() ([]*model.APIBusinessFunction, error) {
	return s.functionDao.GetAll()
}

// GetByGroupID 根据业务组ID获取业务功能列表
func (s *APIBusinessFunctionServiceImpl) GetByGroupID(groupID uint) ([]*model.APIBusinessFunction, error) {
	return s.functionDao.GetByGroupID(groupID)
}

// GetList 分页查询业务功能列表
func (s *APIBusinessFunctionServiceImpl) GetList(groupID *uint, keyword string, page, pageSize int) ([]*model.APIBusinessFunction, int64, error) {
	return s.functionDao.GetList(groupID, keyword, page, pageSize)
}

// UpdateStatus 更新状态
func (s *APIBusinessFunctionServiceImpl) UpdateStatus(id uint, status int) error {
	return s.functionDao.UpdateStatus(id, status)
}

// CreateIfNotExists 检查同名功能，不存在则创建，存在则返回错误
func (s *APIBusinessFunctionServiceImpl) CreateIfNotExists(groupID uint, name string) (*model.APIBusinessFunction, error) {
	// 检查功能是否已存在
	existFunc, err := s.functionDao.GetByGroupIDAndName(groupID, name)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	// 如果已存在，返回错误
	if existFunc != nil {
		return nil, errors.New("业务功能名称已存在: " + name)
	}

	// 检查业务组是否存在
	cmdbGroupDao := cmdbdao.NewCmdbGroupDao()
	cmdbGroup, err := cmdbGroupDao.GetCmdbGroupById(groupID)
	if err != nil || cmdbGroup.ID == 0 {
		return nil, errors.New("业务组不存在")
	}

	// 创建新功能
	function := &model.APIBusinessFunction{
		Name:      name,
		GroupID:   groupID,
		GroupName: cmdbGroup.Name,
		Status:    1, // 默认启用
	}

	if err := s.functionDao.Create(function); err != nil {
		return nil, err
	}

	return function, nil
}

// CheckNameExists 检查名称是否已存在（在同一业务组内）
func (s *APIBusinessFunctionServiceImpl) CheckNameExists(groupID uint, name string, excludeID uint) (bool, error) {
	function, err := s.functionDao.GetByGroupIDAndName(groupID, name)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return false, err
	}

	// 如果找到记录，检查是否排除ID
	if function != nil {
		if excludeID > 0 && function.ID == excludeID {
			return false, nil // 排除当前记录
		}
		return true, nil // 存在同名记录
	}

	return false, nil
}
