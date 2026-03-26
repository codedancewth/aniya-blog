package repository

import (
	"github.com/cworld1/aniya-blog/backend/internal/database"
	"github.com/cworld1/aniya-blog/backend/internal/models"
	"gorm.io/gorm"
)

// CommentRepository 评论数据访问层
type CommentRepository struct {
	db *gorm.DB
}

// NewCommentRepository 创建评论数据访问层
func NewCommentRepository() *CommentRepository {
	return &CommentRepository{
		db: database.GetDB(),
	}
}

// Create 创建评论
func (r *CommentRepository) Create(comment *models.Comment) error {
	return r.db.Create(comment).Error
}

// FindByID 根据 ID 查找评论
func (r *CommentRepository) FindByID(id uint) (*models.Comment, error) {
	var comment models.Comment
	err := r.db.Preload("User").Preload("Replies").First(&comment, id).Error
	return &comment, err
}

// Update 更新评论
func (r *CommentRepository) Update(comment *models.Comment) error {
	return r.db.Save(comment).Error
}

// Delete 删除评论
func (r *CommentRepository) Delete(id uint) error {
	return r.db.Delete(&models.Comment{}, id).Error
}

// ListByPostID 根据文章 ID 获取评论列表
func (r *CommentRepository) ListByPostID(postID uint, page, pageSize int) ([]*models.Comment, int64, error) {
	var comments []*models.Comment
	var total int64

	offset := (page - 1) * pageSize

	// 先查询总数
	countQuery := r.db.Model(&models.Comment{}).
		Where("post_id = ? AND status = ? AND parent_id IS NULL", postID, 1)
	
	if err := countQuery.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 查询评论（不 Preload User，避免 NULL 用户 ID 导致的问题）
	query := r.db.Model(&models.Comment{}).
		Where("post_id = ? AND status = ? AND parent_id IS NULL", postID, 1).
		Order("created_at DESC").
		Offset(offset).
		Limit(pageSize)

	err := query.Find(&comments).Error
	return comments, total, err
}

// List 获取评论列表
func (r *CommentRepository) List(page, pageSize int, status *int) ([]*models.Comment, int64, error) {
	var comments []*models.Comment
	var total int64

	offset := (page - 1) * pageSize

	query := r.db.Model(&models.Comment{}).Preload("User").Preload("Post")

	if status != nil {
		query = query.Where("status = ?", *status)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	err := query.Order("created_at DESC").Offset(offset).Limit(pageSize).Find(&comments).Error
	return comments, total, err
}

// CountByPostID 根据文章 ID 获取评论数量
func (r *CommentRepository) CountByPostID(postID uint) (int64, error) {
	var count int64
	err := r.db.Model(&models.Comment{}).Where("post_id = ? AND status = ?", postID, 1).Count(&count).Error
	return count, err
}

// IncrementLikeCount 增加点赞数
func (r *CommentRepository) IncrementLikeCount(id uint) error {
	return r.db.Model(&models.Comment{}).Where("id = ?", id).Update("like_count", gorm.Expr("like_count + 1")).Error
}
