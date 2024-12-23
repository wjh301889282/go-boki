package config

import (
	"fmt"
	"log" // 导入日志包，用于输出错误日志

	"github.com/spf13/viper" // 导入 Viper 库，Viper 是一个配置管理工具
)

// Config 结构体定义了应用程序的配置，包括应用信息、数据库配置和 Redis 配置信息
type Config struct {
	App struct {
		Name string // 应用的名称
		Port string // 应用的端口
	}
	Database struct {
		Dsn          string // 数据库的连接字符串
		MaxIdleConns int    // 数据库连接池中的最大空闲连接数
		MaxOpenConns int    // 数据库连接池中的最大打开连接数
	}
	Redis struct {
		Addr     string // Redis 服务器地址
		DB       int    // Redis 数据库索引
		Password string // Redis 密码
	}
}

// AppConfig 是一个全局配置实例，保存从配置文件中读取的配置信息
var AppConfig *Config

// InitConfig 初始化配置文件并加载配置
func InitConfig() {
	// 设置配置文件的名称和类型（配置文件应该为 config.yml）
	viper.SetConfigName("config")
	viper.SetConfigType("yml")
	viper.AddConfigPath("./config") // 配置文件存放的路径

	// 读取配置文件
	if err := viper.ReadInConfig(); err != nil {
		// 如果读取配置文件出错，输出错误并终止程序
		log.Fatalf("Error reading config file: %v", err)
	}

	// 初始化 AppConfig 结构体
	AppConfig = &Config{}

	// 将读取到的配置文件数据反序列化到 AppConfig 结构体中
	if err := viper.Unmarshal(AppConfig); err != nil {
		// 如果解码失败，输出错误并终止程序
		log.Fatalf("Unable to decode into struct: %v", err)
	}

	fmt.Println("开始初始化数据库")
	// 初始化数据库和 Redis 连接
	initDB()
	fmt.Println("开始初始化redis")
	initRedis()
	fmt.Println("redis和mysql都连接成功了！！！")
}
