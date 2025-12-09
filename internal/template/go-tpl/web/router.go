package web

import (
	"{{.ModulePath}}/infra/monitor"
	"{{.ModulePath}}/infra/logger"
	"{{.ModulePath}}/logic/auth"
	"{{.ModulePath}}/logic/user"
	"{{.ModulePath}}/web/base"
	"{{.ModulePath}}/web/middleware"
	authHandler "{{.ModulePath}}/web/rest/auth"
	userHandler "{{.ModulePath}}/web/rest/user"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type App struct {
	engine *gin.Engine
}

func New() *App {
	gin.SetMode(gin.ReleaseMode)
	engine := gin.New()

	// 中间件
	engine.Use(gin.Logger())
	engine.Use(gin.Recovery())
	engine.Use(corsMiddleware())

	// Prometheus metrics
	engine.GET("/metrics", gin.WrapH(promhttp.Handler()))

	// 健康检查
	engine.GET("/health", func(c *gin.Context) {
		base.Success(c, gin.H{
			"status": "ok",
			"time":   time.Now().Unix(),
		})
	})

	// API路由
	api := engine.Group("/api")
	{
		// 认证路由（不需要token）
		authSvc := auth.NewService()
		authH := authHandler.NewHandler(authSvc)
		authGroup := api.Group("/auth")
		{
			authGroup.POST("/register", authH.Register)
			authGroup.POST("/login", authH.Login)
			authGroup.POST("/refresh", authH.RefreshToken)
		}

		// 需要认证的路由
		api.Use(middleware.AuthMiddleware())
		{
			// 用户路由
			userSvc := user.NewService()
			userH := userHandler.NewHandler(userSvc)
			userGroup := api.Group("/user")
			{
				userGroup.GET("/profile", userH.GetProfile)
				userGroup.GET("/:id", userH.GetUser)
			}
		}
	}

	return &App{
		engine: engine,
	}
}

func (a *App) Run(addr string) {
	logger.Log.Info().Msgf("Server starting on %s", addr)
	a.engine.Run(addr)
}

func corsMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}

		c.Next()
	}
}