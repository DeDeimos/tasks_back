package controller

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"awesomeProject/internal/app/ds"
	"awesomeProject/internal/app/repository"

	"github.com/gin-gonic/gin"
)

// @Summary Get Requests
// @Security ApiKeyAuth
// @Description Get all requests
// @Tags Requests
// @ID get-requests
// @Produce json
// @Success 200 {object} ds.Request
// @Failure 400 {object} ds.Request "Некорректный запрос"
// @Failure 404 {object} ds.Request "Некорректный запрос"
// @Failure 500 {object} ds.Request "Ошибка сервера"
// @Router /requests [get]
func GetAllRequests(repository *repository.Repository, c *gin.Context) {

	userID, contextError := c.Value("userID").(uint)
	if !contextError {
		c.JSON(http.StatusBadRequest, gin.H{
			"Status":  "Failed",
			"Message": "ошибка при авторизации",
		})
		return
	}

	log.Println(userID)
	fmt.Println(contextError)

	userRole, contextError := c.Value("userRole").(string)
	log.Println(userRole)
	fmt.Println(contextError)
	if !contextError {
		c.JSON(http.StatusBadRequest, gin.H{
			"Status":  "Failed",
			"Message": "ошибка при авторизации",
		})
		return
	}

	status := c.DefaultQuery("status", "")
	dateFrom := c.DefaultQuery("startDate", "")
	dateTo := c.DefaultQuery("endDate", "")
	const timeFormat = "2006-01-02 15:04:05"

	// user, err := repository.FindByID(userID)
	// if err == gorm.ErrRecordNotFound {
	// 	c.JSON(http.StatusBadRequest, gin.H{
	// 		"Status":  "Failed",
	// 		"Message": "неверное значение id",
	// 	})
	// 	return
	// }

	// log.Println(user)

	if dateFrom == "" && dateTo == "" {
		if userRole == ds.USER_ROLE_MODERATOR {
			requests, err := repository.FindAllByModeratorID(userID, status, nil, nil)
			if err != nil {
				c.JSON(http.StatusInternalServerError, err)
				return
			}
			c.JSON(http.StatusOK, requests)
			return
		}
		requests, err := repository.FindAllByUserID(userID, status, nil, nil)
		if err != nil {
			c.JSON(http.StatusInternalServerError, err)
			return
		}
		c.JSON(http.StatusOK, requests)
		return

	}

	timeFrom, err := time.Parse(timeFormat, dateFrom)
	if err != nil {
		timeFrom = time.Unix(0, 0)
	}
	timeTo, err := time.Parse(timeFormat, dateTo)
	if err != nil {
		timeTo = time.Now()
	}
	if userRole == ds.USER_ROLE_MODERATOR {
		requests, err := repository.FindAllByModeratorID(userID, status, &timeFrom, &timeTo)
		if err != nil {
			c.JSON(http.StatusInternalServerError, err)
			return
		}
		c.JSON(http.StatusOK, requests)
		return
	}
	requests, err := repository.FindAllByUserID(userID, status, &timeFrom, &timeTo)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err)
		return
	}
	c.JSON(http.StatusOK, requests)

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

// @Summary Delete Request by ID
// @Security ApiKeyAuth
// @Description Delete request by ID
// @Tags Requests
// @ID delete-request-by-id
// @Accept       json
// @Produce      json
// @Param        id   path      int  true  "ID заявки"
// @Success 200 {string} string
// @Failure 400 {object} ds.Request "Некорректный запрос"
// @Failure 404 {object} ds.Request "Некорректный запрос"
// @Failure 500 {object} ds.Request "Ошибка сервера"
// @Router /requests/delete/{id} [delete]
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

// @Summary Update Request Status By User
// @Security ApiKeyAuth
// @Description Update request status by user
// @Tags Requests
// @ID update-request-status-by-user
// @Accept       json
// @Produce      json
// @Param        id   path      int  true  "ID заявки"
// @Success 200 {string} string
// @Failure 400 {object} ds.Request "Некорректный запрос"
// @Failure 404 {object} ds.Request "Некорректный запрос"
// @Failure 500 {object} ds.Request "Ошибка сервера"
// @Router /requests/{id}/user/update-status [put]
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

// @Summary Update Request Status By Moderator
// @Security ApiKeyAuth
// @Description Update request by moderator
// @Tags Requests
// @ID update-request-status-by-moderator
// @Accept       json
// @Produce      json
// @Param        id   path      int  true  "ID заявки"
// @Param input body ds.Status true "status info"
// @Success 200 {string} string
// @Failure 400 {object} ds.Request "Некорректный запрос"
// @Failure 404 {object} ds.Request "Некорректный запрос"
// @Failure 500 {object} ds.Request "Ошибка сервера"
// @Router /requests/{id}/admin/update-status [put]
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

	userID, contextError := c.Value("userID").(uint)
	if !contextError {
		c.JSON(http.StatusBadRequest, gin.H{
			"Status":  "Failed",
			"Message": "ошибка при авторизации",
		})
		return
	}

	log.Println(userID)

	err = repository.UpdateAdminRequestStatus(id, status.Status, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status": "updated",
	})
}
