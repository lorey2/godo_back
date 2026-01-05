package handlers

import (
	"net/http"
	"github.com/gin-gonic/gin"
	"goback/database"
	"goback/models"
)

func GetTasks(c *gin.Context) {
	userID, _ := c.Get("userID")
	var tasks []models.Task
	database.DB.Where("user_id = ?", userID).Order("position asc").Find(&tasks)
	c.JSON(http.StatusOK, tasks)
}

func AddTask(c *gin.Context) {
	userID, _ := c.Get("userID")
	var newTask models.Task
	if err := c.BindJSON(&newTask); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	newTask.Done = false
	newTask.UserID = uint(userID.(float64))
	if newTask.Category == "" { newTask.Category = "General" }
	if newTask.Priority == 0 { newTask.Priority = 2 }

	database.DB.Create(&newTask)
	c.JSON(http.StatusCreated, newTask)
}

func UpdateTask(c *gin.Context) {
	userID, _ := c.Get("userID")
	id := c.Param("id")
	var task models.Task
	if err := database.DB.Where("id = ? AND user_id = ?", id, userID).First(&task).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Task not found"})
		return
	}
	if err := c.BindJSON(&task); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	database.DB.Save(&task)
	c.JSON(http.StatusOK, task)
}

func DeleteTask(c *gin.Context) {
	userID, _ := c.Get("userID")
	id := c.Param("id")
	result := database.DB.Where("id = ? AND user_id = ?", id, userID).Delete(&models.Task{})
	if result.RowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Task not found"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Deleted"})
}

func ReorderTasks(c *gin.Context) {
	userID, _ := c.Get("userID")
	type OrderUpdate struct {
		ID       uint `json:"id"`
		Position int  `json:"position"`
	}
	var updates []OrderUpdate
	if err := c.BindJSON(&updates); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	for _, update := range updates {
		database.DB.Model(&models.Task{}).Where("id = ? AND user_id = ?", update.ID, userID).Update("position", update.Position)
	}
	c.Status(http.StatusOK)
}
