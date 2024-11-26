package utils

import (
	"fmt"
	"gopkg.in/gomail.v2"
	"math/rand"
	"regexp"
	"time"
)

// GenerateVerificationCode 随机生成一个6位数的验证码。
func GenerateVerificationCode() string {
	rand.Seed(time.Now().UnixNano())
	code := fmt.Sprintf("%06d", rand.Intn(1000000))
	return code
}

// IsValidEmail 验证邮箱格式
func IsValidEmail(email string) bool {
	regex := `^[a-z0-9._%+-]+@[a-z0-9.-]+\.[a-z]{2,}$`
	re := regexp.MustCompile(regex)
	return re.MatchString(email)
}

// SendEmail 发送邮件
func SendEmail(to, subject, body string) error {
	fromEmail := "svj04us_c4t0bak@a.web3woolbox.com"
	fromPassword := "Cgy88888888"
	smtpServer := "mail.a.web3woolbox.com"
	smtpPort := 25

	// 创建邮件消息
	m := gomail.NewMessage()
	m.SetHeader("From", fromEmail)
	m.SetHeader("To", to)
	m.SetHeader("Subject", subject)
	m.SetBody("text/plain", body)

	// 配置 SMTP 连接
	d := gomail.NewDialer(smtpServer, smtpPort, fromEmail, fromPassword)

	// 发送邮件并打印详细错误
	if err := d.DialAndSend(m); err != nil {
		fmt.Printf("邮件发送失败: %v\n", err)
		return err
	}
	fmt.Println("邮件发送成功")

	return nil
}
