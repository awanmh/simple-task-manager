package main

import (
    "context"
    "fmt"
    "log"
    "os"
    "time"

    "simple-task-manager/internal/core/usecase"
    handler "simple-task-manager/internal/infra/delivery/http"
    "simple-task-manager/internal/infra/delivery/middleware"
    "simple-task-manager/internal/infra/repository"
    "simple-task-manager/internal/infra/scheduler"

    "github.com/gin-contrib/cors"
    "github.com/gin-gonic/gin"
    "github.com/jackc/pgx/v5/pgxpool"
    "github.com/joho/godotenv"
    "github.com/robfig/cron/v3"

    // Swagger imports
    _ "simple-task-manager/docs" // generated docs folder
    swaggerFiles "github.com/swaggo/files"
    ginSwagger "github.com/swaggo/gin-swagger"
)

// @title           Simple Task Manager API
// @version         1.0
// @description     REST API untuk manajemen tugas dengan fitur Auth, Subtasks, dan Recurring.
// @termsOfService  http://swagger.io/terms/
//
// @contact.name    API Support
// @contact.email   support@swagger.io
//
// @license.name    Apache 2.0
// @license.url     http://www.apache.org/licenses/LICENSE-2.0.html
//
// @host            localhost:8080
// @BasePath        /
//
// @securityDefinitions.apikey BearerAuth
// @in              header
// @name            Authorization

// Config holds environment-backed configuration.
type Config struct {
    DBURL     string
    JWTSecret string
    Port      string
    Timeout   time.Duration
}

func main() {
    cfg := mustLoadConfig()
    dbPool := mustInitDB(cfg.DBURL)
    defer dbPool.Close()

    startCron(dbPool)

    userRepo := repository.NewUserRepository(dbPool)
    taskRepo := repository.NewTaskRepository(dbPool)

    userUseCase := usecase.NewUserUsecase(userRepo, cfg.Timeout, cfg.JWTSecret)
    taskUseCase := usecase.NewTaskUsecase(taskRepo, cfg.Timeout)

    userHandler := &handler.UserHandler{UserUseCase: userUseCase}
    taskHandler := &handler.TaskHandler{TaskUseCase: taskUseCase}

    r := setupRouter(cfg.JWTSecret, userHandler, taskHandler)

    log.Printf("Server running on port %s", cfg.Port)
    if err := r.Run(":" + cfg.Port); err != nil {
        log.Fatalf("Failed to run server: %v", err)
    }
}

// mustLoadConfig loads .env (if present) and returns validated config.
func mustLoadConfig() Config {
    if err := godotenv.Load(); err != nil {
        log.Println("No .env file found, using default env variables")
    }

    dbURL := os.Getenv("DB_URL")
    if dbURL == "" {
        log.Fatal("DB_URL is not set in .env")
    }

    jwtSecret := os.Getenv("JWT_SECRET")
    if jwtSecret == "" {
        log.Fatal("JWT_SECRET is not set in .env")
    }

    port := os.Getenv("PORT")
    if port == "" {
        port = "8080"
    }

    return Config{
        DBURL:     dbURL,
        JWTSecret: jwtSecret,
        Port:      port,
        Timeout:   2 * time.Second,
    }
}

// mustInitDB initializes a pgx pool or fatals.
func mustInitDB(dbURL string) *pgxpool.Pool {
    dbPool, err := pgxpool.New(context.Background(), dbURL)
    if err != nil {
        log.Fatalf("Unable to connect to database: %v\n", err)
    }
    return dbPool
}

// startCron initializes and starts the recurring task scheduler.
func startCron(dbPool *pgxpool.Pool) {
    taskScheduler := scheduler.NewTaskScheduler(dbPool)

    c := cron.New()
    _, _ = c.AddFunc("@every 1m", func() {
        taskScheduler.ProcessRecurringTasks()
    })
    c.Start()

    fmt.Println(">>> Cron Scheduler Started!")
}

// setupRouter wires middlewares, routes, and swagger.
func setupRouter(jwtSecret string, userHandler *handler.UserHandler, taskHandler *handler.TaskHandler) *gin.Engine {
    r := gin.Default()
    r.Use(corsMiddleware())

    // Public routes
    r.POST("/register", userHandler.Register)
    r.POST("/login", userHandler.Login)

    // Protected task routes
    protected := r.Group("/tasks")
    protected.Use(middleware.AuthMiddleware(jwtSecret))
    {
        protected.POST("/", taskHandler.Create)
        protected.GET("/", taskHandler.Fetch)
        protected.PUT("/:id", taskHandler.UpdateStatus)
        protected.DELETE("/:id", taskHandler.Delete)
        protected.POST("/:id/subtasks", taskHandler.AddSubtask)
    }

    // Subtask routes (protected per-handler)
    r.PUT("/subtasks/:sub_id", middleware.AuthMiddleware(jwtSecret), taskHandler.ToggleSubtask)
    r.DELETE("/subtasks/:sub_id", middleware.AuthMiddleware(jwtSecret), taskHandler.DeleteSubtask)

    // Swagger
    r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

    return r
}

// corsMiddleware returns a configured CORS middleware.
func corsMiddleware() gin.HandlerFunc {
    config := cors.DefaultConfig()
    config.AllowAllOrigins = true
    config.AllowMethods = []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"}
    config.AllowHeaders = []string{"Origin", "Content-Type", "Authorization", "Accept"}
    return cors.New(config)
}
