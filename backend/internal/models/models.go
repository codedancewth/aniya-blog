package models

import (
	"time"

	"gorm.io/gorm"
)

// BaseModel 基础模型
type BaseModel struct {
	ID        uint           `gorm:"primarykey" json:"id"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

// User 用户模型
type User struct {
	BaseModel
	Username    string     `gorm:"uniqueIndex;size:50;not null" json:"username"`
	Password    string     `gorm:"size:255;not null" json:"-"`
	Email       string     `gorm:"uniqueIndex;size:100" json:"email"`
	Avatar      string     `gorm:"size:255" json:"avatar"`
	Nickname    string     `gorm:"size:50" json:"nickname"`
	Role        string     `gorm:"size:20;default:'user'" json:"role"` // admin, editor, user
	Status      int        `gorm:"default:1" json:"status"`            // 1: active, 0: inactive
	LastLoginAt *time.Time `json:"last_login_at"`
	LastLoginIP string     `gorm:"size:50" json:"last_login_ip"`
}

// Post 文章模型
type Post struct {
	BaseModel
	Title        string     `gorm:"size:200;not null" json:"title"`
	Slug         string     `gorm:"uniqueIndex;size:200;not null" json:"slug"`
	Description  string     `gorm:"size:500" json:"description"`
	Content      string     `gorm:"type:text" json:"content"`
	ContentHTML  string     `gorm:"type:text" json:"content_html"`
	CoverImage   string     `gorm:"size:255" json:"cover_image"`
	AuthorID     uint       `gorm:"not null" json:"author_id"`
	Author       User       `gorm:"foreignKey:AuthorID" json:"author,omitempty"`
	Status       int        `gorm:"default:1" json:"status"` // 1: published, 0: draft, 2: archived
	PublishedAt  *time.Time `json:"published_at"`
	ViewCount    int64      `gorm:"default:0" json:"view_count"`
	CommentCount int64      `gorm:"default:0" json:"comment_count"`
	LikeCount    int64      `gorm:"default:0" json:"like_count"`
	Tags         []Tag      `gorm:"many2many:post_tags;" json:"tags,omitempty"`
	CategoryID   *uint      `json:"category_id"`
	Category     *Category  `gorm:"foreignKey:CategoryID" json:"category,omitempty"`
	Language     string     `gorm:"size:10;default:'zh-CN'" json:"language"`
	IsTop        bool       `gorm:"default:false" json:"is_top"`
	CustomData   string     `gorm:"type:text" json:"custom_data"` // 存储 frontmatter 数据
}

// Category 分类模型
type Category struct {
	BaseModel
	Name        string     `gorm:"uniqueIndex;size:50;not null" json:"name"`
	Slug        string     `gorm:"uniqueIndex;size:50;not null" json:"slug"`
	Description string     `gorm:"size:255" json:"description"`
	ParentID    *uint      `json:"parent_id"`
	Parent      *Category  `gorm:"foreignKey:ParentID" json:"parent,omitempty"`
	Children    []Category `gorm:"foreignKey:ParentID" json:"children,omitempty"`
	SortOrder   int        `gorm:"default:0" json:"sort_order"`
	PostCount   int64      `gorm:"default:0" json:"post_count"`
}

// Tag 标签模型
type Tag struct {
	BaseModel
	Name        string `gorm:"uniqueIndex;size:50;not null" json:"name"`
	Slug        string `gorm:"uniqueIndex;size:50;not null" json:"slug"`
	Description string `gorm:"size:255" json:"description"`
	PostCount   int64  `gorm:"default:0" json:"post_count"`
}

// Comment 评论模型
type Comment struct {
	BaseModel
	Content     string    `gorm:"type:text;not null" json:"content"`
	PostID      uint      `gorm:"not null;index" json:"post_id"`
	Post        Post      `gorm:"foreignKey:PostID" json:"post,omitempty"`
	UserID      *uint     `json:"user_id"`
	User        *User     `gorm:"foreignKey:UserID" json:"user,omitempty"`
	ParentID    *uint     `json:"parent_id"`
	Parent      *Comment  `gorm:"foreignKey:ParentID" json:"parent,omitempty"`
	Replies     []Comment `gorm:"foreignKey:ParentID" json:"replies,omitempty"`
	AuthorName  string    `gorm:"size:50" json:"author_name"`
	AuthorEmail string    `gorm:"size:100" json:"author_email"`
	AuthorURL   string    `gorm:"size:255" json:"author_url"`
	AuthorIP    string    `gorm:"size:50" json:"author_ip"`
	Agent       string    `gorm:"size:255" json:"agent"`   // User-Agent
	Status      int       `gorm:"default:1" json:"status"` // 1: approved, 0: pending, 2: spam
	LikeCount   int64     `gorm:"default:0" json:"like_count"`
	IsAdmin     bool      `gorm:"default:false" json:"is_admin"`
}

// PageView 页面浏览记录
type PageView struct {
	BaseModel
	Path      string `gorm:"size:255;not null;index" json:"path"`
	IP        string `gorm:"size:50" json:"ip"`
	UserAgent string `gorm:"size:255" json:"user_agent"`
	Referer   string `gorm:"size:255" json:"referer"`
	Country   string `gorm:"size:50" json:"country"`
	Province  string `gorm:"size:50" json:"province"`
	City      string `gorm:"size:50" json:"city"`
}

// Link 友情链接
type Link struct {
	BaseModel
	Name        string `gorm:"size:100;not null" json:"name"`
	URL         string `gorm:"size:255;not null" json:"url"`
	Logo        string `gorm:"size:255" json:"logo"`
	Description string `gorm:"size:255" json:"description"`
	Status      int    `gorm:"default:1" json:"status"` // 1: active, 0: inactive
	SortOrder   int    `gorm:"default:0" json:"sort_order"`
}

// Config 站点配置
type Config struct {
	BaseModel
	Key   string `gorm:"uniqueIndex;size:100;not null" json:"key"`
	Value string `gorm:"type:text" json:"value"`
	Type  string `gorm:"size:20" json:"type"` // string, number, boolean, json
}

// Like 点赞记录
type Like struct {
	BaseModel
	PostID    uint   `gorm:"not null;index" json:"post_id"`
	UserID    *uint  `json:"user_id"`
	IP        string `gorm:"size:50" json:"ip"`
	UserAgent string `gorm:"size:255" json:"user_agent"`
}

// 表名自定义
func (Post) TableName() string {
	return "posts"
}

func (Category) TableName() string {
	return "categories"
}

func (Tag) TableName() string {
	return "tags"
}

func (Comment) TableName() string {
	return "comments"
}

func (PageView) TableName() string {
	return "page_views"
}

func (Link) TableName() string {
	return "links"
}

func (Config) TableName() string {
	return "configs"
}

func (Like) TableName() string {
	return "likes"
}
