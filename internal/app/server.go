package app

import (
	"awesomeProject/internal/app/ds"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

func findTask(tasks []ds.Task, id uint) (*ds.Task, error) {
	for _, task := range tasks {
		if task.Task_id == id {
			return &task, nil
		}
	}
	return nil, fmt.Errorf("Задание с ID %d не найдено", id)
}

func (a *Application) StartServer() {

	// tasks := [5]ds.Task{
	// 	{1, "Математика", "Задание 1", "Решение интегралов", "/resources/math.jpg", "Вычислите определенный интеграл функции f(x) = x^4 - 3x^3 + 2x^2 - 5x + 8 в пределах от 0 до 5. Затем найдите корни уравнения: 2y^2 - 7y + 6 = 0. Для каждого корня вычислите значение функции g(x) = sin(x) + cos(x) и умножьте его на 10. Найдите среднее арифметическое полученных значений. В конечном итоге, представьте ответ в виде массива чисел."},
	// 	{2, "География", "Задание 2", "", "/resources/geography.jpg", "Изучите маршрут следующей путешественницы: она начала свой путь в городе A, затем двигалась на север на 100 километров и прибыла в город B. В B она изменила направление и двигалась на восток на 50 километров до города C. Затем она отправилась на юг на 75 километров и добралась до города D. Наконец, она двинулась на запад на 30 километров и достигла своего конечного пункта E. Вычислите общее расстояние, которое она прошла, и определите её текущее местоположение (город E) в координатах широты и долготы. Затем найдите расстояние между городами B и D и угол между направлениями движения в городах B и C в градусах."},
	// 	{3, "История", "Задание 3", "Первая мировая война", "/resources/history.jpg", "Исследуйте исторический период Первой мировой войны (1914-1918). Опишите причины начала войны, ключевые события, включая битвы и дипломатические переговоры, а также итоги и последствия этой войны для мировой истории. Уделите особое внимание роли великих держав, таких как Германия, Россия, Франция и Великобритания, в ходе конфликта. Назовите выдающихся лидеров, политиков и генералов, чьи действия оказали влияние на исход войны. Включите в ваш отчет даты ключевых событий и географические области, где разворачивались боевые действия."},
	// 	{4, "Физика", "Задание 4", "", "/resources/physic.jpg", "Изучите движение тела, брошенного вертикально вверх с начальной скоростью 20 м/с. Рассмотрите его движение в зависимости от времени. Вычислите момент времени, когда тело достигнет максимальной высоты и определите это значение высоты. Затем рассмотрите падение тела обратно на землю и найдите время, через которое оно упадет на землю, а также скорость удара о землю. Объясните, как воздействие силы тяжести и начальной скорости влияет на движение тела. Используйте законы Ньютона и уравнения движения для решения задачи. Укажите все известные параметры и формулы, которые использовались в расчетах."},
	// 	{5, "Музыка", "Задание 5", "Развитие рок-музыки", "/resources/music.jpg", "Исследуйте развитие жанра рок-музыки в 20-м и 21-м веках. Опишите ключевые моменты в истории рока, начиная с его зарождения в 1950-х годах. Укажите на влияние различных поджанров, таких как классический рок, психоделия, панк-рок, их представителей и особенности звучания. Рассмотрите также роль технологических инноваций, таких как электронные инструменты и интернет, в развитии и распространении рока. Уделите внимание выдающимся музыкантам, группам и альбомам, которые сделали значительный вклад в этот жанр. Укажите ключевые тренды и изменения в музыкальной индустрии, связанные с рок-музыкой, и их воздействие на музыкальную культуру."},
	// }

	log.Println("Server start up")

	r := gin.Default()

	// r.GET("/ping", func(c *gin.Context) {
	// 	c.JSON(http.StatusOK, gin.H{
	// 		"message": "pong",
	// 	})
	// })

	r.LoadHTMLGlob("templates/*")

	// r.GET("/tasks", func(c *gin.Context) {
	// 	c.HTML(http.StatusOK, "tasks.tmpl", gin.H{
	// 		"tasks": tasks,
	// 	})
	// })

	r.GET("/task/:id", func(c *gin.Context) {
		tasks, err := a.repository.GetActiveTasks()
		// Получите ID задания из URL
		taskID := c.Param("id")

		ID, errNum := strconv.Atoi(taskID)
		if errNum != nil {
			fmt.Println("Ошибка при преобразовании строки в int:", errNum)
			return
		}

		task, err := findTask(*tasks, uint(ID))
		if err != nil {
			fmt.Println("Ошибка:", err)
			return
		}
		// Здесь вы можете использовать ID, чтобы найти соответствующее задание в массиве tasks или базе данных
		// Загрузите данные о задании (название, цену и т.д.)

		// Отобразите страницу задания с подробностями
		c.HTML(http.StatusOK, "task.tmpl", gin.H{
			"task": task, // Здесь предполагается, что у вас есть переменная task с данными о задании
		})
	})

	r.GET("/tasks", func(c *gin.Context) {

		tasks, err := a.repository.GetActiveTasks()
		if err != nil {
			log.Println("Error with running\nServer")
			return
		}
		log.Println(tasks)

		searchQuery := c.DefaultQuery("q", "")
		var foundTasks []ds.Task
		for _, task := range *tasks {
			if strings.HasPrefix(strings.ToLower(task.Name), strings.ToLower(searchQuery)) {
				foundTasks = append(foundTasks, task)
			}
		}
		data := gin.H{
			"tasks":  foundTasks,
			"search": searchQuery,
		}
		c.HTML(http.StatusOK, "tasks.tmpl", data)
	})

	r.POST("/task/:id/delete", func(context *gin.Context) {
		log.Println("Hello there")
		// tasks, err := a.repository.GetActiveTasks()

		// if err != nil {
		// 	log.Println("Error with running\nServer down")
		// 	return
		// }

		id, err := strconv.Atoi(context.Param("id"))

		// log.Println(id)
		log.Println("General Kenobi")
		if err != nil {
			log.Println("err != nil")
			context.AbortWithStatus(404)
			return
		}

		// if id < 0 {
		// 	log.Println("id < 0")
		// 	context.AbortWithStatus(404)
		// 	return
		// }

		// if len(*tasks) == 0 {
		// 	log.Println("len(*tasks)")
		// 	context.AbortWithStatus(404)
		// 	return
		// }
		log.Println("No mistakes")

		// var activeTasks []ds.Task

		// for _, task := range *tasks {
		// 	if task.Task_id != uint(id) {
		// 		activeTasks = append(activeTasks, task)
		// 	} else {
		// 		taskToDelete = task
		// 	}
		// }

		// var db *sql.DB
		// _ = godotenv.Load()
		// db, err = sql.Open("postgres", dsn.FromEnv())
		// if err != nil {
		// 	log.Fatal(err)
		// }
		// defer db.Close()

		// _, err = a.repository. db.Exec("UPDATE tasks SET status = $1 WHERE id = $2", "del", taskToDelete.Task_id)
		// if err != nil {
		// 	context.AbortWithStatus(500)
		// 	return
		// }

		a.repository.DeleteTask(id)

		context.Redirect(http.StatusMovedPermanently, "/tasks")

		// context.HTML(http.StatusOK, "tasks.tmpl", gin.H{
		// 	"tasks": activeTasks,
		// })

		// context.HTML(http.StatusOK, "tasks.tmpl", gin.H{
		// 	"tasks": activeTasks,
		// })
	})

	r.Static("/resources", "./resources")

	r.Run()

	log.Println("Server down")
}
