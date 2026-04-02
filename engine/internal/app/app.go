package app

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/lieying/engine/internal/config"
	"github.com/lieying/engine/internal/handlers"
	"github.com/lieying/engine/internal/models"
	"github.com/lieying/engine/internal/scheduler"
	"github.com/lieying/engine/pkg/ai"
	"github.com/lieying/engine/pkg/database"
	"github.com/lieying/engine/pkg/mq"
	"github.com/lieying/engine/pkg/network"
	"github.com/lieying/engine/pkg/redis"
	"gorm.io/gorm"
)

type Application struct {
	config         *config.Config
	router         *gin.Engine
	db             *gorm.DB
	redis          *redis.Client
	mq             *mq.Client
	networkManager *network.NetworkManager
	aiManager      *ai.AIClientManager
	scheduler      *scheduler.TaskScheduler
	server         *http.Server
}

func New(cfg *config.Config) (*Application, error) {
	if cfg.Server.Mode == "release" {
		gin.SetMode(gin.ReleaseMode)
	}

	db, err := database.New(cfg.Database)
	if err != nil {
		return nil, fmt.Errorf("database init failed: %w", err)
	}

	if err := db.AutoMigrate(
		&models.Task{},
		&models.Vulnerability{},
		&models.Target{},
		&models.POC{},
		&models.ScanResult{},
	); err != nil {
		return nil, fmt.Errorf("migrate failed: %w", err)
	}

	redisClient, err := redis.New(cfg.Redis)
	if err != nil {
		return nil, fmt.Errorf("redis init failed: %w", err)
	}

	mqClient, err := mq.New(cfg.RabbitMQ)
	if err != nil {
		return nil, fmt.Errorf("mq init failed: %w", err)
	}

	networkManager, err := network.NewNetworkManager(cfg.Network)
	if err != nil {
		return nil, fmt.Errorf("network manager init failed: %w", err)
	}

	aiManager := ai.NewAIClientManager(cfg.AI, networkManager.GetClient())

	taskScheduler := scheduler.New(db, redisClient, mqClient, cfg.Engine)

	router := gin.Default()
	registerRoutes(router, db, redisClient, aiManager, taskScheduler)

	server := &http.Server{
		Addr:    fmt.Sprintf(":%d", cfg.Server.Port),
		Handler: router,
	}

	return &Application{
		config:         cfg,
		router:         router,
		db:             db,
		redis:          redisClient,
		mq:             mqClient,
		networkManager: networkManager,
		aiManager:      aiManager,
		scheduler:      taskScheduler,
		server:         server,
	}, nil
}

func (a *Application) Start() error {
	go a.scheduler.Start()
	return a.server.ListenAndServe()
}

func (a *Application) Stop() error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	a.scheduler.Stop()
	if err := a.mq.Close(); err != nil {
		return err
	}
	return a.server.Shutdown(ctx)
}

func registerRoutes(r *gin.Engine, db *gorm.DB, redis *redis.Client, aiManager *ai.AIClientManager, scheduler *scheduler.TaskScheduler) {
	api := r.Group("/api/v1")
	{
		taskHandler := handlers.NewTaskHandler(db, scheduler)
		tasks := api.Group("/tasks")
		{
			tasks.POST("", taskHandler.Create)
			tasks.GET("", taskHandler.List)
			tasks.GET("/:id", taskHandler.Get)
			tasks.POST("/:id/start", taskHandler.Start)
			tasks.POST("/:id/stop", taskHandler.Stop)
		}

		vulnHandler := handlers.NewVulnHandler(db)
		vulns := api.Group("/vulnerabilities")
		{
			vulns.GET("", vulnHandler.List)
			vulns.GET("/:id", vulnHandler.Get)
		}

		pocHandler := handlers.NewPOCHandler(db)
		pocs := api.Group("/pocs")
		{
			pocs.GET("", pocHandler.List)
			pocs.POST("", pocHandler.Create)
			pocs.GET("/:id", pocHandler.Get)
		}

		aiHandler := handlers.NewAIHandler(aiManager)
		ai := api.Group("/ai")
		{
			ai.POST("/poc/generate", aiHandler.GeneratePOC)
			ai.POST("/report/generate", aiHandler.GenerateReport)
			ai.POST("/attack-path/analyze", aiHandler.AnalyzeAttackPath)
			ai.GET("/status", aiHandler.Status)
		}
	}
}
