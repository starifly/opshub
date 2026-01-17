package server

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/ydcloud-dy/opshub/pkg/response"
	"github.com/ydcloud-dy/opshub/plugins/task/model"
	"gorm.io/gorm"
)

type Handler struct {
	db *gorm.DB
}

func NewHandler(db *gorm.DB) *Handler {
	return &Handler{db: db}
}

// ==================== 任务作业 ====================

func (h *Handler) ListJobTasks(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "10"))
	keyword := c.Query("keyword")
	taskType := c.Query("taskType")
	status := c.Query("status")

	var jobTasks []*model.JobTask
	var total int64

	query := h.db.Model(&model.JobTask{}).Where("deleted_at IS NULL")

	if keyword != "" {
		query = query.Where("name LIKE ?", "%"+keyword+"%")
	}
	if taskType != "" {
		query = query.Where("task_type = ?", taskType)
	}
	if status != "" {
		query = query.Where("status = ?", status)
	}

	query.Count(&total)
	offset := (page - 1) * pageSize
	query.Order("created_at DESC").Limit(pageSize).Offset(offset).Find(&jobTasks)

	response.Success(c, gin.H{
		"list":     jobTasks,
		"total":    total,
		"page":     page,
		"pageSize": pageSize,
	})
}

func (h *Handler) GetJobTask(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 32)
	var jobTask model.JobTask
	if err := h.db.Where("id = ? AND deleted_at IS NULL", id).First(&jobTask).Error; err != nil {
		response.ErrorCode(c, http.StatusNotFound, "任务不存在")
		return
	}
	response.Success(c, jobTask)
}

func (h *Handler) CreateJobTask(c *gin.Context) {
	var jobTask model.JobTask
	if err := c.ShouldBindJSON(&jobTask); err != nil {
		response.ErrorCode(c, http.StatusBadRequest, "参数错误")
		return
	}
	jobTask.Status = "pending"
	jobTask.CreatedBy = 1 // TODO: 从JWT获取
	if err := h.db.Create(&jobTask).Error; err != nil {
		response.ErrorCode(c, http.StatusInternalServerError, "创建失败")
		return
	}
	response.Success(c, jobTask)
}

func (h *Handler) UpdateJobTask(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 32)
	var jobTask model.JobTask
	if err := h.db.Where("id = ? AND deleted_at IS NULL", id).First(&jobTask).Error; err != nil {
		response.ErrorCode(c, http.StatusNotFound, "任务不存在")
		return
	}
	if err := c.ShouldBindJSON(&jobTask); err != nil {
		response.ErrorCode(c, http.StatusBadRequest, "参数错误")
		return
	}
	h.db.Save(&jobTask)
	response.Success(c, jobTask)
}

func (h *Handler) DeleteJobTask(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 32)
	h.db.Delete(&model.JobTask{}, id)
	response.SuccessWithMessage(c, "删除成功", nil)
}

// ==================== 任务模板 ====================

func (h *Handler) ListJobTemplates(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "10"))
	keyword := c.Query("keyword")
	category := c.Query("category")

	var templates []*model.JobTemplate
	var total int64

	query := h.db.Model(&model.JobTemplate{}).Where("deleted_at IS NULL")

	if keyword != "" {
		query = query.Where("name LIKE ? OR code LIKE ?", "%"+keyword+"%", "%"+keyword+"%")
	}
	if category != "" {
		query = query.Where("category = ?", category)
	}

	query.Count(&total)
	offset := (page - 1) * pageSize
	query.Order("sort ASC, created_at DESC").Limit(pageSize).Offset(offset).Find(&templates)

	response.Success(c, gin.H{
		"list":     templates,
		"total":    total,
		"page":     page,
		"pageSize": pageSize,
	})
}

func (h *Handler) GetAllJobTemplates(c *gin.Context) {
	category := c.Query("category")
	var templates []*model.JobTemplate
	query := h.db.Where("deleted_at IS NULL AND status = ?", 1)
	if category != "" {
		query = query.Where("category = ?", category)
	}
	query.Order("sort ASC").Find(&templates)
	response.Success(c, templates)
}

