package handlers

import (
	"strconv"

	"github.com/cworld1/aniya-blog/backend/internal/models"
	"github.com/cworld1/aniya-blog/backend/internal/repository"
	"github.com/cworld1/aniya-blog/backend/pkg/response"
	"github.com/cworld1/aniya-blog/backend/pkg/validator"
	"github.com/gin-gonic/gin"
)

// CategoryHandler 分类处理器
type CategoryHandler struct {
	categoryRepo *repository.CategoryRepository
}

// NewCategoryHandler 创建分类处理器
func NewCategoryHandler(categoryRepo *repository.CategoryRepository) *CategoryHandler {
	return &CategoryHandler{
		categoryRepo: categoryRepo,
	}
}

// ListCategories 获取分类列表
// @Summary 获取分类列表
// @Tags 分类
// @Produce json
// @Param page query int false "页码" default(1)
// @Param pageSize query int false "每页数量" default(20)
// @Success 200 {object} response.Response
// @Router /api/v1/categories [get]
func (h *CategoryHandler) ListCategories(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "20"))

	categories, total, err := h.categoryRepo.List(page, pageSize)
	if err != nil {
		response.Error(c, response.ERROR, "failed to get categories")
		return
	}

	response.PageSuccess(c, categories, total, page, pageSize)
}

// GetAllCategories 获取所有分类（树形结构）
// @Summary 获取所有分类
// @Tags 分类
// @Produce json
// @Success 200 {object} response.Response
// @Router /api/v1/categories/tree [get]
func (h *CategoryHandler) GetAllCategories(c *gin.Context) {
	categories, err := h.categoryRepo.GetTree()
	if err != nil {
		response.Error(c, response.ERROR, "failed to get categories")
		return
	}

	response.Success(c, categories)
}

// GetCategory 获取分类详情
// @Summary 获取分类详情
// @Tags 分类
// @Produce json
// @Param slug path string true "分类 slug"
// @Success 200 {object} response.Response
// @Router /api/v1/categories/:slug [get]
func (h *CategoryHandler) GetCategory(c *gin.Context) {
	slug := c.Param("slug")

	category, err := h.categoryRepo.FindBySlug(slug)
	if err != nil {
		response.Error(c, response.ERROR, "category not found")
		return
	}

	response.Success(c, category)
}

// CreateCategoryRequest 创建分类请求
type CreateCategoryRequest struct {
	Name        string `json:"name" binding:"required"`
	Slug        string `json:"slug"`
	Description string `json:"description"`
	ParentID    *uint  `json:"parent_id"`
	SortOrder   int    `json:"sort_order"`
}

// CreateCategory 创建分类
// @Summary 创建分类
// @Tags 分类
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body CreateCategoryRequest true "分类信息"
// @Success 200 {object} response.Response
// @Router /api/v1/categories [post]
func (h *CategoryHandler) CreateCategory(c *gin.Context) {
	var req CreateCategoryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, response.ERROR, validator.GetErrorMsg(err))
		return
	}

	category := &models.Category{
		Name:        req.Name,
		Slug:        req.Slug,
		Description: req.Description,
		ParentID:    req.ParentID,
		SortOrder:   req.SortOrder,
	}

	if err := h.categoryRepo.Create(category); err != nil {
		response.Error(c, response.ERROR, "failed to create category")
		return
	}

	response.Success(c, category)
}

// UpdateCategory 更新分类
// @Summary 更新分类
// @Tags 分类
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param slug path string true "分类 slug"
// @Param request body UpdateCategoryRequest true "分类信息"
// @Success 200 {object} response.Response
// @Router /api/v1/categories/:slug [put]
func (h *CategoryHandler) UpdateCategory(c *gin.Context) {
	slug := c.Param("slug")

	category, err := h.categoryRepo.FindBySlug(slug)
	if err != nil {
		response.Error(c, response.ERROR, "category not found")
		return
	}

	var req UpdateCategoryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, response.ERROR, validator.GetErrorMsg(err))
		return
	}

	if req.Name != "" {
		category.Name = req.Name
	}
	if req.Slug != "" {
		category.Slug = req.Slug
	}
	if req.Description != "" {
		category.Description = req.Description
	}
	if req.ParentID != nil {
		category.ParentID = req.ParentID
	}
	category.SortOrder = req.SortOrder

	if err := h.categoryRepo.Update(category); err != nil {
		response.Error(c, response.ERROR, "failed to update category")
		return
	}

	response.Success(c, category)
}

// DeleteCategory 删除分类
// @Summary 删除分类
// @Tags 分类
// @Produce json
// @Security BearerAuth
// @Param slug path string true "分类 slug"
// @Success 200 {object} response.Response
// @Router /api/v1/categories/:slug [delete]
func (h *CategoryHandler) DeleteCategory(c *gin.Context) {
	slug := c.Param("slug")

	category, err := h.categoryRepo.FindBySlug(slug)
	if err != nil {
		response.Error(c, response.ERROR, "category not found")
		return
	}

	if err := h.categoryRepo.Delete(category.ID); err != nil {
		response.Error(c, response.ERROR, "failed to delete category")
		return
	}

	response.Success(c, nil)
}

// UpdateCategoryRequest 更新分类请求
type UpdateCategoryRequest struct {
	Name        string `json:"name"`
	Slug        string `json:"slug"`
	Description string `json:"description"`
	ParentID    *uint  `json:"parent_id"`
	SortOrder   int    `json:"sort_order"`
}
