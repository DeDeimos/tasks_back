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

// @Summary Get Task by ID
// @Description Show task by ID
// @Tags Tasks
// @ID get-task-by-id
// @Accept       json
// @Produce      json
// @Param        id   path      int  true  "ID задания"
// @Success 200 {object} ds.Task
// @Failure 400 {object} ds.Task "Некорректный запрос"
// @Failure 404 {object} ds.Task "Некорректный запрос"
// @Failure 500 {object} ds.Task "Ошибка сервера"
// @Router /tasks/{id} [get]
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

// @Summary Get Task by request ID
// @Security ApiKeyAuth
// @Description Show task by ID of request
// @Tags Tasks
// @ID get-task-by-id-of-request
// @Accept       json
// @Produce      json
// @Param        id   path      int  true  "ID заявки"
// @Success 200 {object} ds.Task
// @Failure 400 {object} ds.Task "Некорректный запрос"
// @Failure 404 {object} ds.Task "Некорректный запрос"
// @Failure 500 {object} ds.Task "Ошибка сервера"
// @Router /tasks/request/{id} [get]
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

// @Summary Get Tasks
// @Description Get all tasks
// @Tags Tasks
// @ID get-tasks
// @Produce json
// @Success 200 {object} ds.Task
// @Failure 400 {object} ds.Task "Некорректный запрос"
// @Failure 404 {object} ds.Task "Некорректный запрос"
// @Failure 500 {object} ds.Task "Ошибка сервера"
// @Router /tasks [get]
func GetAllTasks(repository *repository.Repository, c *gin.Context) {
	title := c.DefaultQuery("title", "")
	tasks, err := repository.GetAllTasks("active", title)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err)
		return
	}
	log.Println(1)
	userID, contextError := c.Value("userID").(uint)
	log.Println(userID)
	log.Println(contextError)
	if !contextError {
		c.JSON(http.StatusOK, gin.H{
			"ActiveRequestID": nil,
			"tasks":           tasks,
		})
		return
	}
	log.Println(2)
	log.Println(userID)
	fmt.Println(contextError)

	draftID := repository.GetDraftUser(int(userID))
	if draftID == -1 {
		c.JSON(http.StatusOK, gin.H{
			"ActiveRequestID": nil,
			"tasks":           tasks,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"ActiveRequestID": draftID,
		"tasks":           tasks,
	})
}

// @Summary Delete task by ID
// @Security ApiKeyAuth
// @Description Delete task by ID
// @Tags Tasks
// @ID delete-task-by-id
// @Accept       json
// @Produce      json
// @Param        id   path      int  true  "ID задания"
// @Success 200 {string} string
// @Failure 400 {object} ds.Task "Некорректный запрос"
// @Failure 404 {object} ds.Task "Некорректный запрос"
// @Failure 500 {object} ds.Task "Ошибка сервера"
// @Router /tasks/delete/{id} [delete]
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

// @Summary create task
// @Security ApiKeyAuth
// @Description create task
// @Tags Tasks
// @ID create-task
// @Accept json
// @Produce json
// @Param input body ds.Task true "task info"
// @Success 200 {string} string
// @Failure 400 {object} ds.Task "Некорректный запрос"
// @Failure 404 {object} ds.Task "Некорректный запрос"
// @Failure 500 {object} ds.Task "Ошибка сервера"
// @Router /tasks/create [post]
func CreateTask(repository *repository.Repository, c *gin.Context) {
	var task ds.Task

	// Попробуйте извлечь JSON-данные из тела запроса и привести их к структуре Task
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

// @Summary update task
// @Security ApiKeyAuth
// @Description update task
// @Tags Tasks
// @ID update-task
// @Accept       json
// @Produce      json
// @Param        id   path      int  true  "ID задания"
// @Param input body ds.Task true "task info"
// @Success 200 {string} string
// @Failure 400 {object} ds.Task "Некорректный запрос"
// @Failure 404 {object} ds.Task "Некорректный запрос"
// @Failure 500 {object} ds.Task "Ошибка сервера"
// @Router /tasks/update/{id} [put]
func UpdateTask(repository *repository.Repository, c *gin.Context) {
	// Извлекаем id задания из параметра запроса
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

	// Попробуем извлечь JSON-данные из тела запроса и привести их к структуре Task
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

// @Summary add task to request
// @Security ApiKeyAuth
// @Description add task to request
// @Tags Tasks
// @ID add-task-to-request
// @Accept       json
// @Produce      json
// @Param        id   path      int  true  "ID задания"
// @Success 200 {string} string
// @Failure 400 {object} ds.Task "Некорректный запрос"
// @Failure 404 {object} ds.Task "Некорректный запрос"
// @Failure 500 {object} ds.Task "Ошибка сервера"
// @Router /tasks/{id}/add-to-request [post]
func AddTaskToRequest(repository *repository.Repository, c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusInternalServerError, err)
		return
	}
	userID, contextError := c.Value("userID").(uint)
	if !contextError {
		c.JSON(http.StatusBadRequest, gin.H{
			"Status":  "Failed",
			"Message": "ошибка при авторизации",
		})
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
	err = repository.AddTaskToRequest(id, int(userID))
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

// @Summary Add task image
// @Security ApiKeyAuth
// @Description Add an image to a specific task by ID.
// @Tags Tasks
// @ID add-task-image
// @Accept       json
// @Produce      json
// @Param        id   path      int  true  "ID задания"
// @Param image formData file true "Image file to be uploaded"
// @Success 200 {string} string
// @Failure 400 {object} ds.Task "Некорректный запрос"
// @Failure 404 {object} ds.Task "Некорректный запрос"
// @Failure 500 {object} ds.Task "Ошибка сервера"
// @Router /tasks/{id}/addImage [post]
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
