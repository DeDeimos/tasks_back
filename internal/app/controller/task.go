package controller

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"

	"awesomeProject/internal/app/ds"
	"awesomeProject/internal/app/repository"

	"github.com/gin-gonic/gin"
)

func GetTaskByID(repository *repository.Repository, c *gin.Context) {
	var task *ds.Task

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusInternalServerError, err)
		return
	}

	if id < 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"Status":  "Failed",
			"Message": "неверное значение id",
		})
		return
	}

	task, err = repository.GetTaskByID(uint(id))
	if err != nil {
		c.JSON(http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusOK, task)
}

func GetTasksByRequestID(repository *repository.Repository, c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"Status":  "Failed",
			"Message": "неверное значение id",
		})
		return
	}

	task, err := repository.GetRequestByID(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusOK, task)
}

func GetAllTasks(repository *repository.Repository, c *gin.Context) {

	tasks, err := repository.GetAllTasks()
	if err != nil {
		c.JSON(http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusOK, tasks)
}

func DeleteTask(repository *repository.Repository, c *gin.Context) {

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusInternalServerError, err)
		return
	}

	if id < 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"Status":  "Failed",
			"Message": "неверное значение id",
		})
		return
	}

	err = repository.DeleteTask(uint(id))

	if err != nil {
		c.JSON(http.StatusBadRequest, err)
		return
	}
	c.JSON(http.StatusOK, "deleted successful")
}

func CreateTask(repository *repository.Repository, c *gin.Context) {
	var task ds.Task

	// Попробуйте извлечь JSON-данные из тела запроса и привести их к структуре Consultation
	if err := c.ShouldBindJSON(&task); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"Status":  "Failed",
			"Message": "неверные данные задания",
		})
		return
	}

	err := repository.CreateTask(task)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"task":   task,
		"status": "added",
	})
}

func UpdateTask(repository *repository.Repository, c *gin.Context) {
	// Извлекаем id консультации из параметра запроса
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusInternalServerError, err)
		return
	}

	// Проверяем, что id неотрицательный
	if id < 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"Status":  "Failed",
			"Message": "неверное значение id",
		})
		return
	}

	// Попробуем извлечь JSON-данные из тела запроса и привести их к структуре Consultation
	var updatedTask ds.Task
	if err := c.ShouldBindJSON(&updatedTask); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"Status":  "Failed",
			"Message": "неверные данные задания",
		})
		return
	}
	fmt.Println(updatedTask)
	// Обновляем консультацию в репозитории
	err = repository.UpdateTask(id, updatedTask)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status": "updated",
	})
}

func AddTaskToRequest(repository *repository.Repository, c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusInternalServerError, err)
		return
	}

	if id < 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"Status":  "Failed",
			"Message": "неверное значение id",
		})
		return
	}
	log.Println("id > 0")
	err = repository.AddTaskToRequest(id, 1)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"Status":  "Failed",
			"Message": "неверное значение id",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status": "added to request",
	})
}

func AddTaskImage(repository *repository.Repository, c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusInternalServerError, err)
		return
	}

	if id < 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"Status":  "Failed",
			"Message": "неверное значение id",
		})
		return
	}
	// Чтение изображения из запроса
	image, err := c.FormFile("image")
	log.Println(image, err)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid image"})
		return
	}

	// Чтение содержимого изображения в байтах
	file, err := image.Open()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка при открытии"})
		return
	}
	defer file.Close()

	imageBytes, err := io.ReadAll(file)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка чтения"})
		return
	}
	// Получение Content-Type из заголовков запроса
	contentType := image.Header.Get("Content-Type")

	// Вызов функции репозитория для добавления изображения
	err = repository.AddTaskImage(id, imageBytes, contentType)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка сервера"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Image uploaded successfully"})

}
