package config

import (
	"exchangeapp/global" // 导入全局包，使用全局变量 RedisDB 保存 Redis 客户端连接
	"log"                // 导入日志包，用于记录错误信息

	"github.com/go-redis/redis" // 导入 Redis 客户端库，用于操作 Redis
)

// initRedis 初始化 Redis 客户端并进行连接
func initRedis() {
	// 从配置中获取 Redis 服务器的地址、数据库索引和密码
	addr := AppConfig.Redis.Addr
	db := AppConfig.Redis.DB
	password := AppConfig.Redis.Password

	// 创建一个新的 Redis 客户端实例，传入配置项
	RedisClient := redis.NewClient(&redis.Options{
		Addr:     addr,     // Redis 服务器的地址
		DB:       db,       // Redis 使用的数据库索引
		Password: password, // Redis 连接密码
	})

	// 使用 Ping 命令测试是否能成功连接到 Redis 服务器
	_, err := RedisClient.Ping().Result()
	if err != nil {
		// 如果连接失败，输出错误日志并终止程序
		log.Fatalf("Failed to connect to Redis, got error: %v", err)
	}

	// 将 Redis 客户端实例赋值给全局变量 global.RedisDB，供其他部分使用
	global.RedisDB = RedisClient
}
