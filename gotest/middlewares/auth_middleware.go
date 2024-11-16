package middlewares

import (
	"exchangeapp/utils" // 导入 utils 包，处理 JWT 解析
	"net/http"          // 导入 net/http 包，用于 HTTP 状态码

	"github.com/gin-gonic/gin" // 导入 Gin 框架
)

// AuthMiddleWare 返回一个 Gin 中间件，用于处理身份验证
func AuthMiddleWare() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// 从请求头中获取 Authorization 字段
		token := ctx.GetHeader("Authorization")

		// 如果没有提供 Authorization Header，则返回 401 未授权错误
		if token == "" {
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Missing Authorization Header"})
			ctx.Abort() // 中止当前请求的处理
			return
		}

		// 解析 JWT，获取用户名和验证 token 是否有效
		username, err := utils.ParseJWT(token)

		// 如果解析失败，返回 401 未授权错误
		if err != nil {
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			ctx.Abort() // 中止当前请求的处理
			return
		}

		// 如果 token 有效，将用户名存入上下文中
		ctx.Set("username", username)

		// 继续处理请求
		ctx.Next()
	}
}
