package controller

import (
	"agri-trace/middleware"
	"agri-trace/model"
	"agri-trace/pkg/captcha"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type AuthController struct {
	DB *gorm.DB
}

// GetCaptcha 获取验证码
func (ctrl *AuthController) GetCaptcha(c *gin.Context) {
	key, code := captcha.GenerateCaptcha()
	
	// 保存到数据库
	cap := model.Captcha{
		Key:       key,
		Code:      code,
		ExpiredAt: time.Now().Add(5 * time.Minute),
	}
	ctrl.DB.Create(&cap)

	// 为了简化前端对接，这里直接返回验证码文本（实际应返回base64图片）
	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"data": gin.H{
			"key":          key,
			"captcha_code": code, // 模拟图片内容
		},
	})
}

type LoginRequest struct {
	Username   string `json:"username" binding:"required"`
	Password   string `json:"password" binding:"required"`
	CaptchaKey string `json:"captcha_key"`
	Captcha    string `json:"captcha"`
}

func (ctrl *AuthController) Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "msg": "请求参数错误"})
		return
	}

	// 校验验证码 (非空时校验)
	if req.CaptchaKey != "" && req.Captcha != "" {
		var cap model.Captcha
		if err := ctrl.DB.Where("`key` = ? AND code = ? AND used = 0 AND expired_at > ?", req.CaptchaKey, req.Captcha, time.Now()).First(&cap).Error; err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"code": 400, "msg": "验证码错误或已过期"})
			return
		}
		ctrl.DB.Model(&cap).Update("used", 1)
	}

	var user model.User
	if err := ctrl.DB.Where("username = ?", req.Username).First(&user).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"code": 401, "msg": "用户名或密码错误"})
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password)); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"code": 401, "msg": "用户名或密码错误"})
		return
	}

	token, err := middleware.GenerateToken(user.ID, user.Username, user.Role)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "msg": "令牌生成失败"})
		return
	}

	// 更新最后登录时间
	ctrl.DB.Model(&user).Update("last_login_at", time.Now())

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  "登录成功",
		"data": gin.H{
			"token":     token,
			"user_id":   user.ID,
			"username":  user.Username,
			"real_name": user.RealName,
			"role":      user.Role,
			"avatar":    user.Avatar,
		},
	})
}

type RegisterRequest struct {
	Username   string `json:"username" binding:"required,min=3,max=64"`
	Password   string `json:"password" binding:"required,min=6"`
	RealName   string `json:"real_name" binding:"required"`
	Role       string `json:"role" binding:"required"`
	Phone      string `json:"phone"`
	CaptchaKey string `json:"captcha_key" binding:"required"`
	Captcha    string `json:"captcha" binding:"required"`
}

func (ctrl *AuthController) Register(c *gin.Context) {
	var req RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "msg": err.Error()})
		return
	}

	// 校验验证码
	var cap model.Captcha
	if err := ctrl.DB.Where("`key` = ? AND code = ? AND used = 0 AND expired_at > ?", req.CaptchaKey, req.Captcha, time.Now()).First(&cap).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "msg": "验证码错误或已过期"})
		return
	}
	ctrl.DB.Model(&cap).Update("used", 1)

	var count int64
	ctrl.DB.Model(&model.User{}).Where("username = ?", req.Username).Count(&count)
	if count > 0 {
		c.JSON(http.StatusConflict, gin.H{"code": 409, "msg": "用户名已存在"})
		return
	}

	hash, _ := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)

	user := model.User{
		Username:     req.Username,
		PasswordHash: string(hash),
		RealName:     req.RealName,
		Role:         req.Role,
		Phone:        req.Phone,
		Status:       1,
	}
	if err := ctrl.DB.Create(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "msg": "用户创建失败"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"code": 200, "msg": "注册成功"})
}

func (ctrl *AuthController) GetProfile(c *gin.Context) {
	userID, _ := c.Get("user_id")
	var user model.User
	if err := ctrl.DB.First(&user, userID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"code": 404, "msg": "用户不存在"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"code": 200, "data": user})
}
