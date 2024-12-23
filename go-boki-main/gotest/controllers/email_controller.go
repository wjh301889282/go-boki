package controllers

import (
	"exchangeapp/global"
	"exchangeapp/models/user"
	"exchangeapp/rsp"
	"exchangeapp/utils"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

// SendVerificationCode 发送邮件
func SendVerificationCode(ctx *gin.Context) {
	type Request struct {
		Email string `json:"email" binding:"required,email"`
	}

	var req Request

	// 解析请求体
	if err := ctx.ShouldBindJSON(&req); err != nil {
		panic(rsp.NewErrorResponse(40001, "请求参数错误: "+err.Error(), req))
	}

	// 验证邮箱格式
	if &req.Email == nil || !utils.IsValidEmail(req.Email) {
		panic(rsp.NewErrorResponse(40002, "邮箱格式无效", req))
	}

	// 检查邮箱是否已注册
	var existingUser user.User
	if err := global.Db.Where("email = ?", req.Email).First(&existingUser).Error; err == nil {
		panic(rsp.NewErrorResponse(40002, err.Error(), req))
	}

	fmt.Println("开始生成验证码")
	// 生成验证码
	code := utils.GenerateVerificationCode()

	fmt.Println("开始存入redis")
	// 将验证码存入 Redis，设置 300 秒过期时间
	redisKey := fmt.Sprintf("verification:%s", req.Email)
	err := global.RedisDB.Set(redisKey, code, time.Minute*5).Err()
	if err != nil {
		panic(rsp.NewErrorResponse(50001, err.Error(), req))
	}

	// 发送验证码到邮箱
	err = utils.SendEmail(req.Email, "注册验证码", fmt.Sprintf("您的验证码是：%s，有效期为 5 分钟", code))
	if err != nil {
		panic(rsp.NewErrorResponse(50002, err.Error(), req))
	}

	// 返回响应
	ctx.JSON(http.StatusOK, gin.H{"message": "验证码发送成功，请检查您的邮箱"})
}
