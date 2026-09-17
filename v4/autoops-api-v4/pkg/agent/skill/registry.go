// Skill 注册器
package skill

import (
	"sync"

	"dodevops-api/pkg/log"
)

// SkillRegistry Skill注册中心
type SkillRegistry struct {
	skills map[string]Skill
	mu     sync.RWMutex
}

var defaultRegistry = NewSkillRegistry()

// GetDefaultRegistry 获取默认注册器
func GetDefaultRegistry() *SkillRegistry {
	return defaultRegistry
}

// NewSkillRegistry 创建新的注册器
func NewSkillRegistry() *SkillRegistry {
	return &SkillRegistry{
		skills: make(map[string]Skill),
	}
}

// Register 注册技能
func (r *SkillRegistry) Register(skill Skill) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.skills[skill.Name()] = skill
	log.Log().Infof("[SkillRegistry] 注册技能: %s - %s", skill.Name(), skill.Description())
}

// Find 根据意图查找合适的技能
func (r *SkillRegistry) Find(intent *Intent) Skill {
	r.mu.RLock()
	defer r.mu.RUnlock()

	for _, skill := range r.skills {
		if skill.CanHandle(intent) {
			log.Log().Infof("[SkillRegistry] 找到匹配技能: %s", skill.Name())
			return skill
		}
	}
	return nil
}

// GetAll 获取所有技能
func (r *SkillRegistry) GetAll() []Skill {
	r.mu.RLock()
	defer r.mu.RUnlock()

	skills := make([]Skill, 0, len(r.skills))
	for _, skill := range r.skills {
		skills = append(skills, skill)
	}
	return skills
}

// GetByName 根据名称获取技能
func (r *SkillRegistry) GetByName(name string) Skill {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return r.skills[name]
}
