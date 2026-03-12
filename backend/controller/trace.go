package controller

import (
	"agri-trace/model"
	"agri-trace/pkg/blockchain"
	"agri-trace/pkg/ipfs"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// TraceController 溯源核心控制器
type TraceController struct {
	DB               *gorm.DB
	BlockchainClient blockchain.Client
	IPFSClient       ipfs.Client
}

// ==================== 仪表盘大屏 ====================

// GetDashboardStats 获取大屏统计数据
func (ctrl *TraceController) GetDashboardStats(c *gin.Context) {
	var stats []model.DashboardStat
	// 取最近15天数据
	if err := ctrl.DB.Order("stat_date DESC").Limit(15).Find(&stats).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "msg": "查询统计数据失败"})
		return
	}
	
	// 为了图表展示，需要倒序（时间正序）
	for i, j := 0, len(stats)-1; i < j; i, j = i+1, j-1 {
		stats[i], stats[j] = stats[j], stats[i]
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"data": stats,
	})
}

// ==================== 批次管理 ====================

// CreateBatchRequest 创建批次请求体
type CreateBatchRequest struct {
	ProductName string  `json:"product_name" binding:"required"`
	ProductType string  `json:"product_type" binding:"required"`
	Quantity    float64 `json:"quantity" binding:"required,gt=0"`
	Unit        string  `json:"unit" binding:"required"`
	OriginInfo  string  `json:"origin_info" binding:"required"`
}

// CreateBatch 创建农产品批次
func (ctrl *TraceController) CreateBatch(c *gin.Context) {
	var req CreateBatchRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "msg": err.Error()})
		return
	}

	farmerID, _ := c.Get("user_id")
	// 生成唯一批次号
	batchNo := fmt.Sprintf("BATCH-%s-%s",
		time.Now().Format("20060102"),
		uuid.New().String()[:8])
	
	traceCode := fmt.Sprintf("AKS%s%s", time.Now().Format("200601"), uuid.New().String()[:4])

	batch := model.AgriBatch{
		BatchNo:     batchNo,
		TraceCode:   traceCode,
		ProductName: req.ProductName,
		ProductType: req.ProductType,
		Quantity:    req.Quantity,
		Unit:        req.Unit,
		OriginInfo:  req.OriginInfo,
		FarmerID:    farmerID.(uint),
		Status:      0,
	}

	if err := ctrl.DB.Create(&batch).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "msg": "批次创建失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  "批次创建成功",
		"data": gin.H{"batch_no": batchNo, "trace_code": traceCode, "batch_id": batch.ID},
	})
}

// ListBatches 获取批次列表
func (ctrl *TraceController) ListBatches(c *gin.Context) {
	var batches []model.AgriBatch
	query := ctrl.DB.Model(&model.AgriBatch{})

	// 角色过滤：种植户只能看自己的批次
	role, _ := c.Get("role")
	if role == "farmer" {
		userID, _ := c.Get("user_id")
		query = query.Where("farmer_id = ?", userID)
	}

	if err := query.Order("created_at DESC").Find(&batches).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "msg": "查询失败"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"code": 200, "data": batches})
}

// GetBatchDetail 获取批次详情及完整溯源链路
func (ctrl *TraceController) GetBatchDetail(c *gin.Context) {
	id := c.Param("id")
	var batch model.AgriBatch
	if err := ctrl.DB.First(&batch, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"code": 404, "msg": "批次不存在"})
		return
	}

	// 查所有关联记录
	var records []model.TraceRecord
	ctrl.DB.Where("batch_id = ?", batch.ID).Order("operation_time ASC").Find(&records)

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"data": gin.H{
			"batch":   batch,
			"records": records,
		},
	})
}

// ==================== 溯源记录管理 ====================

// AddTraceRecordRequest 添加溯源记录请求体
type AddTraceRecordRequest struct {
	BatchID       uint                   `json:"batch_id" binding:"required"`
	NodeType      string                 `json:"node_type" binding:"required"`
	OperationTime string                 `json:"operation_time" binding:"required"`
	Location      string                 `json:"location" binding:"required"`
	EnvData       map[string]interface{} `json:"env_data"`
}

