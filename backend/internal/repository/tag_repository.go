package repository

import (
	"github.com/cworld1/aniya-blog/backend/internal/database"
	"github.com/cworld1/aniya-blog/backend/internal/models"
	"gorm.io/gorm"
)

// TagRepository 标签数据访问层
type TagRepository struct {
	db *gorm.DB
}

// NewTagRepository 创建标签数据访问层
func NewTagRepository() *TagRepository {
	return &TagRepository{
		db: database.GetDB(),
	}
}

// Create 创建标签
func (r *TagRepository) Create(tag *models.Tag) error {
	return r.db.Create(tag).Error
}

// FindByID 根据 ID 查找标签
func (r *TagRepository) FindByID(id uint) (*models.Tag, error) {
	var tag models.Tag
	err := r.db.First(&tag, id).Error
	return &tag, err
}

// FindBySlug 根据 slug 查找标签
func (r *TagRepository) FindBySlug(slug string) (*models.Tag, error) {
	var tag models.Tag
	err := r.db.Where("slug = ?", slug).First(&tag).Error
	return &tag, err
}

// FindByName 根据名称查找标签
func (r *TagRepository) FindByName(name string) (*models.Tag, error) {
	var tag models.Tag
	err := r.db.Where("name = ?", name).First(&tag).Error
	return &tag, err
}

// Update 更新标签
func (r *TagRepository) Update(tag *models.Tag) error {
	return r.db.Save(tag).Error
}

// Delete 删除标签
func (r *TagRepository) Delete(id uint) error {
	return r.db.Delete(&models.Tag{}, id).Error
}

// List 获取标签列表
func (r *TagRepository) List(page, pageSize int) ([]*models.Tag, int64, error) {
	var tags []*models.Tag
	var total int64

	offset := (page - 1) * pageSize

	if err := r.db.Model(&models.Tag{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	err := r.db.Offset(offset).Limit(pageSize).Order("post_count DESC").Find(&tags).Error
	return tags, total, err
}

// GetAll 获取所有标签
func (r *TagRepository) GetAll() ([]*models.Tag, error) {
	var tags []*models.Tag
	err := r.db.Order("post_count DESC").Find(&tags).Error
	return tags, err
}

// GetOrCreate 获取或创建标签
func (r *TagRepository) GetOrCreate(name, slug string) (*models.Tag, error) {
	tag, err := r.FindBySlug(slug)
	if err == nil {
		return tag, nil
	}

	newTag := &models.Tag{
		Name: name,
		Slug: slug,
	}

	if err := r.Create(newTag); err != nil {
		return nil, err
	}

	return newTag, nil
}
