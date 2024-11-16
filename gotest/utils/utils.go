package utils

import (
	"errors"
	"exchangeapp/global"
	"time"

	"github.com/golang-jwt/jwt"
	"golang.org/x/crypto/bcrypt"
)

// HashPassword 用于加密密码，返回加密后的哈希值
func HashPassword(pwd string) (string, error) {
	// 使用 bcrypt 生成加密后的密码哈希，成本因子为 12
	hash, err := bcrypt.GenerateFromPassword([]byte(pwd), 12)
	// 返回加密后的密码和可能发生的错误
	return string(hash), err
}

// GenerateJWT 用于生成一个带有用户名和过期时间（72小时）的 JWT
func GenerateJWT(username string) (string, error) {
	// 创建一个新的 JWT，使用 HS256 算法，payload 包含用户名和过期时间
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"username": username,                              // 用户名
		"exp":      time.Now().Add(time.Hour * 72).Unix(), // 设置 JWT 过期时间为当前时间的 72 小时后
	})

	// 使用 secret 密钥对 token 进行签名，返回签名后的 JWT 字符串
	signedToken, err := token.SignedString([]byte(global.JwtSecret))
	// 返回带有 "Bearer " 前缀的 JWT 和可能发生的错误
	return signedToken, err
}

// CheckPassword 用于验证输入的密码与存储的哈希密码是否匹配
func CheckPassword(password string, hash string) bool {
	// 比较输入的密码和存储的哈希密码是否一致
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	// 如果密码匹配，err 会是 nil，返回 true；否则返回 false
	return err == nil
}

// ParseJWT 用于解析 JWT，并从中提取用户名
func ParseJWT(tokenString string) (string, error) {
	// 检查 token 字符串是否以 "Bearer " 开头，若是则去掉 "Bearer " 前缀
	//if len(tokenString) > 7 && tokenString[:7] == "Bearer " {
	//	tokenString = tokenString[7:]
	//}

	// 解析 JWT，并验证签名方法是否是 HMAC（HS256）
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// 检查签名方法是否是 HMAC
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected Signing Method") // 如果签名方法不正确，返回错误
		}
		// 返回用于签名的密钥（这里是 "secret"）
		return []byte(global.JwtSecret), nil
	})

	// 如果解析失败，返回错误
	if err != nil {
		return "", err
	}

	// 如果 JWT 有效，提取其中的 claims（负载数据）
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		// 获取用户名
		username, ok := claims["username"].(string)
		// 如果用户名不是字符串类型，返回错误
		if !ok {
			return "", errors.New("username claim is not a string")
		}
		// 返回用户名
		return username, nil
	}

	// 如果 JWT 无效，返回错误
	return "", err
}