// AddTraceRecord 添加溯源节点记录并上链
func (ctrl *TraceController) AddTraceRecord(c *gin.Context) {
	var req AddTraceRecordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "msg": err.Error()})
		return
	}

	// 查询批次
	var batch model.AgriBatch
	if err := ctrl.DB.First(&batch, req.BatchID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"code": 404, "msg": "批次不存在"})
		return
	}

	// 解析操作时间
	opTime, err := time.Parse("2006-01-02 15:04:05", req.OperationTime)
	if err != nil {
		opTime = time.Now()
	}

	operatorID, _ := c.Get("user_id")

	// 构建上链数据
	payload := blockchain.TracePayload{
		BatchNo:       batch.BatchNo,
		NodeType:      req.NodeType,
		OperatorID:    operatorID.(uint),
		OperationTime: opTime,
		Location:      req.Location,
		EnvData:       req.EnvData,
	}

	// 调用区块链（含降级机制）
	txResult, err := ctrl.BlockchainClient.SubmitTrace(payload)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "msg": "上链失败: " + err.Error()})
		return
	}

	// 写入 MySQL
	record := model.TraceRecord{
		BatchID:       req.BatchID,
		NodeType:      req.NodeType,
		OperatorID:    operatorID.(uint),
		OperationTime: opTime,
		Location:      req.Location,
		EnvData:       model.JSONMap(req.EnvData),
		TxHash:        txResult.TxHash,
		BlockHeight:   txResult.BlockHeight,
	}
	if err := ctrl.DB.Create(&record).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "msg": "记录保存失败"})
		return
	}

	// 更新批次状态
	statusMap := map[string]int8{
		"planting": 0, "harvesting": 1, "inspecting": 2,
		"packing": 2, "processing": 2, "transporting": 3, "retailing": 4,
	}
	if newStatus, ok := statusMap[req.NodeType]; ok {
		ctrl.DB.Model(&batch).Update("status", newStatus)
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  "溯源记录添加成功",
		"data": gin.H{
			"record_id":    record.ID,
			"tx_hash":      txResult.TxHash,
			"block_height": txResult.BlockHeight,
			"is_mock":      txResult.IsMock,
		},
	})
}

// ListTraceRecords 获取溯源记录列表
func (ctrl *TraceController) ListTraceRecords(c *gin.Context) {
	batchID := c.Query("batch_id")
	nodeType := c.Query("node_type")
	
	query := ctrl.DB.Model(&model.TraceRecord{})
	if batchID != "" {
		query = query.Where("batch_id = ?", batchID)
	}
	if nodeType != "" {
		query = query.Where("node_type = ?", nodeType)
	}

	var records []model.TraceRecord
	if err := query.Order("operation_time ASC").Find(&records).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "msg": "查询失败"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"code": 200, "data": records})
}

// ==================== 溯源查询与电子合格证 ====================

// QueryByTraceCode 根据溯源码查询完整溯源链
func (ctrl *TraceController) QueryByTraceCode(c *gin.Context) {
	code := c.Param("trace_code")
	var batch model.AgriBatch
	if err := ctrl.DB.Where("trace_code = ?", code).First(&batch).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"code": 404, "msg": "未找到该溯源码对应的产品信息"})
		return
	}

	// 增加一次查询统计（简单实现）
	today := time.Now().Format("2006-01-02")
	ctrl.DB.Exec("UPDATE dashboard_stats SET total_queries = total_queries + 1 WHERE stat_date = ?", today)

	var records []model.TraceRecord
	ctrl.DB.Where("batch_id = ?", batch.ID).Order("operation_time ASC").Find(&records)

	// 查电子合格证
	var cert model.Certificate
	ctrl.DB.Where("batch_id = ?", batch.ID).First(&cert)

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"data": gin.H{
			"batch":       batch,
			"records":     records,
			"certificate": cert,
		},
	})
}

// GetCertificate 获取电子合格证
func (ctrl *TraceController) GetCertificate(c *gin.Context) {
	batchID := c.Param("batch_id")
	var cert model.Certificate
	if err := ctrl.DB.Where("batch_id = ?", batchID).First(&cert).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"code": 404, "msg": "该批次暂无电子合格证"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"code": 200, "data": cert})
}

// GetBlockInfo 查询区块信息
func (ctrl *TraceController) GetBlockInfo(c *gin.Context) {
	// 模拟返回区块信息（实际应调用 blockchainClient）
	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"data": gin.H{
			"channel":        "agrichannel",
			"height":         2156,
			"current_hash":   "0x8b3a4f9e2d1c5b7a6f8e9d0c1b2a3f4e5d6c7b8a9f0e1d2c3b4a5f6e7d8c9b0a",
			"previous_hash":  "0x7a2b3e8d1c0f5a6b4c9d8e7f1a2b3c4d5e6f7a8b9c0d1e2f3a4b5c6d7e8f9a0b",
			"total_tx":       15890,
			"nodes":          []string{"peer0.org1.example.com", "peer0.org2.example.com", "peer0.org3.example.com"},
			"consensus_type": "Raft",
			"crypto_type":    "国密SM2/SM3",
		},
	})
}

// UpdateBatch 更新批次信息
func (ctrl *TraceController) UpdateBatch(c *gin.Context) {
	id := c.Param("id")
	var batch model.AgriBatch
	if err := ctrl.DB.First(&batch, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"code": 404, "msg": "批次不存在"})
		return
	}

	var updates map[string]interface{}
	if err := c.ShouldBindJSON(&updates); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "msg": err.Error()})
		return
	}

	// 不允许修改批次号和溯源码
	delete(updates, "batch_no")
	delete(updates, "trace_code")

	if err := ctrl.DB.Model(&batch).Updates(updates).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "msg": "更新失败"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"code": 200, "msg": "更新成功"})
}

// DeleteBatch 删除批次
func (ctrl *TraceController) DeleteBatch(c *gin.Context) {
	id := c.Param("id")
	if err := ctrl.DB.Delete(&model.AgriBatch{}, id).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "msg": "删除失败"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"code": 200, "msg": "删除成功"})
}
