package gorm

import (
	"exchangeapp/global"
	"exchangeapp/models/user"
	"fmt"
	"log"
	"reflect"
)

func InitGORM() {
	entities := []interface{}{
		&user.User{},
		// 更多结构体
	}

	var failedEntities []string

	for _, entity := range entities {
		entityName := reflect.TypeOf(entity).Elem().Name()

		if err := global.Db.AutoMigrate(entity); err != nil {
			log.Printf("迁移失败: %s, 错误: %v", entityName, err)
			failedEntities = append(failedEntities, entityName)
			continue
		}

		fmt.Printf("成功迁移: %s\n", entityName)
	}

	if len(failedEntities) > 0 {
		log.Fatalf("以下结构体迁移失败: %v", failedEntities)
	}

	fmt.Println("数据库连接和所有结构体自动迁移成功!")
}
