package router

import (
	"agri-trace/controller"
	"agri-trace/middleware"

	"github.com/gin-gonic/gin"
)

func Setup(
	authCtrl *controller.AuthController,
	traceCtrl *controller.TraceController,
	adminCtrl *controller.AdminController,
) *gin.Engine {
	r := gin.Default()

	r.Use(middleware.CORS())

	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok", "service": "agri-trace-backend"})
	})

	v1 := r.Group("/api/v1")
	{
		// ==================== 无需鉴权 ====================
		auth := v1.Group("/auth")
		{
			auth.GET("/captcha", authCtrl.GetCaptcha)
			auth.POST("/login", authCtrl.Login)
			auth.POST("/register", authCtrl.Register)
		}

		// 溯源公开查询（消费者扫码入口，一物一码）
		v1.GET("/trace/public/:trace_code", traceCtrl.QueryByTraceCode)
		// 区块链节点信息（大屏展示）
		v1.GET("/block/info", traceCtrl.GetBlockInfo)
		// 仪表盘统计（大屏展示）
		v1.GET("/dashboard/stats", traceCtrl.GetDashboardStats)

		// ==================== 需要JWT鉴权 ====================
		authorized := v1.Group("/")
		authorized.Use(middleware.JWTAuth())
		{
			// 个人信息
			authorized.GET("/auth/profile", authCtrl.GetProfile)

			// 批次管理
			authorized.POST("/batches", middleware.RoleRequired("farmer", "admin"), traceCtrl.CreateBatch)
			authorized.GET("/batches", traceCtrl.ListBatches)
			authorized.GET("/batches/:id", traceCtrl.GetBatchDetail)
			authorized.PUT("/batches/:id", middleware.RoleRequired("farmer", "admin"), traceCtrl.UpdateBatch)
			authorized.DELETE("/batches/:id", middleware.RoleRequired("admin"), traceCtrl.DeleteBatch)

			// 溯源记录管理
			authorized.POST("/trace/records", traceCtrl.AddTraceRecord)
			authorized.GET("/trace/records", traceCtrl.ListTraceRecords)

			// 电子合格证
			authorized.GET("/certificates/:batch_id", traceCtrl.GetCertificate)
			authorized.POST("/certificates", middleware.RoleRequired("inspector", "admin"), adminCtrl.CreateCertificate)

			// ==================== 管理员专属 ====================
			admin := authorized.Group("/admin")
			admin.Use(middleware.RoleRequired("admin"))
			{
				admin.GET("/users", adminCtrl.ListUsers)
				admin.POST("/users", adminCtrl.CreateUser)
				admin.PUT("/users/:id", adminCtrl.UpdateUser)
				admin.DELETE("/users/:id", adminCtrl.DeleteUser)
			}
		}
	}

	return r
}
