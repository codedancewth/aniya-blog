package handlers

import (
	"strings"

	"github.com/cworld1/aniya-blog/backend/internal/models"
	"github.com/cworld1/aniya-blog/backend/internal/repository"
	"github.com/cworld1/aniya-blog/backend/pkg/response"
	"github.com/gin-gonic/gin"
)

// PageViewHandler 页面浏览处理器
type PageViewHandler struct {
	pageViewRepo *repository.PageViewRepository
	postRepo     *repository.PostRepository
}

// NewPageViewHandler 创建页面浏览处理器
func NewPageViewHandler(pageViewRepo *repository.PageViewRepository, postRepo *repository.PostRepository) *PageViewHandler {
	return &PageViewHandler{
		pageViewRepo: pageViewRepo,
		postRepo:     postRepo,
	}
}

// RecordPageViewRequest 记录页面浏览请求
type RecordPageViewRequest struct {
	Path string `json:"path" binding:"required"`
}

// RecordPageView 记录页面浏览
// @Summary 记录页面浏览
// @Tags 页面浏览
// @Accept json
// @Produce json
// @Param request body RecordPageViewRequest true "页面路径"
// @Success 200 {object} response.Response
// @Router /api/v1/pageviews [post]
func (h *PageViewHandler) RecordPageView(c *gin.Context) {
	var req RecordPageViewRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, response.ERROR, "invalid request")
		return
	}

	pageView := &models.PageView{
		Path:      req.Path,
		IP:        c.ClientIP(),
		UserAgent: c.GetHeader("User-Agent"),
		Referer:   c.GetHeader("Referer"),
	}

	if err := h.pageViewRepo.Create(pageView); err != nil {
		response.Error(c, response.ERROR, "failed to record page view")
		return
	}

	// 获取浏览次数
	count, _ := h.pageViewRepo.CountByPath(req.Path)

	response.Success(c, gin.H{
		"path":  req.Path,
		"count": count,
	})
}

// GetPageViewCount 获取页面浏览次数
// @Summary 获取页面浏览次数
// @Tags 页面浏览
// @Produce json
// @Param path query string true "页面路径"
// @Success 200 {object} response.Response
// @Router /api/v1/pageviews/count [get]
func (h *PageViewHandler) GetPageViewCount(c *gin.Context) {
	path := c.Query("path")
	if path == "" {
		response.Error(c, response.ERROR, "path is required")
		return
	}

	count, err := h.pageViewRepo.CountByPath(path)
	if err != nil {
		response.Error(c, response.ERROR, "failed to get page view count")
		return
	}

	response.Success(c, gin.H{
		"path":  path,
		"count": count,
	})
}

// GetPageViewCounts 批量获取页面浏览次数
// @Summary 批量获取页面浏览次数
// @Tags 页面浏览
// @Produce json
// @Param paths query string true "页面路径列表，逗号分隔"
// @Success 200 {object} response.Response
// @Router /api/v1/pageviews/counts [get]
func (h *PageViewHandler) GetPageViewCounts(c *gin.Context) {
	paths := c.Query("paths")
	if paths == "" {
		response.Error(c, response.ERROR, "paths is required")
		return
	}

	pathList := splitString(paths, ",")
	counts := make(map[string]int64)

	for _, path := range pathList {
		count, _ := h.pageViewRepo.CountByPath(path)
		counts[path] = count
	}

	response.Success(c, counts)
}

// GetSiteStats 获取站点统计信息
// @Summary 获取站点统计信息
// @Tags 页面浏览
// @Produce json
// @Success 200 {object} response.Response
// @Router /api/v1/pageviews/stats [get]
func (h *PageViewHandler) GetSiteStats(c *gin.Context) {
	totalViews, _ := h.pageViewRepo.GetTotalViews()
	totalPosts, _ := h.postRepo.List(1, 1, nil)
	totalPostCount, _ := h.pageViewRepo.GetTotalPosts()

	response.Success(c, gin.H{
		"total_views":     totalViews,
		"total_posts":     len(totalPosts),
		"total_post_count": totalPostCount,
	})
}

func splitString(s, sep string) []string {
	return strings.Split(s, sep)
}
