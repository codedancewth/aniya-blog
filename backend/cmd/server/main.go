package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"

	"github.com/cworld1/aniya-blog/backend/internal/config"
	"github.com/cworld1/aniya-blog/backend/internal/database"
	"github.com/cworld1/aniya-blog/backend/internal/handlers"
	"github.com/cworld1/aniya-blog/backend/internal/middleware"
	"github.com/cworld1/aniya-blog/backend/internal/repository"
	"github.com/cworld1/aniya-blog/backend/pkg/auth"
	"github.com/gin-gonic/gin"
	"github.com/rs/cors"
)

// @title Aniya Blog API
// @version 1.0.0
// @description Aniya Blog 后端 API 服务
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host localhost:8080
// @BasePath /api/v1

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description 使用 JWT 令牌进行认证，格式为 "Bearer {token}"
func main() {
	// 解析命令行参数
	configPath := flag.String("config", "", "config file path")
	flag.Parse()

	// 加载配置
	cfg := config.Load()

	// 初始化数据库
	if err := database.Init(&cfg.Database); err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	defer database.Close()

	// 初始化仓库
	userRepo := repository.NewUserRepository()
	postRepo := repository.NewPostRepository()
	tagRepo := repository.NewTagRepository()
	categoryRepo := repository.NewCategoryRepository()
	commentRepo := repository.NewCommentRepository()
	pageViewRepo := repository.NewPageViewRepository()

	// 初始化 JWT 管理器
	jwtManager := auth.NewJWTManager(cfg.JWT.Secret, cfg.JWT.ExpireTime)

	// 初始化处理器
	authHandler := handlers.NewAuthHandler(userRepo, jwtManager)
	postHandler := handlers.NewPostHandler(postRepo, tagRepo, categoryRepo)
	tagHandler := handlers.NewTagHandler(tagRepo)
	categoryHandler := handlers.NewCategoryHandler(categoryRepo)
	commentHandler := handlers.NewCommentHandler(commentRepo, postRepo)
	pageViewHandler := handlers.NewPageViewHandler(pageViewRepo, postRepo)

	// 设置 Gin 模式
	gin.SetMode(cfg.Server.Mode)

	// 创建 Gin 路由器
	r := gin.Default()

	// 应用 CORS 中间件
	c := cors.New(cors.Options{
		AllowedOrigins:   cfg.Server.AllowOrigins,
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "PATCH", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300,
	})

	// 应用中间件
	r.Use(func() gin.HandlerFunc {
		return func(c *gin.Context) {
			corsHandler := c.Handler()
			corsHandler.ServeHTTP(c.Writer, c.Request)
		}
	}())
	r.Use(middleware.Logger())
	r.Use(middleware.Recovery())

	// 健康检查
	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status": "ok",
		})
	})

	// API v1 路由组
	v1 := r.Group("/api/v1")
	{
		// 认证路由（无需登录）
		auth := v1.Group("/auth")
		{
			auth.POST("/login", authHandler.Login)
			auth.POST("/register", authHandler.Register)
		}

		// 公开路由
		public := v1.Group("")
		{
			// 文章
			public.GET("/posts", postHandler.ListPosts)
			public.GET("/posts/:id", postHandler.GetPost)
			public.GET("/posts/slug/:slug", postHandler.GetPostBySlug)
			public.GET("/posts/search", postHandler.SearchPosts)
			public.GET("/tags/:tagSlug/posts", postHandler.ListPostsByTag)
			public.GET("/categories/:categorySlug/posts", postHandler.ListPostsByCategory)

			// 标签
			public.GET("/tags", tagHandler.ListTags)
			public.GET("/tags/all", tagHandler.GetAllTags)
			public.GET("/tags/:slug", tagHandler.GetTag)

			// 分类
			public.GET("/categories", categoryHandler.ListCategories)
			public.GET("/categories/tree", categoryHandler.GetAllCategories)
			public.GET("/categories/:slug", categoryHandler.GetCategory)

			// 评论
			public.GET("/posts/:post_id/comments", commentHandler.ListCommentsByPost)

			// 页面浏览
			public.POST("/pageviews", pageViewHandler.RecordPageView)
			public.GET("/pageviews/count", pageViewHandler.GetPageViewCount)
			public.GET("/pageviews/counts", pageViewHandler.GetPageViewCounts)
			public.GET("/pageviews/stats", pageViewHandler.GetSiteStats)
		}

		// 需要认证的路由
		protected := v1.Group("")
		protected.Use(middleware.JWTAuth(jwtManager))
		{
			// 认证
			protected.GET("/auth/me", authHandler.GetCurrentUser)
			protected.POST("/auth/refresh", authHandler.RefreshToken)
			protected.POST("/auth/change-password", authHandler.ChangePassword)

			// 文章管理
			protected.POST("/posts", postHandler.CreatePost)
			protected.PUT("/posts/:id", postHandler.UpdatePost)
			protected.DELETE("/posts/:id", postHandler.DeletePost)

			// 标签管理（需要管理员权限）
			protected.POST("/tags", tagHandler.CreateTag)
			protected.PUT("/tags/:slug", tagHandler.UpdateTag)
			protected.DELETE("/tags/:slug", tagHandler.DeleteTag)

			// 分类管理（需要管理员权限）
			protected.POST("/categories", categoryHandler.CreateCategory)
			protected.PUT("/categories/:slug", categoryHandler.UpdateCategory)
			protected.DELETE("/categories/:slug", categoryHandler.DeleteCategory)

			// 评论管理
			protected.POST("/comments", commentHandler.CreateComment)
			protected.PUT("/comments/:id", commentHandler.UpdateComment)
			protected.DELETE("/comments/:id", commentHandler.DeleteComment)
			protected.POST("/comments/:id/like", commentHandler.LikeComment)
		}

		// 需要管理员权限的路由
		admin := v1.Group("/admin")
		admin.Use(middleware.JWTAuth(jwtManager), middleware.RequireRole("admin", "editor"))
		{
			// 管理员专用 API 可以在这里添加
			admin.GET("/users", func(c *gin.Context) {
				users, total, err := userRepo.List(1, 10)
				if err != nil {
					c.JSON(500, gin.H{"error": "failed to get users"})
					return
				}
				c.JSON(200, gin.H{"users": users, "total": total})
			})
		}
	}

	// 启动服务器
	addr := fmt.Sprintf(":%s", cfg.Server.Port)
	log.Printf("Starting server on %s", addr)

	// 使用 CORS 包装器启动
	log.Fatal(http.ListenAndServe(addr, c.Handler(r)))
}
