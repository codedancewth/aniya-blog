package handlers

import (
	"regexp"
	"strconv"
	"strings"

	"github.com/cworld1/aniya-blog/backend/internal/models"
	"github.com/cworld1/aniya-blog/backend/internal/repository"
	"github.com/cworld1/aniya-blog/backend/pkg/response"
	"github.com/cworld1/aniya-blog/backend/pkg/validator"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// PostHandler 文章处理器
type PostHandler struct {
	postRepo     *repository.PostRepository
	tagRepo      *repository.TagRepository
	categoryRepo *repository.CategoryRepository
}

// NewPostHandler 创建文章处理器
func NewPostHandler(postRepo *repository.PostRepository, tagRepo *repository.TagRepository, categoryRepo *repository.CategoryRepository) *PostHandler {
	return &PostHandler{
		postRepo:     postRepo,
		tagRepo:      tagRepo,
		categoryRepo: categoryRepo,
	}
}

// CreatePostRequest 创建文章请求
type CreatePostRequest struct {
	Title       string   `json:"title" binding:"required"`
	Slug        string   `json:"slug"`
	Description string   `json:"description"`
	Content     string   `json:"content"`
	ContentHTML string   `json:"content_html"`
	CoverImage  string   `json:"cover_image"`
	CategoryID  *uint    `json:"category_id"`
	TagNames    []string `json:"tag_names"`
	Language    string   `json:"language"`
	IsTop       bool     `json:"is_top"`
	CustomData  string   `json:"custom_data"`
	Status      int      `json:"status"` // 1: published, 0: draft
}

// UpdatePostRequest 更新文章请求
type UpdatePostRequest struct {
	Title       string   `json:"title"`
	Slug        string   `json:"slug"`
	Description string   `json:"description"`
	Content     string   `json:"content"`
	ContentHTML string   `json:"content_html"`
	CoverImage  string   `json:"cover_image"`
	CategoryID  *uint    `json:"category_id"`
	TagNames    []string `json:"tag_names"`
	Language    string   `json:"language"`
	IsTop       bool     `json:"is_top"`
	CustomData  string   `json:"custom_data"`
	Status      int      `json:"status"`
}

// CreatePost 创建文章
// @Summary 创建文章
// @Tags 文章
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body CreatePostRequest true "文章信息"
// @Success 200 {object} response.Response
// @Router /api/v1/posts [post]
func (h *PostHandler) CreatePost(c *gin.Context) {
	var req CreatePostRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, response.ERROR, validator.GetErrorMsg(err))
		return
	}

	// 生成 slug
	slug := req.Slug
	if slug == "" {
		slug = uuid.New().String()[:8]
	}

	// 获取或创建标签
	var tags []models.Tag
	for _, tagName := range req.TagNames {
		tagSlug := generateSlug(tagName)
		tag, err := h.tagRepo.GetOrCreate(tagName, tagSlug)
		if err != nil {
			response.Error(c, response.ERROR, "failed to create tag")
			return
		}
		tags = append(tags, *tag)
	}

	// 创建文章
	userID, _ := c.Get("user_id")
	post := &models.Post{
		Title:       req.Title,
		Slug:        slug,
		Description: req.Description,
		Content:     req.Content,
		ContentHTML: req.ContentHTML,
		CoverImage:  req.CoverImage,
		AuthorID:    userID.(uint),
		Language:    req.Language,
		IsTop:       req.IsTop,
		CustomData:  req.CustomData,
		Status:      req.Status,
		Tags:        tags,
	}

	if req.CategoryID != nil {
		post.CategoryID = req.CategoryID
	}

	if err := h.postRepo.Create(post); err != nil {
		response.Error(c, response.ERROR, "failed to create post")
		return
	}

	response.Success(c, post)
}

// GetPost 获取文章详情
// @Summary 获取文章详情
// @Tags 文章
// @Produce json
// @Param id path uint true "文章 ID"
// @Success 200 {object} response.Response
// @Router /api/v1/posts/:id [get]
func (h *PostHandler) GetPost(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.Error(c, response.ERROR, "invalid post id")
		return
	}

	post, err := h.postRepo.FindByID(uint(id))
	if err != nil {
		response.Error(c, response.POST_NOT_FOUND, "post not found")
		return
	}

	response.Success(c, post)
}

// GetPostBySlug 根据 slug 获取文章详情
// @Summary 根据 slug 获取文章详情
// @Tags 文章
// @Produce json
// @Param slug path string true "文章 slug"
// @Success 200 {object} response.Response
// @Router /api/v1/posts/slug/:slug [get]
func (h *PostHandler) GetPostBySlug(c *gin.Context) {
	slug := c.Param("slug")

	post, err := h.postRepo.FindBySlug(slug)
	if err != nil {
		response.Error(c, response.POST_NOT_FOUND, "post not found")
		return
	}

	// 增加浏览次数
	_ = h.postRepo.IncrementViewCount(post.ID)

	response.Success(c, post)
}

