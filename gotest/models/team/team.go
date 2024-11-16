package team

import (
	"exchangeapp/models/user"
	"gorm.io/gorm"
)

// Team 团队模型，存储团队信息
type Team struct {
	gorm.Model
	Name        string       `gorm:"type:varchar(255);not null;uniqueIndex:idx_owner_name"` // 与 OwnerID 组合唯一
	Description string       `gorm:"type:text"`                                             // 团队描述信息
	OwnerID     uint         `gorm:"not null;uniqueIndex:idx_owner_name"`                   // 与 Name 组合唯一
	Owner       user.User    `gorm:"foreignKey:OwnerID"`                                    // 团队拥有者，外键
	Members     []TeamMember `gorm:"foreignKey:TeamID"`
}
