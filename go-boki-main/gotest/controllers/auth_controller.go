package controllers

import (
	"errors"
	"exchangeapp/global" // 引入全局包，用于访问数据库实例
	"exchangeapp/models/user"
	"exchangeapp/rsp"   // 引入错误处理包
	"exchangeapp/utils" // 引入工具包，处理密码加密、JWT 生成等操作
	"fmt"
	"github.com/gin-gonic/gin" // 引入 Gin 框架，用于处理 Web 请求和响应
	"github.com/go-redis/redis"
	"net/http" // 引入 HTTP 包，用于定义 HTTP 状态码
)

// Register 处理用户注册请求
// @Summary 用户注册接口
// @Description 用户通过用户名和密码注册
// @Tags Auth
// @Accept json
// @Produce json
// @Success 200 {string} string "注册成功，返回 JWT token"
// @Failure 400 {string} string "请求参数错误"
// @Failure 500 {string} string "服务器内部错误"
// @Router /api/auth/register [post]
func Register(ctx *gin.Context) {
	type RegisterRequest struct {
		Account struct {
			ID       *uint  `json:"id"`
			Type     int    `json:"type"`
			Email    string `json:"email"`
			Username string `json:"username"`
			Password string `json:"password"`
		} `json:"account"`
		Captcha string `json:"captcha"` // 添加 captcha 字段
	}

	var req RegisterRequest

	// 绑定 JSON 数据到 req 结构体
	if err := ctx.ShouldBindJSON(&req); err != nil {
		// 使用 panic 抛出错误，传入错误信息、请求数据
		panic(rsp.NewErrorResponse(30001, err.Error(), req))
	}

	// 从 req.Account 中提取用户数据
	user := user.User{
		Email:    &req.Account.Email,
		Username: req.Account.Username,
		Password: req.Account.Password,
	}

	// 验证邮箱格式
	if &user.Email == nil || !utils.IsValidEmail(*user.Email) {
		// 使用 panic 抛出错误，传入错误信息、请求数据
		panic(rsp.NewErrorResponse(40002, "邮箱格式无效", req))
	}

	// 检查邮箱是否已注册
	if err := global.Db.Where("email = ?", user.Email).First(&user).Error; err == nil {
		// 使用 panic 抛出错误，传入错误信息、请求数据
		panic(rsp.NewErrorResponse(40003, "邮箱已注册", req))
	}

	// 从 Redis 中获取验证码
	redisKey := fmt.Sprintf("verification:%s", *user.Email)
	storedCode, err := global.RedisDB.Get(redisKey).Result()

	fmt.Println(redisKey)
	fmt.Println(storedCode)
	if errors.Is(err, redis.Nil) {
		// Redis 中没有该验证码，说明验证码过期或不存在
		panic(rsp.NewErrorResponse(40005, "验证码已过期或无效", redis.Nil))
	} else if err != nil {
		// Redis 查询错误
		panic(rsp.NewErrorResponse(50003, "验证码验证失败", redis.Nil))
	}

	// 比较验证码
	if storedCode != req.Captcha {
		// 使用 panic 抛出错误，传入错误信息、请求数据
		panic(rsp.NewErrorResponse(40006, "验证码错误", req))
	}

	// 密码加密
	hashedPwd, err := utils.HashPassword(user.Password)
	if err != nil {
		// 使用 panic 抛出错误，传入错误信息、请求数据
		panic(rsp.NewErrorResponse(50002, err.Error(), hashedPwd))
	}
	user.Password = hashedPwd

	// 生成 JWT
	token, err := utils.GenerateJWT(user.Username)
	if err != nil {
		panic(rsp.NewErrorResponse(50001, err.Error(), token))
	}

	// 插入用户数据
	if err := global.Db.Create(&user).Error; err != nil {
		panic(rsp.NewErrorResponse(20003, err.Error(), user))
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
		Username string `json:"username" binding:"required"` // 必填字段
		Password string `json:"password" binding:"required"` // 必填字段
	}

	// 绑定 JSON 数据到 input 结构体
	if err := ctx.ShouldBindJSON(&input); err != nil {
		panic(rsp.NewErrorResponse(30001, err.Error(), input)) // 使用 panic 抛出错误
	}

	var user user.User

	// 从数据库中查询用户
	if err := global.Db.Where("username = ?", input.Username).First(&user).Error; err != nil {
		panic(rsp.NewErrorResponse(10002, err.Error(), input)) // 使用 panic 抛出错误
	}

	// 校验密码
	if !utils.CheckPassword(input.Password, user.Password) {
		panic(rsp.NewErrorResponse(10003, input.Password, input)) // 使用 panic 抛出错误
	}

	// 生成 JWT
	token, err := utils.GenerateJWT(user.Username)
	if err != nil {
		panic(rsp.NewErrorResponse(50001, err.Error(), input)) // 使用 panic 抛出错误
	}

	// 返回成功响应，包含生成的 JWT token
	ctx.JSON(http.StatusOK, gin.H{"token": token})
}
