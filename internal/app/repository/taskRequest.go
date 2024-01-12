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

func (r *Repository) ChangeOrder(idT int, idR int, newOrder int) error {

	var TaskRequests []ds.TaskRequest
	err := r.db.Where("request_id = ?", idR).Find(&TaskRequests).Order("order").Error
	if err != nil {
		return err
	}
	var oldOrder int
	var updateIndex int
	for i, taskRequest := range TaskRequests {
		if taskRequest.Task_id == idT {
			updateIndex = i
			oldOrder = taskRequest.Order
			break
		}
	}

	TaskRequests[updateIndex].Order = newOrder

	if newOrder > oldOrder {
		for i := updateIndex; i < len(TaskRequests)-1; i++ {
			TaskRequests[i+1].Order = TaskRequests[i+1].Order - 1
		}
	} else {
		for i := updateIndex; i > 0; i-- {
			TaskRequests[i-1].Order = TaskRequests[i-1].Order + 1
		}
	}

	for _, taskRequest := range TaskRequests {
		err = r.db.Save(&taskRequest).Error
		if err != nil {
			return err
		}
	}

	return nil

}
