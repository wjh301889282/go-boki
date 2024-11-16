package team

import (
	"exchangeapp/models/user"
	"gorm.io/gorm"
)

// TeamMember 团队成员模型，存储团队成员信息以及权限
type TeamMember struct {
	gorm.Model
	TeamID      uint      `gorm:"not null"`                                             // 团队ID，外键
	UserID      uint      `gorm:"not null"`                                             // 用户ID，外键
	Role        string    `gorm:"type:enum('owner','admin','member');default:'member'"` // 成员角色
	Permissions string    `gorm:"type:text"`                                            // 权限（JSON 格式存储）
	Team        Team      `gorm:"foreignKey:TeamID"`                                    // 关联的团队
	User        user.User `gorm:"foreignKey:UserID"`                                    // 关联的用户
}
