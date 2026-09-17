// 基础 Skill 实现
package skill

// BaseSkill 基础技能实现（提供默认方法）
type BaseSkill struct {
	name        string
	description string
	category    SkillCategory
}

func (b *BaseSkill) Name() string {
	return b.name
}

func (b *BaseSkill) Description() string {
	return b.description
}

func (b *BaseSkill) Category() SkillCategory {
	return b.category
}

// Validate 默认验证实现（子类可覆盖）
func (b *BaseSkill) Validate(req *SkillRequest) error {
	return nil
}
