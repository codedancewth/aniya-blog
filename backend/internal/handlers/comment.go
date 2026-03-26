package handlers

import (
	"strconv"

	"github.com/cworld1/aniya-blog/backend/internal/models"
	"github.com/cworld1/aniya-blog/backend/internal/repository"
	"github.com/cworld1/aniya-blog/backend/pkg/response"
	"github.com/cworld1/aniya-blog/backend/pkg/validator"
	"github.com/gin-gonic/gin"
)

// CommentHandler 评论处理器
type CommentHandler struct {
	commentRepo *repository.CommentRepository
	postRepo    *repository.PostRepository
}

// NewCommentHandler 创建评论处理器
func NewCommentHandler(commentRepo *repository.CommentRepository, postRepo *repository.PostRepository) *CommentHandler {
	return &CommentHandler{
		commentRepo: commentRepo,
		postRepo:    postRepo,
	}
}

// CreateCommentRequest 创建评论请求
type CreateCommentRequest struct {
	Content    string `json:"content" binding:"required"`
	PostSlug   string `json:"post_slug" binding:"required"`
	PostID     *uint  `json:"post_id"`  // 兼容旧版本
	ParentID   *uint  `json:"parent_id"`
	AuthorName string `json:"author_name"`
	AuthorEmail string `json:"author_email"`
	AuthorURL  string `json:"author_url"`
}

// CreateComment 创建评论
// @Summary 创建评论
// @Tags 评论
// @Accept json
// @Produce json
// @Param request body CreateCommentRequest true "评论信息"
// @Success 200 {object} response.Response
// @Router /api/v1/comments [post]
func (h *CommentHandler) CreateComment(c *gin.Context) {
	var req CreateCommentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, response.ERROR, validator.GetErrorMsg(err))
		return
	}

	// 检查文章是否存在（优先使用 slug，兼容 post_id）
	var postID uint
	if req.PostSlug != "" {
		post, err := h.postRepo.FindBySlug(req.PostSlug)
		if err != nil {
			response.Error(c, response.POST_NOT_FOUND, "post not found by slug: "+req.PostSlug)
			return
		}
		postID = post.ID
	} else if req.PostID != nil {
		postID = *req.PostID
		_, err := h.postRepo.FindByID(postID)
		if err != nil {
			response.Error(c, response.POST_NOT_FOUND, "post not found")
			return
		}
	} else {
		response.Error(c, response.ERROR, "post_slug or post_id is required")
		return
	}

	// 如果是回复评论，检查父评论是否存在
	if req.ParentID != nil {
		_, err := h.commentRepo.FindByID(*req.ParentID)
		if err != nil {
			response.Error(c, response.COMMENT_NOT_FOUND, "parent comment not found")
			return
		}
	}

	// 获取用户信息（如果已登录）
	var userID *uint
	var authorName string
	var authorEmail string
	var isAdmin bool

	if id, exists := c.Get("user_id"); exists {
		userIDVal := id.(uint)
		userID = &userIDVal
		authorName = c.GetString("username")
		isAdmin = true
	} else {
		authorName = req.AuthorName
		authorEmail = req.AuthorEmail
		isAdmin = false
	}

	comment := &models.Comment{
		Content:     req.Content,
		PostID:      postID,
		ParentID:    req.ParentID,
		UserID:      userID,
		AuthorName:  authorName,
		AuthorEmail: authorEmail,
		AuthorURL:   req.AuthorURL,
		AuthorIP:    c.ClientIP(),
		Agent:       c.GetHeader("User-Agent"),
		Status:      1, // 默认审核通过
		IsAdmin:     isAdmin,
	}

	if err := h.commentRepo.Create(comment); err != nil {
		response.Error(c, response.ERROR, "failed to create comment")
		return
	}

	response.Success(c, comment)
}

