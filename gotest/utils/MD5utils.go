package utils

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"net/url"
	"sort"
)

// EpusdtSign MD加密
func EpusdtSign(params map[string]interface{}, signKey string) string {
	// 1. 参数排序
	keys := make([]string, 0, len(params))
	for k := range params {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	// 2. 组装签名字符串和 URL 编码字符串
	sign := ""
	urls := ""

	for _, key := range keys {
		value := params[key]
		if strVal, ok := value.(string); ok && strVal == "" {
			continue
		}
		if key != "signature" {
			if sign != "" {
				sign += "&"
				urls += "&"
			}
			sign += key + "=" + fmt.Sprintf("%v", value)
			urls += key + "=" + url.QueryEscape(fmt.Sprintf("%v", value))
		}
	}
	fmt.Println(urls)
	// 3. 追加密钥并生成 MD5 签名
	signWithKey := sign + signKey
	hasher := md5.New()
	hasher.Write([]byte(signWithKey))
	md5Sign := hex.EncodeToString(hasher.Sum(nil))
	fmt.Println(md5Sign)
	return md5Sign
}