func (h *Handler) GetJobTemplate(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 32)
	var template model.JobTemplate
	if err := h.db.Where("id = ? AND deleted_at IS NULL", id).First(&template).Error; err != nil {
		response.ErrorCode(c, http.StatusNotFound, "模板不存在")
		return
	}
	response.Success(c, template)
}

func (h *Handler) CreateJobTemplate(c *gin.Context) {
	var template model.JobTemplate
	if err := c.ShouldBindJSON(&template); err != nil {
		response.ErrorCode(c, http.StatusBadRequest, "参数错误")
		return
	}
	template.Status = 1
	template.CreatedBy = 1
	if err := h.db.Create(&template).Error; err != nil {
		response.ErrorCode(c, http.StatusInternalServerError, "创建失败")
		return
	}
	response.Success(c, template)
}

func (h *Handler) UpdateJobTemplate(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 32)
	var template model.JobTemplate
	if err := h.db.Where("id = ? AND deleted_at IS NULL", id).First(&template).Error; err != nil {
		response.ErrorCode(c, http.StatusNotFound, "模板不存在")
		return
	}
	if err := c.ShouldBindJSON(&template); err != nil {
		response.ErrorCode(c, http.StatusBadRequest, "参数错误")
		return
	}
	h.db.Save(&template)
	response.Success(c, template)
}

func (h *Handler) DeleteJobTemplate(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 32)
	h.db.Delete(&model.JobTemplate{}, id)
	response.SuccessWithMessage(c, "删除成功", nil)
}

// ==================== Ansible任务 ====================

func (h *Handler) ListAnsibleTasks(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "10"))
	keyword := c.Query("keyword")
	status := c.Query("status")

	var tasks []*model.AnsibleTask
	var total int64

	query := h.db.Model(&model.AnsibleTask{}).Where("deleted_at IS NULL")

	if keyword != "" {
		query = query.Where("name LIKE ?", "%"+keyword+"%")
	}
	if status != "" {
		query = query.Where("status = ?", status)
	}

	query.Count(&total)
	offset := (page - 1) * pageSize
	query.Order("created_at DESC").Limit(pageSize).Offset(offset).Find(&tasks)

	response.Success(c, gin.H{
		"list":     tasks,
		"total":    total,
		"page":     page,
		"pageSize": pageSize,
	})
}

func (h *Handler) GetAnsibleTask(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 32)
	var task model.AnsibleTask
	if err := h.db.Where("id = ? AND deleted_at IS NULL", id).First(&task).Error; err != nil {
		response.ErrorCode(c, http.StatusNotFound, "任务不存在")
		return
	}
	response.Success(c, task)
}

func (h *Handler) CreateAnsibleTask(c *gin.Context) {
	var ansibleTask model.AnsibleTask
	if err := c.ShouldBindJSON(&ansibleTask); err != nil {
		response.ErrorCode(c, http.StatusBadRequest, "参数错误")
		return
	}
	ansibleTask.Status = "pending"
	if ansibleTask.Fork == 0 {
		ansibleTask.Fork = 5
	}
	if ansibleTask.Timeout == 0 {
		ansibleTask.Timeout = 600
	}
	if ansibleTask.Verbose == "" {
		ansibleTask.Verbose = "v"
	}
	ansibleTask.CreatedBy = 1
	if err := h.db.Create(&ansibleTask).Error; err != nil {
		response.ErrorCode(c, http.StatusInternalServerError, "创建失败")
		return
	}
	response.Success(c, ansibleTask)
}

func (h *Handler) UpdateAnsibleTask(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 32)
	var ansibleTask model.AnsibleTask
	if err := h.db.Where("id = ? AND deleted_at IS NULL", id).First(&ansibleTask).Error; err != nil {
		response.ErrorCode(c, http.StatusNotFound, "任务不存在")
		return
	}
	if err := c.ShouldBindJSON(&ansibleTask); err != nil {
		response.ErrorCode(c, http.StatusBadRequest, "参数错误")
		return
	}
	h.db.Save(&ansibleTask)
	response.Success(c, ansibleTask)
}

func (h *Handler) DeleteAnsibleTask(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 32)
	h.db.Delete(&model.AnsibleTask{}, id)
	response.SuccessWithMessage(c, "删除成功", nil)
}