// GetComment 获取评论详情
// @Summary 获取评论详情
// @Tags 评论
// @Produce json
// @Param id path uint true "评论 ID"
// @Success 200 {object} response.Response
// @Router /api/v1/comments/:id [get]
func (h *CommentHandler) GetComment(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.Error(c, response.ERROR, "invalid comment id")
		return
	}

	comment, err := h.commentRepo.FindByID(uint(id))
	if err != nil {
		response.Error(c, response.COMMENT_NOT_FOUND, "comment not found")
		return
	}

	response.Success(c, comment)
}

// ListCommentsByPost 根据文章 ID 或 slug 获取评论列表
// @Summary 根据文章 ID 或 slug 获取评论列表
// @Tags 评论
// @Produce json
// @Param post_id path string true "文章 ID 或 slug"
// @Param page query int false "页码" default(1)
// @Param pageSize query int false "每页数量" default(10)
// @Success 200 {object} response.Response
// @Router /api/v1/posts/:post_id/comments [get]
func (h *CommentHandler) ListCommentsByPost(c *gin.Context) {
	postParam := c.Param("post_id")
	
	// 尝试解析为数字 ID，如果失败则作为 slug 处理
	var postID uint
	post, err := h.postRepo.FindBySlug(postParam)
	if err != nil {
		// 尝试作为数字 ID 解析
		id, parseErr := strconv.ParseUint(postParam, 10, 32)
		if parseErr != nil {
			response.Error(c, response.ERROR, "invalid post id or slug")
			return
		}
		postID = uint(id)
	} else {
		postID = post.ID
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "10"))

	comments, total, err := h.commentRepo.ListByPostID(uint(postID), page, pageSize)
	if err != nil {
		response.Error(c, response.ERROR, "failed to get comments")
		return
	}

	response.PageSuccess(c, comments, total, page, pageSize)
}

// UpdateComment 更新评论
// @Summary 更新评论
// @Tags 评论
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path uint true "评论 ID"
// @Param request body UpdateCommentRequest true "评论信息"
// @Success 200 {object} response.Response
// @Router /api/v1/comments/:id [put]
func (h *CommentHandler) UpdateComment(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.Error(c, response.ERROR, "invalid comment id")
		return
	}

	comment, err := h.commentRepo.FindByID(uint(id))
	if err != nil {
		response.Error(c, response.COMMENT_NOT_FOUND, "comment not found")
		return
	}

	var req UpdateCommentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, response.ERROR, validator.GetErrorMsg(err))
		return
	}

	if req.Content != "" {
		comment.Content = req.Content
	}
	if req.Status != 0 {
		comment.Status = req.Status
	}

	if err := h.commentRepo.Update(comment); err != nil {
		response.Error(c, response.ERROR, "failed to update comment")
		return
	}

	response.Success(c, comment)
}

// DeleteComment 删除评论
// @Summary 删除评论
// @Tags 评论
// @Produce json
// @Security BearerAuth
// @Param id path uint true "评论 ID"
// @Success 200 {object} response.Response
// @Router /api/v1/comments/:id [delete]
func (h *CommentHandler) DeleteComment(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.Error(c, response.ERROR, "invalid comment id")
		return
	}

	if err := h.commentRepo.Delete(uint(id)); err != nil {
		response.Error(c, response.ERROR, "failed to delete comment")
		return
	}

	response.Success(c, nil)
}

// LikeComment 点赞评论
// @Summary 点赞评论
// @Tags 评论
// @Produce json
// @Param id path uint true "评论 ID"
// @Success 200 {object} response.Response
// @Router /api/v1/comments/:id/like [post]
func (h *CommentHandler) LikeComment(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.Error(c, response.ERROR, "invalid comment id")
		return
	}

	if err := h.commentRepo.IncrementLikeCount(uint(id)); err != nil {
		response.Error(c, response.ERROR, "failed to like comment")
		return
	}

	response.Success(c, nil)
}

// UpdateCommentRequest 更新评论请求
type UpdateCommentRequest struct {
	Content string `json:"content"`
	Status  int    `json:"status"`
}
