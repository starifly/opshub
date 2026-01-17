package asset

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/ydcloud-dy/opshub/internal/biz/asset"
	"github.com/ydcloud-dy/opshub/pkg/response"
)

type HostService struct {
	hostUseCase       *asset.HostUseCase
	credentialUseCase *asset.CredentialUseCase
	cloudUseCase      *asset.CloudAccountUseCase
}

func NewHostService(hostUseCase *asset.HostUseCase, credentialUseCase *asset.CredentialUseCase, cloudUseCase *asset.CloudAccountUseCase) *HostService {
	return &HostService{
		hostUseCase:       hostUseCase,
		credentialUseCase: credentialUseCase,
		cloudUseCase:      cloudUseCase,
	}
}

// CreateHost 创建主机
func (s *HostService) CreateHost(c *gin.Context) {
	var req asset.HostRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ErrorCode(c, http.StatusBadRequest, "参数错误: "+err.Error())
		return
	}

	host, err := s.hostUseCase.Create(c.Request.Context(), &req)
	if err != nil {
		response.ErrorCode(c, http.StatusInternalServerError, "创建失败: "+err.Error())
		return
	}

	response.Success(c, host)
}

// UpdateHost 更新主机
func (s *HostService) UpdateHost(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		response.ErrorCode(c, http.StatusBadRequest, "无效的主机ID")
		return
	}

	var req asset.HostRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ErrorCode(c, http.StatusBadRequest, "参数错误: "+err.Error())
		return
	}

	req.ID = uint(id)
	if err := s.hostUseCase.Update(c.Request.Context(), &req); err != nil {
		response.ErrorCode(c, http.StatusInternalServerError, "更新失败: "+err.Error())
		return
	}

	response.SuccessWithMessage(c, "更新成功", nil)
}

// DeleteHost 删除主机
func (s *HostService) DeleteHost(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		response.ErrorCode(c, http.StatusBadRequest, "无效的主机ID")
		return
	}

	if err := s.hostUseCase.Delete(c.Request.Context(), uint(id)); err != nil {
		response.ErrorCode(c, http.StatusInternalServerError, "删除失败: "+err.Error())
		return
	}

	response.SuccessWithMessage(c, "删除成功", nil)
}

// GetHost 获取主机详情
func (s *HostService) GetHost(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		response.ErrorCode(c, http.StatusBadRequest, "无效的主机ID")
		return
	}

	host, err := s.hostUseCase.GetByID(c.Request.Context(), uint(id))
	if err != nil {
		response.ErrorCode(c, http.StatusNotFound, "主机不存在")
		return
	}

	response.Success(c, host)
}

// ListHosts 主机列表
func (s *HostService) ListHosts(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "10"))
	keyword := c.Query("keyword")

	// 支持分组ID筛选
	var groupID *uint
	if groupIDStr := c.Query("groupId"); groupIDStr != "" {
		id, err := strconv.ParseUint(groupIDStr, 10, 32)
		if err == nil {
			gid := uint(id)
			groupID = &gid
		}
	}

	hosts, total, err := s.hostUseCase.List(c.Request.Context(), page, pageSize, keyword, groupID)
	if err != nil {
		response.ErrorCode(c, http.StatusInternalServerError, "查询失败: "+err.Error())
		return
	}

	response.Success(c, gin.H{
		"list":  hosts,
		"total": total,
		"page":  page,
		"pageSize": pageSize,
	})
}

// CreateCredential 创建凭证
func (s *HostService) CreateCredential(c *gin.Context) {
	var req asset.CredentialRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ErrorCode(c, http.StatusBadRequest, "参数错误: "+err.Error())
		return
	}

	credential, err := s.credentialUseCase.Create(c.Request.Context(), &req)
	if err != nil {
		response.ErrorCode(c, http.StatusInternalServerError, "创建失败: "+err.Error())
		return
	}

	response.Success(c, credential)
}

// UpdateCredential 更新凭证
func (s *HostService) UpdateCredential(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		response.ErrorCode(c, http.StatusBadRequest, "无效的凭证ID")
		return
	}

	var req asset.CredentialRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ErrorCode(c, http.StatusBadRequest, "参数错误: "+err.Error())
		return
	}

	req.ID = uint(id)
	if err := s.credentialUseCase.Update(c.Request.Context(), &req); err != nil {
		response.ErrorCode(c, http.StatusInternalServerError, "更新失败: "+err.Error())
		return
	}

	response.SuccessWithMessage(c, "更新成功", nil)
}

// DeleteCredential 删除凭证
func (s *HostService) DeleteCredential(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		response.ErrorCode(c, http.StatusBadRequest, "无效的凭证ID")
		return
	}

	if err := s.credentialUseCase.Delete(c.Request.Context(), uint(id)); err != nil {
		response.ErrorCode(c, http.StatusInternalServerError, "删除失败: "+err.Error())
		return
	}

	response.SuccessWithMessage(c, "删除成功", nil)
}

// GetCredential 获取凭证详情
func (s *HostService) GetCredential(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		response.ErrorCode(c, http.StatusBadRequest, "无效的凭证ID")
		return
	}

	// 获取解密后的凭证详情（用于编辑时回显私钥）
	credential, err := s.credentialUseCase.GetByIDDecrypted(c.Request.Context(), uint(id))
	if err != nil {
		response.ErrorCode(c, http.StatusNotFound, "凭证不存在")
		return
	}

	response.Success(c, credential)
}

