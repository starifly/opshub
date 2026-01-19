package repository

import (
	"github.com/ydcloud-dy/opshub/plugins/monitor/model"
	"gorm.io/gorm"
)

// DomainMonitorRepository 域名监控数据仓库
type DomainMonitorRepository struct {
	db *gorm.DB
}

// NewDomainMonitorRepository 创建域名监控仓库
func NewDomainMonitorRepository(db *gorm.DB) *DomainMonitorRepository {
	return &DomainMonitorRepository{db: db}
}

// Create 创建域名监控
func (r *DomainMonitorRepository) Create(monitor *model.DomainMonitor) error {
	return r.db.Create(monitor).Error
}

// GetByID 根据ID获取域名监控
func (r *DomainMonitorRepository) GetByID(id uint) (*model.DomainMonitor, error) {
	var monitor model.DomainMonitor
	err := r.db.First(&monitor, id).Error
	if err != nil {
		return nil, err
	}
	return &monitor, nil
}

// GetAll 获取所有域名监控
func (r *DomainMonitorRepository) GetAll() ([]model.DomainMonitor, error) {
	var monitors []model.DomainMonitor
	err := r.db.Order("created_at DESC").Find(&monitors).Error
	return monitors, err
}

// Update 更新域名监控
func (r *DomainMonitorRepository) Update(monitor *model.DomainMonitor) error {
	return r.db.Save(monitor).Error
}

// Delete 删除域名监控
func (r *DomainMonitorRepository) Delete(id uint) error {
	return r.db.Delete(&model.DomainMonitor{}, id).Error
}

// GetByDomain 根据域名获取监控记录
func (r *DomainMonitorRepository) GetByDomain(domain string) (*model.DomainMonitor, error) {
	var monitor model.DomainMonitor
	err := r.db.Where("domain = ?", domain).First(&monitor).Error
	if err != nil {
		return nil, err
	}
	return &monitor, nil
}

// GetStats 获取统计数据
func (r *DomainMonitorRepository) GetStats() (map[string]int64, error) {
	stats := make(map[string]int64)

	// 总数
	var total int64
	if err := r.db.Model(&model.DomainMonitor{}).Count(&total).Error; err != nil {
		return nil, err
	}
	stats["total"] = total

	// 正常数量
	var normal int64
	if err := r.db.Model(&model.DomainMonitor{}).Where("status = ?", "normal").Count(&normal).Error; err != nil {
		return nil, err
	}
	stats["normal"] = normal

	// 异常数量
	var abnormal int64
	if err := r.db.Model(&model.DomainMonitor{}).Where("status = ?", "abnormal").Count(&abnormal).Error; err != nil {
		return nil, err
	}
	stats["abnormal"] = abnormal

	// 暂停数量
	var paused int64
	if err := r.db.Model(&model.DomainMonitor{}).Where("status = ?", "paused").Count(&paused).Error; err != nil {
		return nil, err
	}
	stats["paused"] = paused

	return stats, nil
}
