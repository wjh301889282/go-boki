package controllers

import (
	"exchangeapp/global" // 引入全局包，用于访问数据库实例
	"exchangeapp/models/user"
	"exchangeapp/rsp"          // 引入错误处理包
	"exchangeapp/utils"        // 引入工具包，处理密码加密、JWT 生成等操作
	"github.com/gin-gonic/gin" // 引入 Gin 框架，用于处理 Web 请求和响应
	"net/http"                 // 引入 HTTP 包，用于定义 HTTP 状态码
)

// Register 处理用户注册请求
// @Summary 用户注册接口
// @Description 用户通过用户名和密码注册
// @Tags Auth
// @Accept json
// @Produce json
// @Param body body models.User true "用户信息"
// @Success 200 {string} string "注册成功，返回 JWT token"
// @Failure 400 {string} string "请求参数错误"
// @Failure 500 {string} string "服务器内部错误"
// @Router /api/auth/register [post]
func Register(ctx *gin.Context) {
	var user user.User

	// 绑定 JSON 数据到 user 结构体
	if err := ctx.ShouldBindJSON(&user); err != nil {
		ctx.JSON(http.StatusBadRequest, rsp.NewErrorResponse(30001, err.Error()))
		return
	}

	// 密码加密
	hashedPwd, err := utils.HashPassword(user.Password)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, rsp.NewErrorResponse(50002, err.Error()))
		return
	}

	user.Password = hashedPwd

	// 生成 JWT
	token, err := utils.GenerateJWT(user.Username)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, rsp.NewErrorResponse(50001, err.Error()))
		return
	}

	// 自动迁移数据库
	if err := global.Db.AutoMigrate(&user); err != nil {
		ctx.JSON(http.StatusInternalServerError, rsp.NewErrorResponse(20004, err.Error()))
		return
	}

	// 插入用户数据
	if err := global.Db.Create(&user).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, rsp.NewErrorResponse(20003, err.Error()))
		return
	}

	// 返回生成的 JWT
	ctx.JSON(http.StatusOK, gin.H{"token": token})
}

// Login 处理用户登录请求
// @Summary 用户登录接口
// @Description 用户通过用户名和密码登录，返回 JWT token
// @Tags Auth
// @Accept json
// @Produce json
// @Param username body string true "用户名" // 用户名
// @Param password body string true "密码"   // 用户密码
// @Success 200 {string} string "登录成功，返回 JWT token"
// @Failure 400 {string} string "用户名或密码错误"
// @Failure 401 {string} string "用户未授权"
// @Router /api/auth/login [post]
func Login(ctx *gin.Context) {
	var input struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	// 绑定 JSON 数据到 input 结构体
	if err := ctx.ShouldBindJSON(&input); err != nil {
		ctx.JSON(http.StatusBadRequest, rsp.NewErrorResponse(30001, err.Error()))
		return
	}

	var user user.User

	// 从数据库中查询用户
	if err := global.Db.Where("username = ?", input.Username).First(&user).Error; err != nil {
		ctx.JSON(http.StatusUnauthorized, rsp.NewErrorResponse(10002, "用户不存在"))
		return
	}

	// 校验密码
	if !utils.CheckPassword(input.Password, user.Password) {
		ctx.JSON(http.StatusUnauthorized, rsp.NewErrorResponse(10003, "用户名或密码错误"))
		return
	}

	// 生成 JWT
	token, err := utils.GenerateJWT(user.Username)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, rsp.NewErrorResponse(50001, err.Error()))
		return
	}

	// 返回成功响应，包含生成的 JWT token
	ctx.JSON(http.StatusOK, gin.H{"token": token})
}
