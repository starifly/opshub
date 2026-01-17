package audit

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/ydcloud-dy/opshub/internal/biz/audit"
	"github.com/ydcloud-dy/opshub/pkg/response"
)

type LoginLogService struct {
	useCase *audit.LoginLogUseCase
}

func NewLoginLogService(useCase *audit.LoginLogUseCase) *LoginLogService {
	return &LoginLogService{
		useCase: useCase,
	}
}

// LoginLogListResponse 登录日志列表响应
type LoginLogListResponse struct {
	ID          uint   `json:"id"`
	UserID      uint   `json:"userId"`
	Username    string `json:"username"`
	RealName    string `json:"realName"`
	LoginType   string `json:"loginType"`
	LoginStatus string `json:"loginStatus"`
	LoginTime   string `json:"loginTime"`
	LogoutTime  string `json:"logoutTime"`
	IP          string `json:"ip"`
	Location    string `json:"location"`
	UserAgent   string `json:"userAgent"`
	FailReason  string `json:"failReason"`
}

func toLoginLogListResponse(log *audit.SysLoginLog) LoginLogListResponse {
	resp := LoginLogListResponse{
		ID:          log.ID,
		UserID:      log.UserID,
		Username:    log.Username,
		RealName:    log.RealName,
		LoginType:   log.LoginType,
		LoginStatus: log.LoginStatus,
		LoginTime:   log.LoginTime.Format("2006-01-02 15:04:05"),
		IP:          log.IP,
		Location:    log.Location,
		UserAgent:   log.UserAgent,
		FailReason:  log.FailReason,
	}
	if log.LogoutTime != nil {
		resp.LogoutTime = log.LogoutTime.Format("2006-01-02 15:04:05")
	}
	return resp
}

// ListLoginLogs 登录日志列表
func (s *LoginLogService) ListLoginLogs(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "10"))
	username := c.Query("username")
	loginType := c.Query("loginType")
	loginStatus := c.Query("loginStatus")
	startTime := c.Query("startTime")
	endTime := c.Query("endTime")

	logs, total, err := s.useCase.List(c.Request.Context(), page, pageSize, username, loginType, loginStatus, startTime, endTime)
	if err != nil {
		response.ErrorCode(c, http.StatusInternalServerError, "查询失败: "+err.Error())
		return
	}

	list := make([]LoginLogListResponse, 0, len(logs))
	for _, log := range logs {
		list = append(list, toLoginLogListResponse(log))
	}

	response.Success(c, gin.H{
		"list":     list,
		"page":     page,
		"pageSize": pageSize,
		"total":    total,
	})
}

// GetLoginLog 获取登录日志详情
func (s *LoginLogService) GetLoginLog(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		response.ErrorCode(c, http.StatusBadRequest, "无效的日志ID")
		return
	}

	log, err := s.useCase.GetByID(c.Request.Context(), uint(id))
	if err != nil {
		response.ErrorCode(c, http.StatusNotFound, "日志不存在")
		return
	}

	response.Success(c, toLoginLogListResponse(log))
}

// DeleteLoginLog 删除登录日志
func (s *LoginLogService) DeleteLoginLog(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		response.ErrorCode(c, http.StatusBadRequest, "无效的日志ID")
		return
	}

	if err := s.useCase.Delete(c.Request.Context(), uint(id)); err != nil {
		response.ErrorCode(c, http.StatusInternalServerError, "删除失败: "+err.Error())
		return
	}

	response.SuccessWithMessage(c, "删除成功", nil)
}

// DeleteLoginLogsBatch 批量删除登录日志
func (s *LoginLogService) DeleteLoginLogsBatch(c *gin.Context) {
	var req struct {
		IDs []uint `json:"ids" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ErrorCode(c, http.StatusBadRequest, "参数错误: "+err.Error())
		return
	}

	if err := s.useCase.DeleteBatch(c.Request.Context(), req.IDs); err != nil {
		response.ErrorCode(c, http.StatusInternalServerError, "删除失败: "+err.Error())
		return
	}

	response.SuccessWithMessage(c, "删除成功", nil)
}
