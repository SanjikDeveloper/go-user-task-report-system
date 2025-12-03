package delivery

import (
	"net/http"
	"strconv"
	"strings"
	"tasks-service/internal/models"

	"github.com/gin-gonic/gin"
)

func (h *Handler) createTask(c *gin.Context) {
	h.logger.Info("createTask: request received")
	var task models.Task
	var err error
	if err = c.ShouldBindJSON(&task); err != nil {
		h.logger.Error("createTask: ShouldBindJSON failed: %v", err)
		return
	}

	if err = task.Validate(); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}

	userID, err := getUserId(c)
	if err != nil {
		h.logger.Error("createTask: getUserId failed: %v", err)
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}
	h.logger.Info("createTask: userId retrieved from context: %d", userID)

	task.UserID = userID
	id, err := h.services.TaskTodo.CreateTask(c.Request.Context(), task)
	if err != nil {
		h.logger.Error("createTask: CreateTask failed: %v", err)
		// Проверяем тип ошибки
		if strings.Contains(err.Error(), "does not exist") || strings.Contains(err.Error(), "foreign key constraint") {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "User not found in database. Please ensure user-service and tasks-service use the same database, or sign in again to get a valid token.",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create task"})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"id": id})
}

func (h *Handler) getTasks(c *gin.Context) {
	h.logger.Info("getTasks: request received")
	id := c.Param("id")

	taskId, err := strconv.Atoi(id)
	if err != nil {
		h.logger.Error("getTasks: strconv.Atoi failed: %v", err)
		return
	}

	userID, err := getUserId(c)
	if err != nil {
		h.logger.Error("getTasks: getUserId failed: %v", err)
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}
	h.logger.Info("getTasks: userId retrieved from context: %d", userID)

	task, err := h.services.GetTaskById(c.Request.Context(), taskId, userID)
	if err != nil {
		h.logger.Error("h.services.GetTaskById(...) err:", err)
		return
	}

	c.JSON(http.StatusOK, task)
}

func (h *Handler) updateTask(c *gin.Context) {
	h.logger.Info("updateTask: request received")
	var task models.Task

	if err := c.ShouldBindJSON(&task); err != nil {
		h.logger.Error("updateTask: ShouldBindJSON failed: %v", err)
		return
	}

	id := c.Param("id")
	taskID, err := strconv.Atoi(id)
	if err != nil {
		h.logger.Error("strconv.Atoi(...) err:", err)
		return
	}

	userID, err := getUserId(c)
	if err != nil {
		h.logger.Error("updateTask: getUserId failed: %v", err)
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}
	h.logger.Info("updateTask: userId retrieved from context: %d", userID)

	task.ID = taskID
	task.UserID = userID

	err = h.services.UpdateTask(c.Request.Context(), task)
	if err != nil {
		h.logger.Error("h.services.UpdateTask(...) err:", err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Task updated"})
}

func (h *Handler) deleteTask(c *gin.Context) {
	h.logger.Info("deleteTask: request received")
	id := c.Param("id")
	taskId, err := strconv.Atoi(id)
	if err != nil {
		// TODO: ошибки не надо отправлять на клиента, переделать ответы в ошибках
		h.logger.Error("deleteTask: strconv.Atoi failed: %v", err)
		c.Status(http.StatusInternalServerError)
	}

	userID, err := getUserId(c)
	if err != nil {
		h.logger.Error("deleteTask: getUserId failed: %v", err)
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}
	h.logger.Info("deleteTask: userId retrieved from context: %d", userID)

	err = h.services.DeleteTask(c.Request.Context(), taskId, userID)
	if err != nil {
		h.logger.Error("h.services.DeleteTask(...) err:", err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Task deleted"})
}
