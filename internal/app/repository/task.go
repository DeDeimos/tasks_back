package repository

import (
	"awesomeProject/internal/app/ds"
	"errors"
	"log"
	"time"

	"gorm.io/gorm"
)

func (r *Repository) GetTaskByID(id uint) (*ds.Task, error) {
	task := &ds.Task{}

	err := r.db.First(task, "task_id = ?", id).Error // find product with code D42
	if err != nil {
		return nil, err
	}

	return task, nil
}

func (r *Repository) GetTasksByRequestID(id int) ([]ds.Task, error) {
	var tasks []ds.Task
	var taskRequests []ds.TaskRequest
	err := r.db.Where("requestid = ?", id).Find(&taskRequests).Error
	if err != nil {
		return nil, err
	}
	for _, taskRequest := range taskRequests {
		var task ds.Task
		err = r.db.Where("id = ?", taskRequest.Task_id).First(&task).Error
		if err != nil {
			return nil, err
		}
		tasks = append(tasks, task)
	}
	return tasks, nil
}

func (r *Repository) DeleteTask(id uint) error {
	return r.db.Exec("UPDATE tasks SET status = ? WHERE task_id = ?", ds.TASK_STATUS_DELETED, id).Error
}

func (r *Repository) CreateTask(task ds.Task) error {
	return r.db.Create(&task).Error
}

func (r *Repository) GetAllTasks(status string, subject string) ([]ds.Task, error) {
	// var tasks []ds.Task
	// var request ds.Request
	// // var userID = 1
	// err := r.db.Find(&tasks, "status = 'active'").Error
	// if err != nil {
	// 	return nil, 0, err
	// }
	// err = r.db.Where("status = ? ", "active").First(&request).Error
	// if err != nil {
	// 	return tasks, nil
	// }

	// return tasks, nil
	var tasks []ds.Task
	query := r.db.Table("tasks").Where("status = ?", status).Where("lower(subject) LIKE ?", "%"+subject+"%")
	if err := query.Find(&tasks).Error; err != nil {
		return nil, err
	}
	return tasks, nil
}

func (r *Repository) UpdateTask(id int, task ds.Task) error {
	// Проверяем, существует ли консультация с указанным ID.
	existingTask, err := r.GetTaskByID(uint(id))
	if err != nil {
		return err // Возвращаем ошибку, если консультация не найдена.
	}

	log.Println(id)
	log.Println(existingTask)

	// Обновляем поля существующей консультации.
	existingTask.Name = task.Name
	existingTask.Description = task.Description
	existingTask.Image = task.Image
	existingTask.Status = task.Status

	log.Println(existingTask)

	// Сохраняем обновленную консультацию в базу данных.
	if err := r.db.Model(ds.Task{}).Where("task_id = ?", id).Updates(existingTask).Error; err != nil {
		return err
	}
	return nil
}

func (r *Repository) AddTaskToRequest(taskID int, userID int) error {
	var request ds.Request

	err := r.db.Where("status = ? AND user_id = ?", "draft", userID).FirstOrCreate(&request).Error
	if err != nil {
		return err
	}

	log.Println(request)

	// Если request был создан, установите дополнительные поля
	if request.Status != "draft" {
		request.Status = "draft"
		request.StartDate = time.Now()
		request.UserID = uint(userID)
		err = r.db.Save(&request).Error
		if err != nil {
			return err
		}
	}

	log.Println(request)
	log.Println(taskID)
	var taskRequest ds.TaskRequest
	err = r.db.Where("request_id = ? AND task_id = ?", request.Request_id, taskID).First(&taskRequest).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			// Создайте новую связь между task и request
			taskRequest = ds.TaskRequest{
				Task_id:    taskID,
				Request_id: int(request.Request_id),
			}
			err = r.db.Create(&taskRequest).Error
			if err != nil {
				return err
			}
			log.Println("1")
			return nil
		}
		return err
	}

	log.Println("123132")
	return nil
}

func (r *Repository) AddTaskImage(id int, imageBytes []byte, contentType string) error {
	// Удаление существующего изображения (если есть)
	log.Println(0)
	err := r.minioClient.RemoveServiceImage(id)
	if err != nil {
		// return err
	}
	log.Println(1)
	// Загрузка нового изображения в MinIO
	imageURL, err := r.minioClient.UploadServiceImage(id, imageBytes, contentType)
	if err != nil {
		return err
	}
	log.Println(2)
	// Обновление информации об изображении в БД (например, ссылки на MinIO)
	err = r.db.Model(&ds.Task{}).Where("task_id = ?", id).Update("image", imageURL).Error
	if err != nil {
		return err
	}

	return nil
}
