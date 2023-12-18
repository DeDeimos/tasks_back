package app

import (
	"log"
	"net/http"
	"strconv"
	"strings"

	"awesomeProject/docs"
	"awesomeProject/internal/app/controller"
	"awesomeProject/internal/app/ds"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	_ "github.com/lib/pq"
)

func (a *Application) StartServer() {
	log.Println("Server start up")

	r := gin.Default()
	// c := controller.NewController(a.repository)
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	docs.SwaggerInfo.Title = "RIpPeakBack"
	docs.SwaggerInfo.Description = "rip course project about alpinists and their expeditions"
	docs.SwaggerInfo.Version = "1.0"
	docs.SwaggerInfo.Host = "localhost:8080"
	docs.SwaggerInfo.BasePath = ""

	r.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "https://dedeimos.github.io, http://localhost:1420")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(http.StatusOK)
			return
		}

		c.Next()
	})

	r.Static("/styles", "./resources/styles")
	r.Static("/js", "./resources/js")
	r.Static("/img", "./resources/img")
	r.Static("/hacker", "./resources")
	r.LoadHTMLGlob("templates/*")

	r.GET("/", func(c *gin.Context) {
		var tasks []ds.Task
		tasks, err := a.repository.GetAllTasks("active", "")
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
	AuthorisationGroup := r.Group("/auth")
	{
		AuthorisationGroup.POST("/registration", a.Register)
		AuthorisationGroup.POST("/login", a.Login)
		AuthorisationGroup.POST("/logout", a.Logout)
	}
	// r.POST("/login", a.Login) // там где мы ранее уже заводили эндпоинты
	// r.POST("/sign_up", a.Register)
	// r.POST("/log_out", a.Logout)
	// r.Use(a.WithAuthCheck("user", "admin")).GET("/ping", a.Ping)

	TaskGroup := r.Group("/tasks")
	{
		TaskGroup.GET("/",
			func(c *gin.Context) {
				controller.GetAllTasks(a.repository, c)
			})
		TaskGroup.GET("/:id", func(c *gin.Context) {
			controller.GetTaskByID(a.repository, c)
		})
		TaskGroup.Use(a.WithAuthCheck("admin")).DELETE("/delete/:id", func(c *gin.Context) {
			controller.DeleteTask(a.repository, c)
		})
		TaskGroup.Use(a.WithAuthCheck("admin")).PUT("/update/:id", func(c *gin.Context) {
			controller.UpdateTask(a.repository, c)
		})
		TaskGroup.Use(a.WithAuthCheck("admin")).POST("/create", func(c *gin.Context) {
			controller.CreateTask(a.repository, c)
		})
		TaskGroup.Use(a.WithAuthCheck("user")).POST("/:id/add-to-request", func(c *gin.Context) {
			controller.AddTaskToRequest(a.repository, c)
		})
		TaskGroup.Use(a.WithAuthCheck("admin")).PUT("/tasks/:id/add-image", func(c *gin.Context) {
			controller.AddTaskImage(a.repository, c)
		})

	}

	// 	r.GET("/tasks", func(c *gin.Context) {
	// 	controller.GetAllTasks(a.repository, c)
	// })
	// r.GET("/tasks/:id", func(c *gin.Context) {
	// 	controller.GetTaskByID(a.repository, c)
	// })
	// r.DELETE("/tasks/delete/:id", func(c *gin.Context) {
	// 	controller.DeleteTask(a.repository, c)
	// })
	// r.POST("/tasks/create", func(c *gin.Context) {
	// 	controller.CreateTask(a.repository, c)
	// })
	// r.PUT("/tasks/update/:id", func(c *gin.Context) {
	// 	controller.UpdateTask(a.repository, c)
	// })

	// r.PUT("/tasks/:id/add-image", func(c *gin.Context) {
	// 	controller.AddTaskImage(a.repository, c)
	// })

	// r.POST("/tasks/:id/add-to-request", func(c *gin.Context) {
	// 	controller.AddTaskToRequest(a.repository, c)
	// })

	RequestGroup := r.Group("/requests")
	{
		RequestGroup.Use(a.WithAuthCheck("admin", "user")).GET("/", func(c *gin.Context) {
			controller.GetAllRequests(a.repository, c)
		})
		RequestGroup.Use(a.WithAuthCheck("admin")).GET("/:id", func(c *gin.Context) {
			controller.GetTasksByRequestID(a.repository, c)
		})
		RequestGroup.Use(a.WithAuthCheck("admin")).DELETE("/delete/:id", func(c *gin.Context) {
			controller.DeleteRequest(a.repository, c)
		})
		RequestGroup.Use(a.WithAuthCheck("admin")).PUT("/admin/:id/update-status", func(c *gin.Context) {

			controller.UpdateAdminRequestStatus(a.repository, c)
		})
		RequestGroup.Use(a.WithAuthCheck("user")).PUT("/user/:id/update-status", func(c *gin.Context) {
			controller.UpdateUserRequestStatus(a.repository, c)
		})

	}

	// r.GET("/requests", func(c *gin.Context) {
	// 	controller.GetAllRequests(a.repository, c)
	// })
	// r.GET("/requests/:id", func(c *gin.Context) {
	// 	controller.GetTasksByRequestID(a.repository, c)
	// })
	// r.DELETE("/requests/delete/:id", func(c *gin.Context) {
	// 	controller.DeleteRequest(a.repository, c)
	// })
	// r.PUT("/requests/update/:id", func(c *gin.Context) {
	// 	controller.UpdateRequest(a.repository, c)
	// })
	// r.PUT("/requests/admin/:id/update-status", func(c *gin.Context) {

	// 	controller.UpdateAdminRequestStatus(a.repository, c)
	// })

	// r.PUT("/requests/user/:id/update-status", func(c *gin.Context) {
	// 	controller.UpdateUserRequestStatus(a.repository, c)

	// })
	TaskRequestGroup := r.Group("/task-request")
	{
		TaskRequestGroup.Use(a.WithAuthCheck("user")).DELETE("/delete/task/:id_c/request/:id_r", func(c *gin.Context) {
			controller.DeleteTaskRequest(a.repository, c)
		})
	}
	// r.DELETE("/task-request/delete/task/:id_c/request/:id_r", func(c *gin.Context) {
	// 	controller.DeleteTaskRequest(a.repository, c)
	// })

	r.Run()

	log.Println("Server down")
}

type pingReq struct{}
type pingResp struct {
	Status string `json:"status"`
}

// Ping godoc
// @Summary      Show hello text
// @Description  very very friendly response
// @Tags         Tests
// @Produce      json
// @Success      200  {object}  pingResp
// @Router       /ping/{name} [get]
func (a *Application) Ping(gCtx *gin.Context) {
	name := gCtx.Param("name")
	gCtx.String(http.StatusOK, "Hello %s", name)
}
