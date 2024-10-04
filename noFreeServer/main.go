package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"noFree/config"
	"noFree/database"
	"noFree/handlers"
	"noFree/middleware"
	"noFree/models"
	"noFree/repositories"
	"noFree/services"

	"github.com/gin-gonic/gin"
)

func main() {
	// 加载配置
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	// 初始化数据库连接
	if err := database.InitDB(cfg); err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	defer database.CloseDB()

	db := database.GetDB()

	// 自动迁移数据库结构
	if err := db.AutoMigrate(&models.User{}, &models.Device{}); err != nil {
		log.Fatalf("Failed to migrate database: %v", err)
	}

	// 初始化依赖
	userRepo := repositories.NewUserRepository(db)
	userService := services.NewUserService(userRepo)
	userHandler := handlers.NewUserHandler(userService, cfg)

	// 创建Gin路由
	r := gin.Default()

	// 设置路由
	setupRoutes(r, userHandler, cfg)

	// 创建 HTTP 服务器
	srv := &http.Server{
		Addr:    fmt.Sprintf(":%d", cfg.App.Port),
		Handler: r,
	}

	// 在 goroutine 中启动服务器
	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	// 等待中断信号以优雅地关闭服务器
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down server...")

	// 设置 5 秒的超时时间
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server forced to shutdown:", err)
	}

	log.Println("Server exiting")
}

func setupRoutes(r *gin.Engine, userHandler *handlers.UserHandler, cfg *config.Config) {
	// 公开路由
	r.POST("/api/register", userHandler.Register)
	r.POST("/api/login", userHandler.Login)

	// 认证路由组
	auth := r.Group("/api")
	auth.Use(middleware.AuthMiddleware(cfg), middleware.DeviceAuthMiddleware())
	{
		auth.POST("/validate", userHandler.ValidateDevice)
		auth.DELETE("/device", userHandler.RemoveDevice)
	}
}
