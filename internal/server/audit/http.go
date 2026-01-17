package audit

import (
	"github.com/gin-gonic/gin"
	"github.com/ydcloud-dy/opshub/internal/service/audit"
)

type HTTPService struct {
	operationLogService *audit.OperationLogService
	loginLogService     *audit.LoginLogService
	dataLogService      *audit.DataLogService
}

func NewHTTPService(
	operationLogService *audit.OperationLogService,
	loginLogService *audit.LoginLogService,
	dataLogService *audit.DataLogService,
) *HTTPService {
	return &HTTPService{
		operationLogService: operationLogService,
		loginLogService:     loginLogService,
		dataLogService:      dataLogService,
	}
}

// RegisterRoutes 注册审计模块路由
func (s *HTTPService) RegisterRoutes(r *gin.RouterGroup) {
	// API v1 审计路由
	audit := r.Group("/audit")
	{
		// 操作日志路由
		operationLogs := audit.Group("/operation-logs")
		{
			operationLogs.GET("", s.operationLogService.ListOperationLogs)
			operationLogs.GET("/:id", s.operationLogService.GetOperationLog)
			operationLogs.DELETE("/:id", s.operationLogService.DeleteOperationLog)
			operationLogs.POST("/batch-delete", s.operationLogService.DeleteOperationLogsBatch)
		}

		// 登录日志路由
		loginLogs := audit.Group("/login-logs")
		{
			loginLogs.GET("", s.loginLogService.ListLoginLogs)
			loginLogs.GET("/:id", s.loginLogService.GetLoginLog)
			loginLogs.DELETE("/:id", s.loginLogService.DeleteLoginLog)
			loginLogs.POST("/batch-delete", s.loginLogService.DeleteLoginLogsBatch)
		}

		// 数据日志路由
		dataLogs := audit.Group("/data-logs")
		{
			dataLogs.GET("", s.dataLogService.ListDataLogs)
			dataLogs.GET("/:id", s.dataLogService.GetDataLog)
			dataLogs.DELETE("/:id", s.dataLogService.DeleteDataLog)
			dataLogs.POST("/batch-delete", s.dataLogService.DeleteDataLogsBatch)
		}
	}
}
