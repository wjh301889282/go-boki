package global

import (
	"github.com/go-redis/redis" // 导入 Redis 客户端库
	"gorm.io/gorm"              // 导入 GORM ORM 库，用于数据库操作
)

// 全局变量，用于存储数据库和 Redis 的连接实例
var (
	// Db 是 GORM 的数据库连接实例，能够用来与数据库进行交互
	Db *gorm.DB
	// RedisDB 是 Redis 的客户端连接实例，用于与 Redis 进行交互
	RedisDB   *redis.Client
	JwtSecret = []byte("key") // 可以替换成你自己的秘钥
)
