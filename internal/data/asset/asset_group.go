package asset

import (
	"context"
	"github.com/ydcloud-dy/opshub/internal/biz/asset"
	"gorm.io/gorm"
)

type assetGroupRepo struct {
	db *gorm.DB
}

func NewAssetGroupRepo(db *gorm.DB) asset.AssetGroupRepo {
	return &assetGroupRepo{db: db}
}

func (r *assetGroupRepo) Create(ctx context.Context, group *asset.AssetGroup) error {
	return r.db.WithContext(ctx).Create(group).Error
}

func (r *assetGroupRepo) Update(ctx context.Context, group *asset.AssetGroup) error {
	return r.db.WithContext(ctx).Model(group).Omit("created_at").Updates(group).Error
}

func (r *assetGroupRepo) Delete(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// 检查是否有子分组
		var count int64
		if err := tx.Model(&asset.AssetGroup{}).Unscoped().Where("parent_id = ?", id).Count(&count).Error; err != nil {
			return err
		}
		if count > 0 {
			return gorm.ErrRegistered // 存在子分组，不能删除
		}

		// 硬删除分组
		return tx.Unscoped().Delete(&asset.AssetGroup{}, id).Error
	})
}

func (r *assetGroupRepo) GetByID(ctx context.Context, id uint) (*asset.AssetGroup, error) {
	var group asset.AssetGroup
	err := r.db.WithContext(ctx).First(&group, id).Error
	return &group, err
}

func (r *assetGroupRepo) GetTree(ctx context.Context) ([]*asset.AssetGroup, error) {
	var groups []*asset.AssetGroup
	err := r.db.WithContext(ctx).Order("sort ASC").Find(&groups).Error
	if err != nil {
		return nil, err
	}
	return r.buildTree(groups, 0), nil
}

func (r *assetGroupRepo) buildTree(groups []*asset.AssetGroup, parentID uint) []*asset.AssetGroup {
	var tree []*asset.AssetGroup
	for _, group := range groups {
		if group.ParentID == parentID {
			children := r.buildTree(groups, group.ID)
			if len(children) > 0 {
				group.Children = children
			}
			tree = append(tree, group)
		}
	}
	return tree
}

func (r *assetGroupRepo) GetAll(ctx context.Context) ([]*asset.AssetGroup, error) {
	var groups []*asset.AssetGroup
	err := r.db.WithContext(ctx).Order("sort ASC").Find(&groups).Error
	return groups, err
}

func (r *assetGroupRepo) List(ctx context.Context, page, pageSize int, keyword string) ([]*asset.AssetGroup, int64, error) {
	var groups []*asset.AssetGroup
	var total int64

	query := r.db.WithContext(ctx).Model(&asset.AssetGroup{})
	if keyword != "" {
		query = query.Where("name LIKE ? OR code LIKE ?", "%"+keyword+"%", "%"+keyword+"%")
	}

	err := query.Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	err = query.Order("sort ASC").Offset((page - 1) * pageSize).Limit(pageSize).Find(&groups).Error
	return groups, total, err
}
