package audit

import (
	"context"
	"time"

	"github.com/ydcloud-dy/opshub/internal/biz/audit"
	"gorm.io/gorm"
)

type loginLogRepo struct {
	db *gorm.DB
}

func NewLoginLogRepo(db *gorm.DB) audit.LoginLogRepo {
	return &loginLogRepo{db: db}
}

func (r *loginLogRepo) Create(ctx context.Context, log *audit.SysLoginLog) error {
	return r.db.WithContext(ctx).Create(log).Error
}

func (r *loginLogRepo) GetByID(ctx context.Context, id uint) (*audit.SysLoginLog, error) {
	var log audit.SysLoginLog
	err := r.db.WithContext(ctx).First(&log, id).Error
	return &log, err
}

func (r *loginLogRepo) List(ctx context.Context, page, pageSize int, username, loginType, loginStatus, startTime, endTime string) ([]*audit.SysLoginLog, int64, error) {
	var logs []*audit.SysLoginLog
	var total int64

	query := r.db.WithContext(ctx).Model(&audit.SysLoginLog{})

	if username != "" {
		query = query.Where("username LIKE ?", "%"+username+"%")
	}
	if loginType != "" {
		query = query.Where("login_type = ?", loginType)
	}
	if loginStatus != "" {
		query = query.Where("login_status = ?", loginStatus)
	}
	if startTime != "" {
		t, err := time.Parse("2006-01-02", startTime)
		if err == nil {
			query = query.Where("login_time >= ?", t)
		}
	}
	if endTime != "" {
		t, err := time.Parse("2006-01-02", endTime)
		if err == nil {
			t = t.AddDate(0, 0, 1)
			query = query.Where("login_time < ?", t)
		}
	}

	err := query.Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	err = query.Offset((page-1)*pageSize).Limit(pageSize).
		Order("login_time DESC").
		Find(&logs).Error

	return logs, total, err
}

func (r *loginLogRepo) UpdateLogout(ctx context.Context, userID uint, logoutTime *audit.SysLoginLog) error {
	return r.db.WithContext(ctx).
		Where("user_id = ? AND logout_time IS NULL", userID).
		Order("login_time DESC").
		Limit(1).
		Updates(logoutTime).Error
}

func (r *loginLogRepo) Delete(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Delete(&audit.SysLoginLog{}, id).Error
}

func (r *loginLogRepo) DeleteBatch(ctx context.Context, ids []uint) error {
	if len(ids) == 0 {
		return nil
	}
	return r.db.WithContext(ctx).Delete(&audit.SysLoginLog{}, ids).Error
}
