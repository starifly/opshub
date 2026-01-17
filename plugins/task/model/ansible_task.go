package model

import (
	"time"
)

// AnsibleTask Ansible任务
type AnsibleTask struct {
	ID              uint       `json:"id" gorm:"primaryKey"`
	Name            string     `json:"name" gorm:"size:255;not null" binding:"required"`
	PlaybookContent string     `json:"playbookContent,omitempty" gorm:"type:longtext"`
	PlaybookPath    string     `json:"playbookPath,omitempty" gorm:"size:500"`
	Inventory       string     `json:"inventory,omitempty" gorm:"type:text"` // JSON字符串
	ExtraVars       string     `json:"extraVars,omitempty" gorm:"type:json"` // JSON
	Tags            string     `json:"tags,omitempty" gorm:"size:500"` // 逗号分隔
	Fork            int        `json:"fork" gorm:"default:5"`
	Timeout         int        `json:"timeout" gorm:"default:600"` // 秒
	Verbose         string     `json:"verbose" gorm:"size:20;default:v"` // v, vv, vvv
	Status          string     `json:"status" gorm:"size:50;not null;default:pending;index"` // pending, running, success, failed, cancelled
	LastRunTime     *time.Time `json:"lastRunTime,omitempty"`
	LastRunResult   string     `json:"lastRunResult,omitempty" gorm:"type:json"` // JSON
	CreatedBy       uint       `json:"createdBy" gorm:"not null"`
	CreatedAt       time.Time  `json:"createdAt"`
	UpdatedAt       time.Time  `json:"updatedAt"`
	DeletedAt       *time.Time `json:"deletedAt,omitempty" gorm:"index"`
}

func (AnsibleTask) TableName() string {
	return "ansible_tasks"
}
