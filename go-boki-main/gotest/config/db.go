package config

import (
	"exchangeapp/global" // 导入全局包，使用全局变量 Db 保存数据库连接
	"log"                // 导入日志包，用于记录错误信息
	"time"               // 导入 time 包，用于设置数据库连接的最大生命周期

	"gorm.io/driver/mysql" // 导入 GORM MySQL 驱动
	"gorm.io/gorm"         // 导入 GORM 库，用于数据库操作
)

// initDB 初始化数据库连接并设置相关配置
func initDB() {
	// 从 AppConfig 中获取数据库的 DSN（数据源名称）
	dsn := AppConfig.Database.Dsn

	// 使用 GORM 连接到 MySQL 数据库
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		// 如果数据库连接失败，输出错误并终止程序
		log.Fatalf("数据库连接失败: %v", err)
	}

	// 获取底层的原生数据库连接
	sqlDB, err := db.DB()
	if err != nil {
		// 获取底层数据库连接失败，输出错误并终止程序
		log.Fatalf("数据库连接失败: %v", err)
	}

	// 设置数据库连接池的最大空闲连接数
	sqlDB.SetMaxIdleConns(AppConfig.Database.MaxIdleConns)

	// 设置数据库连接池的最大打开连接数
	sqlDB.SetMaxOpenConns(AppConfig.Database.MaxOpenConns)

	// 设置数据库连接的最大生命周期，超过该时间后，连接会被关闭
	sqlDB.SetConnMaxLifetime(time.Hour)

	// 将数据库连接实例赋值给全局变量 global.Db，供其他部分使用
	global.Db = db
}