// UpdatePost 更新文章
// @Summary 更新文章
// @Tags 文章
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path uint true "文章 ID"
// @Param request body UpdatePostRequest true "文章信息"
// @Success 200 {object} response.Response
// @Router /api/v1/posts/:id [put]
func (h *PostHandler) UpdatePost(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.Error(c, response.ERROR, "invalid post id")
		return
	}

	post, err := h.postRepo.FindByID(uint(id))
	if err != nil {
		response.Error(c, response.POST_NOT_FOUND, "post not found")
		return
	}

	var req UpdatePostRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, response.ERROR, validator.GetErrorMsg(err))
		return
	}

	// 更新字段
	if req.Title != "" {
		post.Title = req.Title
	}
	if req.Slug != "" {
		post.Slug = req.Slug
	}
	if req.Description != "" {
		post.Description = req.Description
	}
	if req.Content != "" {
		post.Content = req.Content
	}
	if req.ContentHTML != "" {
		post.ContentHTML = req.ContentHTML
	}
	if req.CoverImage != "" {
		post.CoverImage = req.CoverImage
	}
	if req.CategoryID != nil {
		post.CategoryID = req.CategoryID
	}
	if req.Language != "" {
		post.Language = req.Language
	}
	post.IsTop = req.IsTop
	if req.CustomData != "" {
		post.CustomData = req.CustomData
	}
	if req.Status != 0 {
		post.Status = req.Status
	}

	// 更新标签
	if req.TagNames != nil {
		var tags []models.Tag
		for _, tagName := range req.TagNames {
			tagSlug := generateSlug(tagName)
			tag, err := h.tagRepo.GetOrCreate(tagName, tagSlug)
			if err != nil {
				response.Error(c, response.ERROR, "failed to create tag")
				return
			}
			tags = append(tags, *tag)
		}
		post.Tags = tags
	}

	if err := h.postRepo.Update(post); err != nil {
		response.Error(c, response.ERROR, "failed to update post")
		return
	}

	response.Success(c, post)
}

// DeletePost 删除文章
// @Summary 删除文章
// @Tags 文章
// @Produce json
// @Security BearerAuth
// @Param id path uint true "文章 ID"
// @Success 200 {object} response.Response
// @Router /api/v1/posts/:id [delete]
func (h *PostHandler) DeletePost(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.Error(c, response.ERROR, "invalid post id")
		return
	}

	if err := h.postRepo.Delete(uint(id)); err != nil {
		response.Error(c, response.ERROR, "failed to delete post")
		return
	}

	response.Success(c, nil)
}

// ListPosts 获取文章列表
// @Summary 获取文章列表
// @Tags 文章
// @Produce json
// @Param page query int false "页码" default(1)
// @Param pageSize query int false "每页数量" default(10)
// @Param status query int false "状态" default(1)
// @Success 200 {object} response.Response
// @Router /api/v1/posts [get]
func (h *PostHandler) ListPosts(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "10"))
	statusStr := c.DefaultQuery("status", "1")

	var status *int
	if statusStr != "" {
		s, _ := strconv.Atoi(statusStr)
		status = &s
	}

	posts, total, err := h.postRepo.List(page, pageSize, status)
	if err != nil {
		response.Error(c, response.ERROR, "failed to get posts")
		return
	}

	response.PageSuccess(c, posts, total, page, pageSize)
}

// SearchPosts 搜索文章
// @Summary 搜索文章
// @Tags 文章
// @Produce json
// @Param q query string true "搜索关键词"
// @Param page query int false "页码" default(1)
// @Param pageSize query int false "每页数量" default(10)
// @Success 200 {object} response.Response
// @Router /api/v1/posts/search [get]
func (h *PostHandler) SearchPosts(c *gin.Context) {
	keyword := c.Query("q")
	if keyword == "" {
		response.Error(c, response.ERROR, "keyword is required")
		return
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "10"))

	posts, total, err := h.postRepo.Search(keyword, page, pageSize)
	if err != nil {
		response.Error(c, response.ERROR, "failed to search posts")
		return
	}

	response.PageSuccess(c, posts, total, page, pageSize)
}

// ListPostsByTag 根据标签获取文章列表
// @Summary 根据标签获取文章列表
// @Tags 文章
// @Produce json
// @Param tagSlug path string true "标签 slug"
// @Param page query int false "页码" default(1)
// @Param pageSize query int false "每页数量" default(10)
// @Success 200 {object} response.Response
// @Router /api/v1/tags/:tagSlug/posts [get]
func (h *PostHandler) ListPostsByTag(c *gin.Context) {
	tagSlug := c.Param("tagSlug")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "10"))

	posts, total, err := h.postRepo.ListByTag(tagSlug, page, pageSize)
	if err != nil {
		response.Error(c, response.ERROR, "failed to get posts by tag")
		return
	}

	response.PageSuccess(c, posts, total, page, pageSize)
}

// ListPostsByCategory 根据分类获取文章列表
// @Summary 根据分类获取文章列表
// @Tags 文章
// @Produce json
// @Param categorySlug path string true "分类 slug"
// @Param page query int false "页码" default(1)
// @Param pageSize query int false "每页数量" default(10)
// @Success 200 {object} response.Response
// @Router /api/v1/categories/:categorySlug/posts [get]
func (h *PostHandler) ListPostsByCategory(c *gin.Context) {
	categorySlug := c.Param("categorySlug")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "10"))

	posts, total, err := h.postRepo.ListByCategory(categorySlug, page, pageSize)
	if err != nil {
		response.Error(c, response.ERROR, "failed to get posts by category")
		return
	}

	response.PageSuccess(c, posts, total, page, pageSize)
}

// generateSlug 生成 slug
func generateSlug(name string) string {
	// 转换为小写
	slug := strings.ToLower(name)
	// 移除特殊字符
	reg := regexp.MustCompile(`[^\w\s-]`)
	slug = reg.ReplaceAllString(slug, "")
	// 替换空格和多个连字符为单个连字符
	reg = regexp.MustCompile(`[-\s]+`)
	slug = reg.ReplaceAllString(slug, "-")
	return strings.TrimSpace(slug)
}
