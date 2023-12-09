package controller

import (
	"awesomeProject/internal/app/repository"
)

type Controller struct {
	Repo *repository.Repository
}

func NewController(repo *repository.Repository) *Controller {
	return &Controller{Repo: repo}
}
