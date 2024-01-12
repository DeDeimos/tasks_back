package controller

import (
	"net/http"
	"strconv"

	"awesomeProject/internal/app/repository"

	"github.com/gin-gonic/gin"
)

// @Summary Delete Task From Request
// @Security ApiKeyAuth
// @Description delete task from request
// @Tags Task-Request
// @ID delete-task-from-request
// @Accept       json
// @Produce      json
// @Param        id_c   path      int  true  "ID задания"
// @Param        id_r   path      int  true  "ID заявки"
// @Success 200 {string} string "Консультация была удалена из заявки"
// @Failure 400 {string} string "Некорректный запрос"
// @Failure 404 {string} string "Некорректный запрос"
// @Failure 500 {string} string "Ошибка сервера"
// @Router /task-request/delete/task/{id_c}/request/{id_r} [delete]
func DeleteTaskRequest(repository *repository.Repository, c *gin.Context) {
	var idT, idR int
	var err error
	idT, err = strconv.Atoi(c.Param("id_c"))
	if err != nil {
		c.JSON(http.StatusInternalServerError, err)
		return
	}

	if idT < 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"Status":  "Failed",
			"Message": "неверное значение id занятия",
		})
		return
	}

	idR, err = strconv.Atoi(c.Param("id_r"))
	if err != nil {
		c.JSON(http.StatusInternalServerError, err)
		return
	}

	if idR < 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"Status":  "Failed",
			"Message": "неверное значение id заявки",
		})
		return
	}

	err = repository.DeleteTaskRequest(idT, idR)

	if err != nil {
		c.JSON(http.StatusBadRequest, err)
		return
	}
	c.JSON(http.StatusOK, "deleted successful")
}

// @Summary Change order of Task in Request
// @Security ApiKeyAuth
// @Description change order of task in request
// @Tags Task-Request
// @ID change-order-of-task-in-request
// @Accept       json
// @Produce      json
// @Param        id_c   path      int  true  "ID задания"
// @Param        id_r   path      int  true  "ID заявки"
// @Param        id_o   path      int  true  "порядок"
// @Success 200 {string} string "Порядок задания был изменен в заявке"
// @Failure 400 {string} string "Некорректный запрос"
// @Failure 404 {string} string "Некорректный запрос"
// @Failure 500 {string} string "Ошибка сервера"
// @Router /task-request/delete/task/{id_c}/request/{id_r}/order/{id_o} [put]
func ChangeOrder(repository *repository.Repository, c *gin.Context) {
	var idT, idR, idO int
	var err error
	idT, err = strconv.Atoi(c.Param("id_c"))
	if err != nil {
		c.JSON(http.StatusInternalServerError, err)
		return
	}

	if idT < 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"Status":  "Failed",
			"Message": "неверное значение id занятия",
		})
		return
	}

	idR, err = strconv.Atoi(c.Param("id_r"))
	if err != nil {
		c.JSON(http.StatusInternalServerError, err)
		return
	}

	if idR < 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"Status":  "Failed",
			"Message": "неверное значение id заявки",
		})
		return
	}

	idO, err = strconv.Atoi(c.Param("id_o"))
	if err != nil {
		c.JSON(http.StatusInternalServerError, err)
		return
	}

	if idO < 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"Status":  "Failed",
			"Message": "неверное значение id заявки",
		})
		return
	}

	err = repository.ChangeOrder(idT, idR, idO)

	if err != nil {
		c.JSON(http.StatusBadRequest, err)
		return
	}
	c.JSON(http.StatusOK, "change successful")
}
