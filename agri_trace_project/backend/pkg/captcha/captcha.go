package captcha

import (
	"math/rand"
	"time"

	"github.com/google/uuid"
)

// GenerateCaptcha 生成一个简单的数字验证码（实际项目中可使用 base64Captcha 库）
// 为了减少依赖，这里使用简单的随机数字字符串模拟
func GenerateCaptcha() (string, string) {
	rand.Seed(time.Now().UnixNano())
	const charset = "0123456789"
	b := make([]byte, 4)
	for i := range b {
		b[i] = charset[rand.Intn(len(charset))]
	}
	key := uuid.New().String()
	return key, string(b)
}
