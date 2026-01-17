package audit

import (
	"context"
	"time"

	"github.com/ydcloud-dy/opshub/internal/biz/audit"
	"gorm.io/gorm"
)

type operationLogRepo struct {
	db *gorm.DB
}

func NewOperationLogRepo(db *gorm.DB) audit.OperationLogRepo {
	return &operationLogRepo{db: db}
}

func (r *operationLogRepo) Create(ctx context.Context, log *audit.SysOperationLog) error {
	return r.db.WithContext(ctx).Create(log).Error
}

func (r *operationLogRepo) GetByID(ctx context.Context, id uint) (*audit.SysOperationLog, error) {
	var log audit.SysOperationLog
	err := r.db.WithContext(ctx).First(&log, id).Error
	return &log, err
}

func (r *operationLogRepo) List(ctx context.Context, page, pageSize int, username, module, action, startTime, endTime string) ([]*audit.SysOperationLog, int64, error) {
	var logs []*audit.SysOperationLog
	var total int64

	query := r.db.WithContext(ctx).Model(&audit.SysOperationLog{})

	if username != "" {
		query = query.Where("username LIKE ?", "%"+username+"%")
	}
	if module != "" {
		query = query.Where("module = ?", module)
	}
	if action != "" {
		query = query.Where("action = ?", action)
	}
	if startTime != "" {
		t, err := time.Parse("2006-01-02", startTime)
		if err == nil {
			query = query.Where("created_at >= ?", t)
		}
	}
	if endTime != "" {
		t, err := time.Parse("2006-01-02", endTime)
		if err == nil {
			// 加一天，包含当天
			t = t.AddDate(0, 0, 1)
			query = query.Where("created_at < ?", t)
		}
	}

	err := query.Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	err = query.Offset((page-1)*pageSize).Limit(pageSize).
		Order("created_at DESC").
		Find(&logs).Error

	return logs, total, err
}

func (r *operationLogRepo) Delete(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Delete(&audit.SysOperationLog{}, id).Error
}

func (r *operationLogRepo) DeleteBatch(ctx context.Context, ids []uint) error {
	if len(ids) == 0 {
		return nil
	}
	return r.db.WithContext(ctx).Delete(&audit.SysOperationLog{}, ids).Error
}
