package http

import (
    "fmt"
    "net/http"
    "strconv"

    "simple-task-manager/internal/core/domain"
    "simple-task-manager/internal/core/usecase"

    "github.com/gin-gonic/gin"
)

// --- User Handler ---

type UserHandler struct {
    UserUseCase *usecase.UserUsecase
}

// Register user baru
func (h *UserHandler) Register(c *gin.Context) {
    var req struct {
        Name     string `json:"name"`
        Email    string `json:"email"`
        Password string `json:"password"`
    }

    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    fmt.Printf("\n>>> DEBUG INPUT REGISTER: Name=%s, Email=%s, PassLen=%d\n\n",
        req.Name, req.Email, len(req.Password))

    user := domain.User{
        Name:     req.Name,
        Email:    req.Email,
        Password: req.Password,
    }

    if err := h.UserUseCase.Register(c.Request.Context(), &user); err != nil {
        fmt.Printf("\n>>> DEBUG ERROR REGISTER: %v\n\n", err)
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    c.JSON(http.StatusCreated, gin.H{
        "message": "User registered successfully",
        "user_id": user.ID,
    })
}

// Login godoc
// @Summary      Login User
// @Description  Masuk dengan email dan password untuk mendapatkan Token JWT
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        request body object{email=string,password=string} true "Login Credentials"
// @Success      200  {object}  map[string]interface{}
// @Failure      400  {object}  map[string]interface{}
// @Failure      401  {object}  map[string]interface{}
// @Router       /login [post]
func (h *UserHandler) Login(c *gin.Context) {
    var req struct {
        Email    string `json:"email"`
        Password string `json:"password"`
    }

    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    res, err := h.UserUseCase.Login(c.Request.Context(), req.Email, req.Password)
    if err != nil {
        fmt.Printf("\n>>> DEBUG ERROR LOGIN: %v\n\n", err)
        c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
        return
    }

    c.JSON(http.StatusOK, res)
}

//
// --- Task Handler ---
//

type TaskHandler struct {
    TaskUseCase *usecase.TaskUsecase
}

// CreateTask godoc
// @Summary      Create New Task
// @Description  Membuat tugas baru dengan opsi Priority dan Recurrence
// @Tags         tasks
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        request body domain.Task true "Task Data"
// @Success      201  {object}  domain.Task
// @Failure      400  {object}  map[string]interface{}
// @Router       /tasks/ [post]
func (h *TaskHandler) Create(c *gin.Context) {
    userID, exists := c.Get("user_id")
    if !exists {
        c.JSON(http.StatusUnauthorized, gin.H{"error": "User not found in context"})
        return
    }

    var task domain.Task
    if err := c.ShouldBindJSON(&task); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    task.UserID = userID.(int64)

    if err := h.TaskUseCase.Create(c.Request.Context(), &task); err != nil {
        fmt.Printf("\n>>> DEBUG ERROR CREATE TASK: %v\n\n", err)
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    c.JSON(http.StatusCreated, task)
}

// FetchTasks godoc
// @Summary      Get All Tasks
// @Description  Mengambil semua tugas milik user yang sedang login
// @Tags         tasks
// @Produce      json
// @Security     BearerAuth
// @Success      200  {array}   domain.Task
// @Router       /tasks/ [get]
func (h *TaskHandler) Fetch(c *gin.Context) {
    userID, exists := c.Get("user_id")
    if !exists {
        c.JSON(http.StatusUnauthorized, gin.H{"error": "User not found in context"})
        return
    }

    tasks, err := h.TaskUseCase.Fetch(c.Request.Context(), userID.(int64))
    if err != nil {
        fmt.Printf("\n>>> DEBUG ERROR FETCH: %v\n\n", err)
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    c.JSON(http.StatusOK, tasks)
}

// Update status task
func (h *TaskHandler) UpdateStatus(c *gin.Context) {
    idParam := c.Param("id")
    id, err := strconv.ParseInt(idParam, 10, 64)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid task ID"})
        return
    }

    userID := c.MustGet("user_id").(int64)

    var req struct {
        Status string `json:"status"`
    }
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    if err := h.TaskUseCase.UpdateStatus(c.Request.Context(), id, userID, req.Status); err != nil {
        fmt.Printf("\n>>> DEBUG ERROR UPDATE: %v\n\n", err)
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    c.JSON(http.StatusOK, gin.H{"message": "Task updated"})
}

// Delete task
func (h *TaskHandler) Delete(c *gin.Context) {
    idParam := c.Param("id")
    id, err := strconv.ParseInt(idParam, 10, 64)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid task ID"})
        return
    }

    userID := c.MustGet("user_id").(int64)

    if err := h.TaskUseCase.Delete(c.Request.Context(), id, userID); err != nil {
        fmt.Printf("\n>>> DEBUG ERROR DELETE: %v\n\n", err)
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    c.JSON(http.StatusOK, gin.H{"message": "Task deleted"})
}

// --- SUBTASK HANDLERS ---

func (h *TaskHandler) AddSubtask(c *gin.Context) {
    idParam := c.Param("id") // Task ID
    taskID, _ := strconv.ParseInt(idParam, 10, 64)

    var req struct{ Title string `json:"title"` }
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    if err := h.TaskUseCase.AddSubtask(c.Request.Context(), taskID, req.Title); err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }
    c.JSON(http.StatusOK, gin.H{"message": "Subtask added"})
}

func (h *TaskHandler) ToggleSubtask(c *gin.Context) {
    idParam := c.Param("sub_id")
    subID, _ := strconv.ParseInt(idParam, 10, 64)

    if err := h.TaskUseCase.ToggleSubtask(c.Request.Context(), subID); err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }
    c.JSON(http.StatusOK, gin.H{"message": "Toggled"})
}

func (h *TaskHandler) DeleteSubtask(c *gin.Context) {
    idParam := c.Param("sub_id")
    subID, _ := strconv.ParseInt(idParam, 10, 64)

    if err := h.TaskUseCase.DeleteSubtask(c.Request.Context(), subID); err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }
    c.JSON(http.StatusOK, gin.H{"message": "Deleted"})
}