// ListCredentials 凭证列表
func (s *HostService) ListCredentials(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "10"))
	keyword := c.Query("keyword")

	credentials, total, err := s.credentialUseCase.List(c.Request.Context(), page, pageSize, keyword)
	if err != nil {
		response.ErrorCode(c, http.StatusInternalServerError, "查询失败: "+err.Error())
		return
	}

	response.Success(c, gin.H{
		"list":     credentials,
		"total":    total,
		"page":     page,
		"pageSize": pageSize,
	})
}

// GetAllCredentials 获取所有凭证（用于下拉选择）
func (s *HostService) GetAllCredentials(c *gin.Context) {
	credentials, err := s.credentialUseCase.GetAll(c.Request.Context())
	if err != nil {
		response.ErrorCode(c, http.StatusInternalServerError, "查询失败: "+err.Error())
		return
	}

	response.Success(c, credentials)
}

// CreateCloudAccount 创建云平台账号
func (s *HostService) CreateCloudAccount(c *gin.Context) {
	var req asset.CloudAccountRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ErrorCode(c, http.StatusBadRequest, "参数错误: "+err.Error())
		return
	}

	account, err := s.cloudUseCase.Create(c.Request.Context(), &req)
	if err != nil {
		response.ErrorCode(c, http.StatusInternalServerError, "创建失败: "+err.Error())
		return
	}

	response.Success(c, account)
}

// UpdateCloudAccount 更新云平台账号
func (s *HostService) UpdateCloudAccount(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		response.ErrorCode(c, http.StatusBadRequest, "无效的账号ID")
		return
	}

	var req asset.CloudAccountRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ErrorCode(c, http.StatusBadRequest, "参数错误: "+err.Error())
		return
	}

	req.ID = uint(id)
	if err := s.cloudUseCase.Update(c.Request.Context(), &req); err != nil {
		response.ErrorCode(c, http.StatusInternalServerError, "更新失败: "+err.Error())
		return
	}

	response.SuccessWithMessage(c, "更新成功", nil)
}

// DeleteCloudAccount 删除云平台账号
func (s *HostService) DeleteCloudAccount(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		response.ErrorCode(c, http.StatusBadRequest, "无效的账号ID")
		return
	}

	if err := s.cloudUseCase.Delete(c.Request.Context(), uint(id)); err != nil {
		response.ErrorCode(c, http.StatusInternalServerError, "删除失败: "+err.Error())
		return
	}

	response.SuccessWithMessage(c, "删除成功", nil)
}

// GetCloudAccount 获取云平台账号详情
func (s *HostService) GetCloudAccount(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		response.ErrorCode(c, http.StatusBadRequest, "无效的账号ID")
		return
	}

	account, err := s.cloudUseCase.GetByID(c.Request.Context(), uint(id))
	if err != nil {
		response.ErrorCode(c, http.StatusNotFound, "账号不存在")
		return
	}

	response.Success(c, account)
}

// ListCloudAccounts 云平台账号列表
func (s *HostService) ListCloudAccounts(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "10"))

	accounts, total, err := s.cloudUseCase.List(c.Request.Context(), page, pageSize)
	if err != nil {
		response.ErrorCode(c, http.StatusInternalServerError, "查询失败: "+err.Error())
		return
	}

	response.Success(c, gin.H{
		"list":     accounts,
		"total":    total,
		"page":     page,
		"pageSize": pageSize,
	})
}

// GetAllCloudAccounts 获取所有启用的云平台账号
func (s *HostService) GetAllCloudAccounts(c *gin.Context) {
	accounts, err := s.cloudUseCase.GetAll(c.Request.Context())
	if err != nil {
		response.ErrorCode(c, http.StatusInternalServerError, "查询失败: "+err.Error())
		return
	}

	response.Success(c, accounts)
}

// ImportFromCloud 从云平台导入主机
func (s *HostService) ImportFromCloud(c *gin.Context) {
	var req asset.CloudImportRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ErrorCode(c, http.StatusBadRequest, "参数错误: "+err.Error())
		return
	}

	if err := s.cloudUseCase.ImportFromCloud(c.Request.Context(), &req, s.hostUseCase); err != nil {
		response.ErrorCode(c, http.StatusInternalServerError, "导入失败: "+err.Error())
		return
	}

	response.SuccessWithMessage(c, "导入成功", nil)
}

// CollectHostInfo 采集主机信息
func (s *HostService) CollectHostInfo(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		response.ErrorCode(c, http.StatusBadRequest, "无效的主机ID")
		return
	}

	if err := s.hostUseCase.CollectHostInfo(c.Request.Context(), uint(id)); err != nil {
		response.ErrorCode(c, http.StatusInternalServerError, "采集失败: "+err.Error())
		return
	}

	response.SuccessWithMessage(c, "采集成功", nil)
}

// TestHostConnection 测试主机连接
func (s *HostService) TestHostConnection(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		response.ErrorCode(c, http.StatusBadRequest, "无效的主机ID")
		return
	}

	if err := s.hostUseCase.TestConnection(c.Request.Context(), uint(id)); err != nil {
		response.ErrorCode(c, http.StatusInternalServerError, "连接失败: "+err.Error())
		return
	}

	response.SuccessWithMessage(c, "连接成功", nil)
}

// BatchCollectHostInfo 批量采集主机信息
func (s *HostService) BatchCollectHostInfo(c *gin.Context) {
	var req struct {
		HostIDs []uint `json:"hostIds" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ErrorCode(c, http.StatusBadRequest, "参数错误: "+err.Error())
		return
	}

	if err := s.hostUseCase.BatchCollectHostInfo(c.Request.Context(), req.HostIDs); err != nil {
		response.ErrorCode(c, http.StatusInternalServerError, "批量采集失败: "+err.Error())
		return
	}

	response.SuccessWithMessage(c, "批量采集完成", nil)
}
