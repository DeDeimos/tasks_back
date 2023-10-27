package controller

import (
	"net/http"
	"strconv"

	"awesomeProject/internal/app/repository"

	"github.com/gin-gonic/gin"
)

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
