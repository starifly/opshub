package server

import (
	"github.com/gin-gonic/gin"
	"github.com/ydcloud-dy/opshub/plugins/monitor/model"
	"gorm.io/gorm"
)

// RegisterRoutes 注册路由
func RegisterRoutes(router *gin.RouterGroup, db *gorm.DB) {
	handler := NewHandler(db)

	// 域名监控路由组
	domains := router.Group("/domains")
	{
		domains.GET("", handler.ListDomains)           // 获取域名监控列表
		domains.GET("/stats", handler.GetStats)        // 获取统计数据
		domains.GET("/:id", handler.GetDomain)         // 获取域名监控详情
		domains.POST("", handler.CreateDomain)         // 创建域名监控
		domains.PUT("/:id", handler.UpdateDomain)      // 更新域名监控
		domains.DELETE("/:id", handler.DeleteDomain)   // 删除域名监控
		domains.POST("/:id/check", handler.CheckDomain) // 立即检查域名
	}
}

// AutoMigrate 自动迁移表结构
func AutoMigrate(db *gorm.DB) error {
	return db.AutoMigrate(
		&model.DomainMonitor{},
	)
}
