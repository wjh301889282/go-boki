package router

import (
	"exchangeapp/controllers" // 导入 controllers 包，处理具体的路由逻辑
	"exchangeapp/controllers/TeamManagement"
	"exchangeapp/middlewares" // 导入 middlewares 包，处理请求中间件
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"time" // 导入 time 包，用于设置跨域的最大缓存时间

	"github.com/gin-contrib/cors" // 导入 CORS 中间件
	"github.com/gin-gonic/gin"    // 导入 Gin 框架
)

// SetupRouter 设置并返回 Gin 路由引擎
func SetupRouter() *gin.Engine {
	// 创建一个默认的 Gin 引擎（包括了日志和恢复中间件）
	r := gin.Default()
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	// 使用 CORS 中间件进行跨域设置
	r.Use(cors.New(cors.Config{
		// 允许的源，只有来自 http://localhost:5173 的请求会被接受
		AllowOrigins: []string{"http://localhost:5173"},
		// 允许的 HTTP 方法
		AllowMethods: []string{"GET", "POST", "OPTIONS"},
		// 允许的请求头部
		AllowHeaders: []string{"Origin", "Content-Type", "Authorization"},
		// 允许暴露的响应头部
		ExposeHeaders: []string{"Content-Length"},
		// 允许携带凭证
		AllowCredentials: true,
		// 设置跨域请求的最大缓存时间，12小时
		MaxAge: 12 * time.Hour,
	}))

	// 创建一个路由组，用于认证相关的路由（以 /api/auth 为前缀）
	auth := r.Group("/api/auth")
	{
		// 登录接口，使用 POST 请求
		auth.POST("/login", controllers.Login)
		// 注册接口，使用 POST 请求
		auth.POST("/register", controllers.Register)
	}

	// 创建一个路由组，用于主要的 API 路由（以 /api 为前缀）
	api := r.Group("/api")
	// 获取汇率接口，使用 GET 请求
	api.GET("/exchangeRates", controllers.GetExchangeRates)

	// 使用 AuthMiddleWare 中间件来保护以下接口，需要身份验证
	api.Use(middlewares.AuthMiddleWare())
	{
		// 创建汇率接口，使用 POST 请求
		api.POST("/exchangeRates", controllers.CreateExchangeRate)
		// 创建文章接口，使用 POST 请求
		api.POST("/articles", controllers.CreateArticle)
		// 获取所有文章接口，使用 GET 请求
		api.GET("/articles", controllers.GetArticles)
		// 根据文章 ID 获取单篇文章，使用 GET 请求
		api.GET("/articles/:id", controllers.GetArticleByID)

		// 点赞文章接口，使用 POST 请求
		api.POST("/articles/:id/like", controllers.LikeArticle)
		// 获取文章的点赞数接口，使用 GET 请求
		api.GET("/articles/:id/like", controllers.GetArticleLikes)
	}

	// TeamManagement 路由分组
	teamMg := r.Group("/teamMg")
	{
		// 创建团队
		teamMg.POST("/createteam", TeamManagement.CreateTeam) // 创建新团队
		// 添加用户到团队
		teamMg.POST("/addusertoteam", TeamManagement.AddUserToTeam) // 将用户添加到指定团队，并设置角色和权限
		// 查询团队成员
		teamMg.GET("/team/:id/members", TeamManagement.GetTeamMembers) // 查询指定团队的所有成员信息
	}
	// 返回设置好的路由引擎
	return r
}
