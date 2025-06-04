package http

import (
	"net/http"

	portIn "hex-todo/internal/core/ports/in"

	"github.com/gin-gonic/gin"
)

type TaskHandler struct {
	Service portIn.TaskService
}

func NewRouter(service portIn.TaskService) *gin.Engine {
	handler := &TaskHandler{Service: service}
	router := gin.Default()

	router.POST("/tasks", handler.CreateTask)
	router.GET("/tasks", handler.ListTasks)

	return router
}

type createTaskRequest struct {
	Title string `json:"title"`
}

func (h *TaskHandler) CreateTask(c *gin.Context) {
	var req createTaskRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "JSON inválido"})
		return
	}

	task, err := h.Service.CreateTask(req.Title)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, task)
}

func (h *TaskHandler) ListTasks(c *gin.Context) {
	tasks, err := h.Service.ListTasks()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "No se pudieron obtener las tareas"})
		return
	}

	c.JSON(http.StatusOK, tasks)
}
