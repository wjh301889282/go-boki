package models

import "gorm.io/gorm"

type User struct {
	gorm.Model
	Username string  `gorm:"unique"` // 用户名
	Password string  // 密码
	Email    *string `gorm:"unique"`        // 邮箱，支持邮箱登录
	Phone    *string `gorm:"unique"`        // 电话号码，支持电话登录
	QQ       *string `gorm:"unique"`        // QQ，支持 QQ 登录
	WeChat   *string `gorm:"unique"`        // 微信，支持微信登录
	Level    int     `gorm:"default:1"`     // 账号等级，默认等级为 1
	IsBanned bool    `gorm:"default:false"` // 是否封禁，默认不封禁
	Pkg      *string `gorm:"unique"`        // 微信，支持微信登录
}
