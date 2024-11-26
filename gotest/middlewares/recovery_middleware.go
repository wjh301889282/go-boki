package middlewares

import (
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

// RecoveryMiddleware 捕获并处理 panic 的中间件
func RecoveryMiddleware() gin.HandlerFunc {
	// 打开日志文件（如果不存在则创建）
	logFile, err := os.OpenFile("panic_errors.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		panic("无法打开或创建日志文件")
	}

	// 创建一个日志记录器
	logger := log.New(logFile, "[Recovery] ", log.LstdFlags|log.Lshortfile)

	return func(ctx *gin.Context) {
		defer func() {
			// 捕获 panic
			if err := recover(); err != nil {
				// 将错误写入日志文件
				logger.Printf("捕获到 panic: %v", err)
				log.Printf("[Recovery] 捕获到 panic:\n%s\n", err)
				// 返回统一的错误响应
				ctx.JSON(http.StatusInternalServerError, gin.H{
					"detail": err, // 可选：生产环境可以隐藏具体错误信息
				})

				// 阻止后续的 Handler 执行
				ctx.Abort()
			}
		}()

		// 执行下一个中间件或 Handler
		ctx.Next()
	}
}
