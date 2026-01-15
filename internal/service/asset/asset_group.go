package asset

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/ydcloud-dy/opshub/internal/biz/asset"
	"github.com/ydcloud-dy/opshub/pkg/response"
)

type AssetGroupService struct {
	groupUseCase *asset.AssetGroupUseCase
}

func NewAssetGroupService(groupUseCase *asset.AssetGroupUseCase) *AssetGroupService {
	return &AssetGroupService{
		groupUseCase: groupUseCase,
	}
}

// CreateGroup 创建分组
func (s *AssetGroupService) CreateGroup(c *gin.Context) {
	var req asset.AssetGroupRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ErrorCode(c, http.StatusBadRequest, "参数错误: "+err.Error())
		return
	}

	group := req.ToModel()
	if err := s.groupUseCase.Create(c.Request.Context(), group); err != nil {
		response.ErrorCode(c, http.StatusInternalServerError, "创建失败: "+err.Error())
		return
	}

	response.Success(c, group)
}

// UpdateGroup 更新分组
func (s *AssetGroupService) UpdateGroup(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		response.ErrorCode(c, http.StatusBadRequest, "无效的分组ID")
		return
	}

	var req asset.AssetGroupRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ErrorCode(c, http.StatusBadRequest, "参数错误: "+err.Error())
		return
	}

	group := req.ToModel()
	group.ID = uint(id)
	if err := s.groupUseCase.Update(c.Request.Context(), group); err != nil {
		response.ErrorCode(c, http.StatusInternalServerError, "更新失败: "+err.Error())
		return
	}

	response.Success(c, group)
}

// DeleteGroup 删除分组
func (s *AssetGroupService) DeleteGroup(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		response.ErrorCode(c, http.StatusBadRequest, "无效的分组ID")
		return
	}

	if err := s.groupUseCase.Delete(c.Request.Context(), uint(id)); err != nil {
		response.ErrorCode(c, http.StatusInternalServerError, "删除失败: "+err.Error())
		return
	}

	response.SuccessWithMessage(c, "删除成功", nil)
}

// GetGroup 获取分组详情
func (s *AssetGroupService) GetGroup(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		response.ErrorCode(c, http.StatusBadRequest, "无效的分组ID")
		return
	}

	group, err := s.groupUseCase.GetByID(c.Request.Context(), uint(id))
	if err != nil {
		response.ErrorCode(c, http.StatusNotFound, "分组不存在")
		return
	}

	response.Success(c, group)
}

// GetGroupTree 获取分组树
func (s *AssetGroupService) GetGroupTree(c *gin.Context) {
	tree, err := s.groupUseCase.GetTree(c.Request.Context())
	if err != nil {
		response.ErrorCode(c, http.StatusInternalServerError, "查询失败: "+err.Error())
		return
	}

	// 转换为VO格式
	var voTree []*asset.AssetGroupInfoVO
	for _, group := range tree {
		voTree = append(voTree, s.groupUseCase.ToInfoVO(group))
	}

	response.Success(c, voTree)
}

// GetParentOptions 获取父级分组选项
func (s *AssetGroupService) GetParentOptions(c *gin.Context) {
	options, err := s.groupUseCase.GetParentOptions(c.Request.Context())
	if err != nil {
		response.ErrorCode(c, http.StatusInternalServerError, "查询失败: "+err.Error())
		return
	}

	response.Success(c, options)
}
