package model

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"time"
)

// JSONMap 自定义JSON类型
type JSONMap map[string]interface{}

func (j JSONMap) Value() (driver.Value, error) {
	if j == nil {
		return nil, nil
	}
	b, err := json.Marshal(j)
	return string(b), err
}

func (j *JSONMap) Scan(value interface{}) error {
	if value == nil {
		*j = nil
		return nil
	}
	var bytes []byte
	switch v := value.(type) {
	case []byte:
		bytes = v
	case string:
		bytes = []byte(v)
	default:
		return fmt.Errorf("JSONMap: 不支持的类型 %T", value)
	}
	return json.Unmarshal(bytes, j)
}

// User 用户角色表
type User struct {
	ID           uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	Username     string    `gorm:"uniqueIndex;size:64;not null" json:"username"`
	PasswordHash string    `gorm:"size:255;not null" json:"-"`
	RealName     string    `gorm:"size:64;not null" json:"real_name"`
	Role         string    `gorm:"size:32;not null" json:"role"`
	Phone        string    `gorm:"size:20" json:"phone"`
	Email        string    `gorm:"size:128" json:"email"`
	Avatar       string    `gorm:"size:512" json:"avatar"`
	OrgName      string    `gorm:"size:128" json:"org_name"`
	IsActive     bool      `gorm:"default:true;not null" json:"is_active"`
	Status       int8      `gorm:"default:1;not null" json:"status"`
	LastLoginAt  time.Time `json:"last_login_at"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

func (User) TableName() string { return "users" }

// Captcha 验证码表
type Captcha struct {
	ID        uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	Key       string    `gorm:"uniqueIndex;size:64;not null" json:"key"`
	Code      string    `gorm:"size:8;not null" json:"code"`
	ExpiredAt time.Time `gorm:"not null" json:"expired_at"`
	Used      int8      `gorm:"default:0;not null" json:"used"`
	CreatedAt time.Time `json:"created_at"`
}

func (Captcha) TableName() string { return "captchas" }

// AgriBatch 农产品批次表
type AgriBatch struct {
	ID          uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	BatchNo     string    `gorm:"uniqueIndex;size:64;not null" json:"batch_no"`
	TraceCode   string    `gorm:"uniqueIndex;size:32;not null" json:"trace_code"`
	ProductName string    `gorm:"size:128;not null" json:"product_name"`
	ProductType string    `gorm:"size:64;not null" json:"product_type"`
	Variety     string    `gorm:"size:64" json:"variety"`
	Quantity    float64   `gorm:"type:decimal(10,2);not null" json:"quantity"`
	Unit        string    `gorm:"size:16;not null" json:"unit"`
	OriginInfo  string    `gorm:"size:255;not null" json:"origin_info"`
	OriginLat   float64   `gorm:"type:decimal(10,7)" json:"origin_lat"`
	OriginLng   float64   `gorm:"type:decimal(10,7)" json:"origin_lng"`
	FarmerID    uint      `gorm:"not null;index" json:"farmer_id"`
	Status      int8      `gorm:"default:0;not null" json:"status"`
	CertNo      string    `gorm:"size:64" json:"cert_no"`
	CoverImage  string    `gorm:"size:512" json:"cover_image"`
	Remark      string    `gorm:"size:512" json:"remark"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

func (AgriBatch) TableName() string { return "agri_batches" }

// DashboardStat 仪表盘统计表
type DashboardStat struct {
	ID            uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	StatDate      string    `gorm:"type:date;uniqueIndex;not null" json:"stat_date"`
	TotalBatches  int       `gorm:"default:0;not null" json:"total_batches"`
	TotalUsers    int       `gorm:"default:0;not null" json:"total_users"`
	TotalVisitors int       `gorm:"default:0;not null" json:"total_visitors"`
	TotalQueries  int       `gorm:"default:0;not null" json:"total_queries"`
	TotalTx       int       `gorm:"default:0;not null" json:"total_tx"`
	TotalAmount   float64   `gorm:"type:decimal(12,2);default:0;not null" json:"total_amount"`
	CreatedAt     time.Time `json:"created_at"`
}

func (DashboardStat) TableName() string { return "dashboard_stats" }

// TraceRecord 溯源节点流转表
type TraceRecord struct {
	ID            uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	BatchID       uint      `gorm:"not null;index" json:"batch_id"`
	NodeType      string    `gorm:"size:32;not null;index" json:"node_type"`
	OperatorID    uint      `gorm:"not null" json:"operator_id"`
	OperationTime time.Time `gorm:"not null" json:"operation_time"`
	Location      string    `gorm:"size:255;not null" json:"location"`
	EnvData       JSONMap   `gorm:"type:json" json:"env_data"`
	TxHash        string    `gorm:"size:128" json:"tx_hash"`
	BlockHeight   int64     `json:"block_height"`
	CreatedAt     time.Time `json:"created_at"`
}

func (TraceRecord) TableName() string { return "trace_records" }

// Certificate 电子合格证表
type Certificate struct {
	ID            uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	BatchID       uint      `gorm:"not null;uniqueIndex" json:"batch_id"`
	CertNo        string    `gorm:"size:64;not null;uniqueIndex" json:"cert_no"`
	ProductName   string    `gorm:"size:128;not null" json:"product_name"`
	ProducerName  string    `gorm:"size:128;not null" json:"producer_name"`
	ProducerAddr  string    `gorm:"size:255" json:"producer_addr"`
	ProducerPhone string    `gorm:"size:20" json:"producer_phone"`
	InspectOrg    string    `gorm:"size:128" json:"inspect_org"`
	Quantity      string    `gorm:"size:64" json:"quantity"`
	IssueDate     string    `gorm:"type:date;not null" json:"issue_date"`
	ValidUntil    string    `gorm:"type:date" json:"valid_until"`
	IssueOrg      string    `gorm:"size:128;not null" json:"issue_org"`
	IssueOrgSeal  string    `gorm:"size:512" json:"issue_org_seal"`
	PesticideOk   bool      `gorm:"default:true;not null" json:"pesticide_ok"`
	HeavyMetalOk  bool      `gorm:"default:true;not null" json:"heavy_metal_ok"`
	MicrobeOk     bool      `gorm:"default:true;not null" json:"microbe_ok"`
	QrCodeUrl     string    `gorm:"size:512" json:"qr_code_url"`
	PdfUrl        string    `gorm:"size:512" json:"pdf_url"`
	TxHash        string    `gorm:"size:128" json:"tx_hash"`
	CreatedAt     time.Time `json:"created_at"`
}

func (Certificate) TableName() string { return "certificates" }
