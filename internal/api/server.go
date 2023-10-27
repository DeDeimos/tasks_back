package app

import (
	"log"
	"net/http"
	"strconv"
	"strings"

	"awesomeProject/internal/app/controller"
	"awesomeProject/internal/app/ds"

	"github.com/gin-gonic/gin"

	_ "github.com/lib/pq"
)

func (a *Application) StartServer() {
	log.Println("Server start up")

	r := gin.Default()

	r.Static("/styles", "./resources/styles")
	r.Static("/js", "./resources/js")
	r.Static("/img", "./resources/img")
	r.Static("/hacker", "./resources")
	r.LoadHTMLGlob("templates/*")

	r.GET("/", func(c *gin.Context) {
		var tasks []ds.Task
		tasks, err := a.repository.GetAllTasks()
		if err != nil { // если не получилось
			log.Printf("cant get product by id %v", err)
			return
		}
		searchQuery := c.DefaultQuery("fsearch", "")

		if searchQuery == "" {
			c.HTML(http.StatusOK, "index.tmpl", gin.H{
				"services": tasks,
			})
			return
		}

		var result []ds.Task

		for _, task := range tasks {
			if strings.Contains(strings.ToLower(task.Name), strings.ToLower(searchQuery)) {
				result = append(result, task)
			}
		}

		c.HTML(http.StatusOK, "index.tmpl", gin.H{
			"services":    result,
			"search_text": searchQuery,
		})
	})

	r.POST("/delete/:id", func(c *gin.Context) {
		id, err := strconv.Atoi(c.Param("id"))

		if err != nil {
			// Обработка ошибки
			log.Printf("cant get task by id %v", err)
			c.Redirect(http.StatusMovedPermanently, "/")
		}
		a.repository.DeleteTask(uint(id))
		c.Redirect(http.StatusMovedPermanently, "/")
	})

	r.GET("/service/:id", func(c *gin.Context) {
		var task *ds.Task

		id, err := strconv.Atoi(c.Param("id"))
		task, err = a.repository.GetTaskByID(uint(id))
		if err != nil {
			// Обработка ошибки
			log.Printf("cant get service by id %v", err)
			return
		}

		c.HTML(http.StatusOK, "card.tmpl", task)
	})

	//------------------------------------------------------------------------------
	r.GET("/tasks", func(c *gin.Context) {
		controller.GetAllTasks(a.repository, c)
	})
	r.GET("/tasks/:id", func(c *gin.Context) {
		controller.GetTaskByID(a.repository, c)
	})
	r.DELETE("/tasks/delete/:id", func(c *gin.Context) {
		controller.DeleteTask(a.repository, c)
	})
	r.POST("/tasks/create", func(c *gin.Context) {
		controller.CreateTask(a.repository, c)
	})
	r.PUT("/tasks/update/:id", func(c *gin.Context) {
		controller.UpdateTask(a.repository, c)
	})

	r.PUT("/tasks/:id/add-image", func(c *gin.Context) {
		controller.AddTaskImage(a.repository, c)
	})

	r.POST("/tasks/:id/add-to-request", func(c *gin.Context) {
		controller.AddTaskToRequest(a.repository, c)
	})

	r.GET("/requests", func(c *gin.Context) {
		controller.GetAllRequests(a.repository, c)
	})
	r.GET("/requests/:id", func(c *gin.Context) {
		controller.GetTasksByRequestID(a.repository, c)
	})
	r.DELETE("/requests/delete/:id", func(c *gin.Context) {
		controller.DeleteRequest(a.repository, c)
	})
	r.PUT("/requests/update/:id", func(c *gin.Context) {
		controller.UpdateRequest(a.repository, c)
	})
	r.PUT("/requests/:role/:id/update-status", func(c *gin.Context) {
		controller.UpdateRequestStatus(a.repository, c)
	})

	r.DELETE("/task-request/delete/task/:id_c/request/:id_r", func(c *gin.Context) {
		controller.DeleteTaskRequest(a.repository, c)
	})

	r.Run()

	log.Println("Server down")
}
