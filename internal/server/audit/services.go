package audit

import (
	"github.com/ydcloud-dy/opshub/internal/biz/audit"
	auditdata "github.com/ydcloud-dy/opshub/internal/data/audit"
	auditservice "github.com/ydcloud-dy/opshub/internal/service/audit"
	"gorm.io/gorm"
)

// NewAuditServices 创建审计模块的所有服务
func NewAuditServices(db *gorm.DB) (
	operationLogService *auditservice.OperationLogService,
	loginLogService *auditservice.LoginLogService,
	dataLogService *auditservice.DataLogService,
) {
	// 初始化Repository
	operationLogRepo := auditdata.NewOperationLogRepo(db)
	loginLogRepo := auditdata.NewLoginLogRepo(db)
	dataLogRepo := auditdata.NewDataLogRepo(db)

	// 初始化UseCase
	operationLogUseCase := audit.NewOperationLogUseCase(operationLogRepo)
	loginLogUseCase := audit.NewLoginLogUseCase(loginLogRepo)
	dataLogUseCase := audit.NewDataLogUseCase(dataLogRepo)

	// 初始化Service
	operationLogService = auditservice.NewOperationLogService(operationLogUseCase)
	loginLogService = auditservice.NewLoginLogService(loginLogUseCase)
	dataLogService = auditservice.NewDataLogService(dataLogUseCase)

	return
}
