package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"backend/audit"
	"backend/auth"
	"backend/config"
	"backend/controllers"
	"backend/database"
	"backend/middleware"
	"backend/realtime"
	"backend/repositories"
	"backend/services"
	"backend/telemetry"

	"github.com/gin-gonic/gin"
	ginSwagger "github.com/swaggo/gin-swagger"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"

	_ "backend/docs"

	swaggerFiles "github.com/swaggo/files"
)

// @title           Engineering Drawing QC System API
// @version         1.0
// @description     This is a quality control system for managing engineering drawings through a review workflow.
// @termsOfService  http://swagger.io/terms/

// @host      localhost:8081
// @BasePath  /api/v1

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization

func main() {
	cfg := config.LoadConfig()

	// Initialize Telemetry
	tp, err := telemetry.InitTracer()
	if err != nil {
		log.Fatalf("Failed to initialize tracer: %v", err)
	}
	defer telemetry.Shutdown(tp)

	// Initialize Database
	database.InitDB(cfg.DBURL)

	// Initialize Services (Dependency Injection)
	realtimeService := realtime.New(cfg)
	defer realtimeService.Shutdown()

	auditService := audit.New(cfg)
	defer auditService.Shutdown()

	// Initialize Repositories
	userRepo := repositories.NewUserRepository(database.DB)
	drawingRepo := repositories.NewDrawingRepository(database.DB)

	// Initialize Casbin
	auth.InitCasbin(database.DB)

	drawingService := services.NewDrawing(drawingRepo, auditService, realtimeService)

	// Initialize Controllers
	authCtrl := controllers.NewAuth(userRepo)
	drawingCtrl := controllers.NewDrawing(drawingRepo, drawingService)
	eventCtrl := controllers.NewEvent(realtimeService)

	r := gin.Default()

	// Middleware
	r.Use(otelgin.Middleware("qc-system"))
	r.Use(gin.Recovery())
	r.Use(CORSMiddleware())

	// Health Check
	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "UP", "time": time.Now()})
	})

	// Public Routes
	api := r.Group("/api/v1")
	{
		api.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
		api.POST("/register", authCtrl.Register)
		api.POST("/login", authCtrl.Login)
	}

	// Protected Routes
	protected := api.Group("/")
	protected.Use(middleware.AuthMiddleware())
	{
		// SSE Events (Real-time)
		protected.GET("/events", eventCtrl.StreamEvents)

		// Drawings
		drawings := protected.Group("/drawings")
		{
			drawings.GET("", middleware.RBACMiddleware("drawings", "view"), drawingCtrl.GetDrawings)
			drawings.POST("", middleware.RBACMiddleware("drawings", "create"), drawingCtrl.CreateDrawing)
			drawings.POST("/:id/claim", middleware.RBACMiddleware("drawings", "claim"), drawingCtrl.ClaimDrawing)
			drawings.POST("/:id/submit", middleware.RBACMiddleware("drawings", "submit"), drawingCtrl.SubmitDrawing)
			drawings.POST("/:id/release", middleware.RBACMiddleware("drawings", "release"), drawingCtrl.ReleaseDrawing)
			drawings.POST("/:id/reject", middleware.RBACMiddleware("drawings", "reject"), drawingCtrl.RejectDrawing)
		}
	}

	// Server setup
	srv := &http.Server{
		Addr:    ":" + cfg.Port,
		Handler: r,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	log.Printf("Server started on port %s", cfg.Port)

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server forced to shutdown:", err)
	}

	log.Println("Server exiting")
}

func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}

		c.Next()
	}
}
