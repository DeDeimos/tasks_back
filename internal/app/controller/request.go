package controller

import (
	"log"
	"net/http"
	"strconv"
	"time"

	"awesomeProject/internal/app/ds"
	"awesomeProject/internal/app/repository"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func GetAllRequests(repository *repository.Repository, c *gin.Context) {

	userID := 3
	status := c.DefaultQuery("status", "")
	dateFrom := c.DefaultQuery("startDate", "")
	dateTo := c.DefaultQuery("endDate", "")
	const timeFormat1 = "2006-01-02T15:04:05.999999Z"
	const timeFormat = "2006-01-02 15:04:05"
	log.Println(dateFrom)
	log.Println(dateTo)
	user, err := repository.FindByID(userID)
	if err == gorm.ErrRecordNotFound {
		c.JSON(http.StatusBadRequest, gin.H{
			"Status":  "Failed",
			"Message": "неверное значение id",
		})
		return
	}

	log.Println(user)

	if dateFrom == "" && dateTo == "" {
		if user.Role == ds.USER_ROLE_MODERATOR {
			requests, err := repository.FindAllByModeratorID(userID, status, nil, nil)
			if err != nil {
				c.JSON(http.StatusInternalServerError, err)
				return
			}
			c.JSON(http.StatusOK, requests)
			return
		} else {
			requests, err := repository.FindAllByUserID(userID, status, nil, nil)
			if err != nil {
				c.JSON(http.StatusInternalServerError, err)
				return
			}
			c.JSON(http.StatusOK, requests)
			return
		}
	}

	timeFrom, err := time.Parse(timeFormat, dateFrom)
	if err != nil {
		timeFrom = time.Unix(0, 0)
	}

	log.Println(timeFrom)
	timeTo, err := time.Parse(timeFormat, dateTo)
	if err != nil {
		timeTo = time.Now()
	}
	log.Println(timeTo)
	if user.Role == ds.USER_ROLE_MODERATOR {
		requests, err := repository.FindAllByModeratorID(userID, status, &timeFrom, &timeTo)
		if err != nil {
			c.JSON(http.StatusInternalServerError, err)
			return
		}
		c.JSON(http.StatusOK, requests)
		return
	} else {
		requests, err := repository.FindAllByUserID(userID, status, &timeFrom, &timeTo)
		if err != nil {
			c.JSON(http.StatusInternalServerError, err)
			return
		}
		c.JSON(http.StatusOK, requests)
	}
	// if status != "" {
	// 	requests, err = repository.GetRequestsByStatus(status)
	// 	if err != nil {
	// 		c.JSON(http.StatusInternalServerError, err)
	// 		return
	// 	}

	// 	c.JSON(http.StatusOK, requests)
	// 	return
	// }
	// log.Println(startFormationDateStr + "ASSDA")
	// if startFormationDateStr != "" {
	// 	var startFormationDate time.Time
	// 	var endFormationDate time.Time
	// 	layout := "2006-01-02 15:04:05.000000"
	// 	startFormationDate, err = time.Parse(layout, startFormationDateStr)
	// 	log.Println(startFormationDate)
	// 	if err != nil {
	// 		c.JSON(http.StatusInternalServerError, err)
	// 		return
	// 	}
	// 	if endFormationDateStr != "" {
	// 		endFormationDate, err = time.Parse(layout, endFormationDateStr)

	// 		if err != nil {
	// 			c.JSON(http.StatusInternalServerError, err)
	// 			return
	// 		}
	// 	}
	// 	log.Panicln("here is not problem")
	// 	requests, err = repository.GetRequestsByDate(startFormationDate, endFormationDate)
	// 	if err != nil {
	// 		c.JSON(http.StatusInternalServerError, err)
	// 		return
	// 	}

	// 	c.JSON(http.StatusOK, requests)
	// 	return
	// }
	// log.Println("go here")

	// user_id := 1

	// requests, err = repository.GetAllRequests(uint(user_id))
	// if err != nil {
	// 	c.JSON(http.StatusInternalServerError, err)
	// 	return
	// }

	// c.JSON(http.StatusOK, requests)
}

func GetRequestByID(repository *repository.Repository, c *gin.Context) {
	var request *ds.Request

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

	request, err = repository.GetRequestByID(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusOK, request)
}

func DeleteRequest(repository *repository.Repository, c *gin.Context) {

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

	err = repository.DeleteRequest(id)

	if err != nil {
		c.JSON(http.StatusBadRequest, err)
		return
	}
	c.JSON(http.StatusOK, "deleted successful")
}

func UpdateRequest(repository *repository.Repository, c *gin.Context) {
	// Извлекаем id request из параметра запроса
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

	// Попробуем извлечь JSON-данные из тела запроса
	var updatedRequest ds.Request
	if err := c.ShouldBindJSON(&updatedRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"Status":  "Failed",
			"Message": "неверные данные консультации",
		})
		return
	}

	err = repository.UpdateRequest(id, updatedRequest)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status": "updated",
	})
}

func UpdateUserRequestStatus(repository *repository.Repository, c *gin.Context) {
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

	request, err := repository.GetRequestByID(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err)
		return
	}

	// Попробуем извлечь JSON-данные из тела запроса - новый статус
	var status ds.Status
	if err := c.ShouldBindJSON(&status); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"Status":  "Failed",
			"Message": "неверные данные статуса занятия",
		})
		return
	}

	log.Println(status)
	if (status.Status == "on_check" || status.Status == "deleted") && request.Status == "draft" {

	} else {
		c.JSON(http.StatusBadRequest, gin.H{
			"Status":  "Not Access",
			"Message": "неверное значение status для user",
		})
		return
	}

	err = repository.UpdateUserRequestStatus(id, status.Status)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status": "updated",
	})
}

func UpdateAdminRequestStatus(repository *repository.Repository, c *gin.Context) {
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

	// if role == "admin" && status.Status != "completed" && status.Status != "rejected" {
	// 	c.JSON(http.StatusBadRequest, gin.H{
	// 		"Status":  "Not Access",
	// 		"Message": "неверное значение status для admin",
	// 	})
	// 	return
	// }

	// if role == "user" && status.Status != "on_check" && status.Status != "deleted" {
	// 	c.JSON(http.StatusBadRequest, gin.H{
	// 		"Status":  "Not Access",
	// 		"Message": "неверное значение status для user",
	// 	})
	// 	return
	// }

	// err = repository.UpdateRequestStatus(id, status.Status)
	// if err != nil {
	// 	c.JSON(http.StatusInternalServerError, err)
	// 	return
	// }

	// c.JSON(http.StatusOK, gin.H{
	// 	"status": "updated",
	// })

	request, err := repository.GetRequestByID(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err)
		return
	}

	// Попробуем извлечь JSON-данные из тела запроса - новый статус
	var status ds.Status
	if err := c.ShouldBindJSON(&status); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"Status":  "Failed",
			"Message": "неверные данные статуса занятия",
		})
		return
	}

	log.Println(status)
	if (status.Status == "completed" || status.Status == "rejected") && request.Status == "on_check" {

	} else {
		c.JSON(http.StatusBadRequest, gin.H{
			"Status":  "Not Access",
			"Message": "неверное значение status для admin",
		})
		return
	}

	err = repository.UpdateAdminRequestStatus(id, status.Status)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status": "updated",
	})
}
