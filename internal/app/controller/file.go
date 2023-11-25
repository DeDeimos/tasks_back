package controller

import (
	minioclient "awesomeProject/internal/minioClient"

	"github.com/gin-gonic/gin"
)

type FileHandler struct{}

func NewFileHandler() *FileHandler {
	return &FileHandler{}
}

func (h *FileHandler) FindFile(c *gin.Context) {
	bucket := c.Param("bucket")
	filename := c.Param("filename")

	contentBytes, contentType, err := minioclient.ReadObject(bucket, filename)
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	c.Data(200, contentType, contentBytes)
}
