package repository

import (
	"github.com/cworld1/aniya-blog/backend/internal/database"
	"github.com/cworld1/aniya-blog/backend/internal/models"
	"gorm.io/gorm"
)

// PageViewRepository 页面浏览数据访问层
type PageViewRepository struct {
	db *gorm.DB
}

// NewPageViewRepository 创建页面浏览数据访问层
func NewPageViewRepository() *PageViewRepository {
	return &PageViewRepository{
		db: database.GetDB(),
	}
}

// Create 创建页面浏览记录
func (r *PageViewRepository) Create(pageView *models.PageView) error {
	return r.db.Create(pageView).Error
}

// CountByPath 根据路径获取浏览次数
func (r *PageViewRepository) CountByPath(path string) (int64, error) {
	var count int64
	err := r.db.Model(&models.PageView{}).Where("path = ?", path).Count(&count).Error
	return count, err
}

// GetTotalViews 获取总浏览次数
func (r *PageViewRepository) GetTotalViews() (int64, error) {
	var count int64
	err := r.db.Model(&models.PageView{}).Count(&count).Error
	return count, err
}

// GetTotalPosts 获取总文章数
func (r *PageViewRepository) GetTotalPosts() (int64, error) {
	var count int64
	err := r.db.Model(&models.Post{}).Where("status = ?", 1).Count(&count).Error
	return count, err
}

// GetViewsByPaths 批量获取路径浏览次数
func (r *PageViewRepository) GetViewsByPaths(paths []string) (map[string]int64, error) {
	result := make(map[string]int64)

	for _, path := range paths {
		count, err := r.CountByPath(path)
		if err != nil {
			continue
		}
		result[path] = count
	}

	return result, nil
}
