package repository

import (
	"time"

	"github.com/cworld1/aniya-blog/backend/internal/database"
	"github.com/cworld1/aniya-blog/backend/internal/models"
	"gorm.io/gorm"
)

// PostRepository 文章数据访问层
type PostRepository struct {
	db *gorm.DB
}

// NewPostRepository 创建文章数据访问层
func NewPostRepository() *PostRepository {
	return &PostRepository{
		db: database.GetDB(),
	}
}

// Create 创建文章
func (r *PostRepository) Create(post *models.Post) error {
	return r.db.Create(post).Error
}

// FindByID 根据 ID 查找文章
func (r *PostRepository) FindByID(id uint) (*models.Post, error) {
	var post models.Post
	err := r.db.Preload("Author").Preload("Tags").Preload("Category").First(&post, id).Error
	return &post, err
}

// FindBySlug 根据 slug 查找文章
func (r *PostRepository) FindBySlug(slug string) (*models.Post, error) {
	var post models.Post
	err := r.db.Preload("Author").Preload("Tags").Preload("Category").Where("slug = ?", slug).First(&post).Error
	return &post, err
}

// Update 更新文章
func (r *PostRepository) Update(post *models.Post) error {
	return r.db.Save(post).Error
}

// Delete 删除文章
func (r *PostRepository) Delete(id uint) error {
	return r.db.Delete(&models.Post{}, id).Error
}

// List 获取文章列表
func (r *PostRepository) List(page, pageSize int, status *int) ([]*models.Post, int64, error) {
	var posts []*models.Post
	var total int64

	offset := (page - 1) * pageSize

	query := r.db.Model(&models.Post{}).Preload("Author").Preload("Tags").Preload("Category")

	if status != nil {
		query = query.Where("status = ?", *status)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	err := query.Order("is_top DESC, published_at DESC").Offset(offset).Limit(pageSize).Find(&posts).Error
	return posts, total, err
}

// ListByTag 根据标签获取文章列表
func (r *PostRepository) ListByTag(tagSlug string, page, pageSize int) ([]*models.Post, int64, error) {
	var posts []*models.Post
	var total int64

	offset := (page - 1) * pageSize

	query := r.db.Model(&models.Post{}).
		Preload("Author").
		Preload("Tags").
		Joins("JOIN post_tags ON post_tags.post_id = posts.id").
		Joins("JOIN tags ON tags.id = post_tags.tag_id").
		Where("tags.slug = ? AND posts.status = ?", tagSlug, 1)

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	err := query.Order("published_at DESC").Offset(offset).Limit(pageSize).Find(&posts).Error
	return posts, total, err
}

// ListByCategory 根据分类获取文章列表
func (r *PostRepository) ListByCategory(categorySlug string, page, pageSize int) ([]*models.Post, int64, error) {
	var posts []*models.Post
	var total int64

	offset := (page - 1) * pageSize

	query := r.db.Model(&models.Post{}).
		Preload("Author").
		Preload("Tags").
		Joins("JOIN categories ON categories.id = posts.category_id").
		Where("categories.slug = ? AND posts.status = ?", categorySlug, 1)

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	err := query.Order("published_at DESC").Offset(offset).Limit(pageSize).Find(&posts).Error
	return posts, total, err
}

// Search 搜索文章
func (r *PostRepository) Search(keyword string, page, pageSize int) ([]*models.Post, int64, error) {
	var posts []*models.Post
	var total int64

	offset := (page - 1) * pageSize

	query := r.db.Model(&models.Post{}).
		Preload("Author").
		Preload("Tags").
		Where("status = ? AND (title LIKE ? OR description LIKE ? OR content LIKE ?)",
			1, "%"+keyword+"%", "%"+keyword+"%", "%"+keyword+"%")

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	err := query.Order("published_at DESC").Offset(offset).Limit(pageSize).Find(&posts).Error
	return posts, total, err
}

// IncrementViewCount 增加浏览次数
func (r *PostRepository) IncrementViewCount(id uint) error {
	return r.db.Model(&models.Post{}).Where("id = ?", id).Update("view_count", gorm.Expr("view_count + 1")).Error
}

// Publish 发布文章
func (r *PostRepository) Publish(id uint) error {
	now := time.Now()
	return r.db.Model(&models.Post{}).Where("id = ?", id).Updates(map[string]interface{}{
		"status":       1,
		"published_at": now,
	}).Error
}

// GetArchives 获取归档列表（按年份和月份分组）
func (r *PostRepository) GetArchives() ([]*ArchiveResult, error) {
	var results []*ArchiveResult

	err := r.db.Model(&models.Post{}).
		Select("DATE_FORMAT(published_at, '%Y-%m') as month, COUNT(*) as count").
		Where("status = ?", 1).
		Group("DATE_FORMAT(published_at, '%Y-%m')").
		Order("month DESC").
		Scan(&results).Error

	return results, err
}

// ArchiveResult 归档结果
type ArchiveResult struct {
	Month string `json:"month"`
	Count int64  `json:"count"`
}

// GetRelatedPosts 获取相关文章
func (r *PostRepository) GetRelatedPosts(postID uint, tagIDs []uint, limit int) ([]*models.Post, error) {
	var posts []*models.Post

	err := r.db.Model(&models.Post{}).
		Preload("Author").
		Preload("Tags").
		Joins("JOIN post_tags ON post_tags.post_id = posts.id").
		Where("post_tags.tag_id IN ? AND posts.id != ? AND posts.status = ?", tagIDs, postID, 1).
		Group("posts.id").
		Order("COUNT(post_tags.tag_id) DESC").
		Limit(limit).
		Find(&posts).Error

	return posts, err
}
