package repository

import (
	"awesomeProject/internal/app/ds"
)

func (r *Repository) DeleteTaskRequest(idT int, idR int) error {

	request := &ds.Request{}

	error := r.db.Preload("Tasks").First(request, "request_id = ?", idR).Error
	if error != nil {
		return error
	}

	if request.Status != "draft" {
		return error
	}

	var taskRequest ds.TaskRequest
	err := r.db.Where("request_id = ? AND task_id = ?", idR, idT).First(&taskRequest).Error
	if err != nil {
		return err
	}
	return r.db.Delete(&taskRequest).Error
}
