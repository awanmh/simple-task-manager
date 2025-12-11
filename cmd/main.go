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

// @contact.name    API Support
// @contact.email   support@swagger.io

// @license.name    Apache 2.0
// @license.url     http://www.apache.org/licenses/LICENSE-2.0.html

// @host            localhost:8080
// @BasePath        /

// @securityDefinitions.apikey BearerAuth
// @in              header
// @name            Authorization
func main() {
    // 1. Load .env file
    if err := godotenv.Load(); err != nil {
        log.Println("No .env file found, using default env variables")
    }

    // 2. Ambil Config dari ENV
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

    timeout := 2 * time.Second

    // 3. Koneksi Database
    dbPool, err := pgxpool.New(context.Background(), dbURL)
    if err != nil {
        log.Fatalf("Unable to connect to database: %v\n", err)
    }
    defer dbPool.Close()

    // --- SETUP CRON JOB (Background Worker) ---
    taskScheduler := scheduler.NewTaskScheduler(dbPool)
    c := cron.New()
    c.AddFunc("@every 1m", func() {
        taskScheduler.ProcessRecurringTasks()
    })
    c.Start()
    fmt.Println(">>> Cron Scheduler Started!")

    // 4. Setup Layers
    userRepo := repository.NewUserRepository(dbPool)
    taskRepo := repository.NewTaskRepository(dbPool)

    userUseCase := usecase.NewUserUsecase(userRepo, timeout, jwtSecret)
    taskUseCase := usecase.NewTaskUsecase(taskRepo, timeout)

    userHandler := &handler.UserHandler{UserUseCase: userUseCase}
    taskHandler := &handler.TaskHandler{TaskUseCase: taskUseCase}

    // 5. Setup Router
    r := gin.Default()

    // --- CORS ---
    config := cors.DefaultConfig()
    config.AllowAllOrigins = true
    config.AllowMethods = []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"}
    config.AllowHeaders = []string{"Origin", "Content-Type", "Authorization", "Accept"}
    r.Use(cors.New(config))

    // Routes
    r.POST("/register", userHandler.Register)
    r.POST("/login", userHandler.Login)

    protected := r.Group("/tasks")
    protected.Use(middleware.AuthMiddleware(jwtSecret))
    {
        protected.POST("/", taskHandler.Create)
        protected.GET("/", taskHandler.Fetch)
        protected.PUT("/:id", taskHandler.UpdateStatus)
        protected.DELETE("/:id", taskHandler.Delete)
        protected.POST("/:id/subtasks", taskHandler.AddSubtask)
    }

    r.PUT("/subtasks/:sub_id", middleware.AuthMiddleware(jwtSecret), taskHandler.ToggleSubtask)
    r.DELETE("/subtasks/:sub_id", middleware.AuthMiddleware(jwtSecret), taskHandler.DeleteSubtask)

    // --- Swagger Route ---
    r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

    // 6. Run Server
    log.Printf("Server running on port %s", port)
    if err := r.Run(":" + port); err != nil {
        log.Fatalf("Failed to run server: %v", err)
    }
}
