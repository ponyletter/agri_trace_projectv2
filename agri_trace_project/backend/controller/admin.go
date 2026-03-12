package controller

import (
	"agri-trace/model"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// AdminController 管理员控制器
type AdminController struct {
	DB *gorm.DB
}

// ==================== 用户管理 ====================

// ListUsersRequest 用户列表查询参数
type ListUsersRequest struct {
	Page     int    `form:"page"`
	PageSize int    `form:"page_size"`
	Keyword  string `form:"keyword"`
	Role     string `form:"role"`
}

// ListUsers 获取用户列表（管理员）
func (ctrl *AdminController) ListUsers(c *gin.Context) {
	var req ListUsersRequest
	c.ShouldBindQuery(&req)
	if req.Page <= 0 {
		req.Page = 1
	}
	if req.PageSize <= 0 {
		req.PageSize = 10
	}

	query := ctrl.DB.Model(&model.User{})
	if req.Keyword != "" {
		query = query.Where("username LIKE ? OR real_name LIKE ?", "%"+req.Keyword+"%", "%"+req.Keyword+"%")
	}
	if req.Role != "" {
		query = query.Where("role = ?", req.Role)
	}

	var total int64
	query.Count(&total)

	var users []model.User
	offset := (req.Page - 1) * req.PageSize
	if err := query.Offset(offset).Limit(req.PageSize).Order("id ASC").Find(&users).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "msg": "查询失败"})
		return
	}

	// 脱敏：不返回密码
	type UserVO struct {
		ID        uint      `json:"id"`
		Username  string    `json:"username"`
		RealName  string    `json:"real_name"`
		Role      string    `json:"role"`
		Phone     string    `json:"phone"`
		Email     string    `json:"email"`
		IsActive  bool      `json:"is_active"`
		CreatedAt time.Time `json:"created_at"`
	}
	vos := make([]UserVO, 0, len(users))
	for _, u := range users {
		vos = append(vos, UserVO{
			ID: u.ID, Username: u.Username, RealName: u.RealName,
			Role: u.Role, Phone: u.Phone, Email: u.Email,
			IsActive: u.IsActive, CreatedAt: u.CreatedAt,
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"data": gin.H{"list": vos, "total": total},
	})
}

// CreateUserRequest 创建用户请求
type CreateUserRequest struct {
	Username string `json:"username" binding:"required,min=3,max=20"`
	Password string `json:"password" binding:"required,min=6"`
	RealName string `json:"real_name" binding:"required"`
	Role     string `json:"role" binding:"required"`
	Phone    string `json:"phone"`
	Email    string `json:"email"`
}

// CreateUser 管理员创建用户
func (ctrl *AdminController) CreateUser(c *gin.Context) {
	var req CreateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "msg": err.Error()})
		return
	}

	// 检查用户名唯一
	var count int64
	ctrl.DB.Model(&model.User{}).Where("username = ?", req.Username).Count(&count)
	if count > 0 {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "msg": "用户名已存在"})
		return
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "msg": "密码加密失败"})
		return
	}

	user := model.User{
		Username:     req.Username,
		PasswordHash: string(hash),
		RealName:     req.RealName,
		Role:         req.Role,
		Phone:        req.Phone,
		Email:        req.Email,
		IsActive:     true,
	}
	if err := ctrl.DB.Create(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "msg": "创建失败"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"code": 200, "msg": "用户创建成功", "data": gin.H{"id": user.ID}})
}

// UpdateUserRequest 更新用户请求
type UpdateUserRequest struct {
	RealName string `json:"real_name"`
	Role     string `json:"role"`
	Phone    string `json:"phone"`
	Email    string `json:"email"`
	IsActive *bool  `json:"is_active"`
	Password string `json:"password"`
}

// UpdateUser 更新用户信息
func (ctrl *AdminController) UpdateUser(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "msg": "无效ID"})
		return
	}

	var req UpdateUserRequest
	c.ShouldBindJSON(&req)

	updates := map[string]interface{}{}
	if req.RealName != "" {
		updates["real_name"] = req.RealName
	}
	if req.Role != "" {
		updates["role"] = req.Role
	}
	if req.Phone != "" {
		updates["phone"] = req.Phone
	}
	if req.Email != "" {
		updates["email"] = req.Email
	}
	if req.IsActive != nil {
		updates["is_active"] = *req.IsActive
	}
	if req.Password != "" {
		hash, _ := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
		updates["password_hash"] = string(hash)
	}

	if err := ctrl.DB.Model(&model.User{}).Where("id = ?", id).Updates(updates).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "msg": "更新失败"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"code": 200, "msg": "更新成功"})
}

// DeleteUser 删除用户
func (ctrl *AdminController) DeleteUser(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "msg": "无效ID"})
		return
	}
	if err := ctrl.DB.Delete(&model.User{}, id).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "msg": "删除失败"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"code": 200, "msg": "删除成功"})
}

// ==================== 电子合格证创建 ====================

// CreateCertificateRequest 创建合格证请求
type CreateCertificateRequest struct {
	BatchID      uint   `json:"batch_id" binding:"required"`
	ProducerName string `json:"producer_name" binding:"required"`
	ProducerPhone string `json:"producer_phone"`
	ProducerAddr  string `json:"producer_addr"`
	InspectOrg   string `json:"inspect_org"`
	PesticideOk  bool   `json:"pesticide_ok"`
	HeavyMetalOk bool   `json:"heavy_metal_ok"`
	MicrobeOk    bool   `json:"microbe_ok"`
	ValidDays    int    `json:"valid_days"`
}

// CreateCertificate 创建电子合格证
func (ctrl *AdminController) CreateCertificate(c *gin.Context) {
	var req CreateCertificateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "msg": err.Error()})
		return
	}

	// 查批次
	var batch model.AgriBatch
	if err := ctrl.DB.First(&batch, req.BatchID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"code": 404, "msg": "批次不存在"})
		return
	}

	validDays := req.ValidDays
	if validDays <= 0 {
		validDays = 365
	}

	certNo := fmt.Sprintf("CERT-AKS-%s-%d", time.Now().Format("20060102"), batch.ID)
	now := time.Now()

	cert := model.Certificate{
		BatchID:      req.BatchID,
		CertNo:       certNo,
		ProductName:  batch.ProductName,
		ProducerName: req.ProducerName,
		ProducerPhone: req.ProducerPhone,
		ProducerAddr:  req.ProducerAddr,
		InspectOrg:   req.InspectOrg,
		PesticideOk:  req.PesticideOk,
		HeavyMetalOk: req.HeavyMetalOk,
		MicrobeOk:    req.MicrobeOk,
		IssueDate:    now.Format("2006-01-02"),
		ValidUntil:   now.AddDate(0, 0, validDays).Format("2006-01-02"),
		IssueOrg:     "新疆阿克苏市农业农村局",
	}

	if err := ctrl.DB.Create(&cert).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "msg": "合格证创建失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"code": 200, "msg": "电子合格证创建成功", "data": cert})
}
