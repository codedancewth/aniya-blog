package repository

import (
	"github.com/cworld1/aniya-blog/backend/internal/database"
	"github.com/cworld1/aniya-blog/backend/internal/models"
	"gorm.io/gorm"
)

// CategoryRepository 分类数据访问层
type CategoryRepository struct {
	db *gorm.DB
}

// NewCategoryRepository 创建分类数据访问层
func NewCategoryRepository() *CategoryRepository {
	return &CategoryRepository{
		db: database.GetDB(),
	}
}

// Create 创建分类
func (r *CategoryRepository) Create(category *models.Category) error {
	return r.db.Create(category).Error
}

// FindByID 根据 ID 查找分类
func (r *CategoryRepository) FindByID(id uint) (*models.Category, error) {
	var category models.Category
	err := r.db.Preload("Parent").Preload("Children").First(&category, id).Error
	return &category, err
}

// FindBySlug 根据 slug 查找分类
func (r *CategoryRepository) FindBySlug(slug string) (*models.Category, error) {
	var category models.Category
	err := r.db.Preload("Parent").Preload("Children").Where("slug = ?", slug).First(&category).Error
	return &category, err
}

// Update 更新分类
func (r *CategoryRepository) Update(category *models.Category) error {
	return r.db.Save(category).Error
}

// Delete 删除分类
func (r *CategoryRepository) Delete(id uint) error {
	return r.db.Delete(&models.Category{}, id).Error
}

// List 获取分类列表
func (r *CategoryRepository) List(page, pageSize int) ([]*models.Category, int64, error) {
	var categories []*models.Category
	var total int64

	offset := (page - 1) * pageSize

	if err := r.db.Model(&models.Category{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	err := r.db.Preload("Parent").Offset(offset).Limit(pageSize).Order("sort_order ASC").Find(&categories).Error
	return categories, total, err
}

// GetAll 获取所有分类（树形结构）
func (r *CategoryRepository) GetAll() ([]*models.Category, error) {
	var categories []*models.Category
	err := r.db.Preload("Parent").Preload("Children").Order("sort_order ASC").Find(&categories).Error
	return categories, err
}

// GetTree 获取分类树
func (r *CategoryRepository) GetTree() ([]*models.Category, error) {
	var rootCategories []*models.Category
	err := r.db.Preload("Children").Where("parent_id IS NULL").Order("sort_order ASC").Find(&rootCategories).Error
	return rootCategories, err
}
