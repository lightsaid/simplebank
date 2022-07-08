package utils

import (
	"math/rand"
	"strings"
	"time"
)

const chars = "qwertyuiopasdfghjklzxcvbnm1234567890"

// 初始化随机种子
func init() {
	rand.Seed(time.Now().UnixNano())
}

// RandomInt 生成一个随机数字，范围（min~max）
func RandomInt(min, max int64) int64 {
	return min + rand.Int63n(max-min+1)
}

// RandomString 生成随机字符串, 长度 = n
func RandomString(n int) string {
	var sb strings.Builder
	k := len(chars)
	for i := 0; i < n; i++ {
		c := chars[rand.Intn(k)]
		sb.WriteByte(c)
	}
	return sb.String()
}

// RandomOwner 随机一个 account 表 owner 字段
func RandomOwner() string {
	return RandomString(6)
}

// RandomMoney 随机一个balance提供给account表
func RandomMoney() int64 {
	return RandomInt(0, 1000)
}

// RandomCurrency 随机一种货币
func RandomCurrency() string {
	currencies := []string{"RMB", "USD", "EUR"}
	n := len(currencies)
	return currencies[rand.Intn(n)]
}
