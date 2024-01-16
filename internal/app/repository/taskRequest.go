package repository

import (
	"awesomeProject/internal/app/ds"
	"log"
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
	err := r.db.Where("request_id = ?", idR).Order("\"order\"").Find(&TaskRequests).Error
	if err != nil {
		return err
	}
	log.Println(TaskRequests)
	log.Println(idT)
	var oldOrder int
	var updateIndex int


	oldState := 0

	for i, taskRequest := range TaskRequests {
		if taskRequest.Task_id == idT {
			updateIndex = i
			oldOrder = taskRequest.Order
		}
		if taskRequest.Order == newOrder {
			oldState = i
		}

	}
	log.Println("oldOrder", oldOrder)
	log.Println("newOrder", newOrder)
	TaskRequests[updateIndex].Order = newOrder
	
	if oldOrder > newOrder {
		for i := oldState; i < updateIndex; i++ {
			TaskRequests[i].Order++
		}
	} else {
		for i := updateIndex + 1; i <= oldState; i++ {
			TaskRequests[i].Order--
		}
	}

	log.Println(TaskRequests)

	for _, taskRequest := range TaskRequests {
		err := r.db.Save(&taskRequest).Error
		if err != nil {
			return err
		}
	}

	return nil

}

