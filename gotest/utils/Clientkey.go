package utils

import (
	"crypto/rand"
	"math/big"
	"strings"
)

// GenerateClientKey 生成包含字母和数字的随机密钥
func GenerateClientKey() (string, error) {
	// 可选的字符集，包括大写字母、小写字母和数字
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

	// 密钥长度为 32 字符
	keyLength := 45
	var sb strings.Builder
	sb.Grow(keyLength)

	// 使用 crypto/rand 生成安全的随机字节
	for i := 0; i < keyLength; i++ {
		// 随机从字符集选择一个字符
		randomByte, err := rand.Int(rand.Reader, big.NewInt(int64(len(charset))))
		if err != nil {
			return "", err
		}
		sb.WriteByte(charset[randomByte.Int64()])
	}

	return sb.String(), nil
}
