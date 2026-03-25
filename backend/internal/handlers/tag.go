package handlers

import (
	"strconv"

	"github.com/cworld1/aniya-blog/backend/internal/models"
	"github.com/cworld1/aniya-blog/backend/internal/repository"
	"github.com/cworld1/aniya-blog/backend/pkg/response"
	"github.com/gin-gonic/gin"
)

// TagHandler 标签处理器
type TagHandler struct {
	tagRepo *repository.TagRepository
}

// NewTagHandler 创建标签处理器
func NewTagHandler(tagRepo *repository.TagRepository) *TagHandler {
	return &TagHandler{
		tagRepo: tagRepo,
	}
}

// ListTags 获取标签列表
// @Summary 获取标签列表
// @Tags 标签
// @Produce json
// @Param page query int false "页码" default(1)
// @Param pageSize query int false "每页数量" default(20)
// @Success 200 {object} response.Response
// @Router /api/v1/tags [get]
func (h *TagHandler) ListTags(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "20"))

	tags, total, err := h.tagRepo.List(page, pageSize)
	if err != nil {
		response.Error(c, response.ERROR, "failed to get tags")
		return
	}

	response.PageSuccess(c, tags, total, page, pageSize)
}

// GetAllTags 获取所有标签
// @Summary 获取所有标签
// @Tags 标签
// @Produce json
// @Success 200 {object} response.Response
// @Router /api/v1/tags/all [get]
func (h *TagHandler) GetAllTags(c *gin.Context) {
	tags, err := h.tagRepo.GetAll()
	if err != nil {
		response.Error(c, response.ERROR, "failed to get all tags")
		return
	}

	response.Success(c, tags)
}

// GetTag 获取标签详情
// @Summary 获取标签详情
// @Tags 标签
// @Produce json
// @Param slug path string true "标签 slug"
// @Success 200 {object} response.Response
// @Router /api/v1/tags/:slug [get]
func (h *TagHandler) GetTag(c *gin.Context) {
	slug := c.Param("slug")

	tag, err := h.tagRepo.FindBySlug(slug)
	if err != nil {
		response.Error(c, response.ERROR, "tag not found")
		return
	}

	response.Success(c, tag)
}

// CreateTagRequest 创建标签请求
type CreateTagRequest struct {
	Name        string `json:"name" binding:"required"`
	Slug        string `json:"slug"`
	Description string `json:"description"`
}

// CreateTag 创建标签
// @Summary 创建标签
// @Tags 标签
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body CreateTagRequest true "标签信息"
// @Success 200 {object} response.Response
// @Router /api/v1/tags [post]
func (h *TagHandler) CreateTag(c *gin.Context) {
	var req CreateTagRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, response.ERROR, "invalid request")
		return
	}

	tag := &models.Tag{
		Name:        req.Name,
		Slug:        req.Slug,
		Description: req.Description,
	}

	if err := h.tagRepo.Create(tag); err != nil {
		response.Error(c, response.ERROR, "failed to create tag")
		return
	}

	response.Success(c, tag)
}

// UpdateTag 更新标签
// @Summary 更新标签
// @Tags 标签
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param slug path string true "标签 slug"
// @Param request body UpdateTagRequest true "标签信息"
// @Success 200 {object} response.Response
// @Router /api/v1/tags/:slug [put]
func (h *TagHandler) UpdateTag(c *gin.Context) {
	slug := c.Param("slug")

	tag, err := h.tagRepo.FindBySlug(slug)
	if err != nil {
		response.Error(c, response.ERROR, "tag not found")
		return
	}

	var req UpdateTagRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, response.ERROR, "invalid request")
		return
	}

	if req.Name != "" {
		tag.Name = req.Name
	}
	if req.Description != "" {
		tag.Description = req.Description
	}

	if err := h.tagRepo.Update(tag); err != nil {
		response.Error(c, response.ERROR, "failed to update tag")
		return
	}

	response.Success(c, tag)
}

// DeleteTag 删除标签
// @Summary 删除标签
// @Tags 标签
// @Produce json
// @Security BearerAuth
// @Param slug path string true "标签 slug"
// @Success 200 {object} response.Response
// @Router /api/v1/tags/:slug [delete]
func (h *TagHandler) DeleteTag(c *gin.Context) {
	slug := c.Param("slug")

	tag, err := h.tagRepo.FindBySlug(slug)
	if err != nil {
		response.Error(c, response.ERROR, "tag not found")
		return
	}

	if err := h.tagRepo.Delete(tag.ID); err != nil {
		response.Error(c, response.ERROR, "failed to delete tag")
		return
	}

	response.Success(c, nil)
}

// UpdateTagRequest 更新标签请求
type UpdateTagRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}
