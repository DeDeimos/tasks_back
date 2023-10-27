package repository

import (
	"awesomeProject/internal/app/ds"
	"time"
)

func (r *Repository) GetRequestByID(id int) (*ds.Request, error) {
	request := &ds.Request{}

	err := r.db.Preload("Tasks").First(request, "request_id = ?", id).Error
	if err != nil {
		return nil, err
	}

	return request, nil
}

func (r *Repository) DeleteRequest(id int) error {
	return r.db.Exec("UPDATE requests SET status = 'deleted' WHERE request_id=?", id).Error
}

func (r *Repository) GetAllRequests(user_id uint) ([]ds.Request, error) {
	var requests []ds.Request
	err := r.db.Preload("Tasks").Where("user_id = ? and status != ?", user_id, "delete").Find(&requests).Error
	if err != nil {
		return nil, err
	}

	return requests, nil
}

func (r *Repository) GetRequestsByStatus(status string) ([]ds.Request, error) {
	var requests []ds.Request
	err := r.db.Preload("Tasks").Where("status = ?", status).Find(&requests).Error
	if err != nil {
		return nil, err
	}

	return requests, nil
}

func (r *Repository) GetRequestsByDate(startDate time.Time, endDate time.Time) ([]ds.Request, error) {
	var requests []ds.Request
	if !endDate.IsZero() {
		err := r.db.Preload("Tasks").Where("formation_date >= ? AND formation_date <= ?", startDate, endDate).Find(&requests).Error
		if err != nil {
			return nil, err
		}
		return requests, nil
	}

	err := r.db.Where("formation_date >= ?", startDate).Find(&requests).Error
	if err != nil {
		return nil, err
	}
	return requests, nil
}

func (r *Repository) UpdateRequest(id int, request ds.Request) error {
	existingRequest, err := r.GetRequestByID(id)
	if err != nil {
		return err // Возвращаем ошибку, если занятие не найдено.
	}

	if err := r.db.Model(ds.Request{}).Where("request_id = ?", id).Updates(existingRequest).Error; err != nil {
		return err
	}
	return nil
}

func (r *Repository) UpdateRequestStatus(id int, status string) error {
	// Проверяем, существует ли занятие с указанным ID.
	existingRequest, err := r.GetRequestByID(id)
	if err != nil {
		return err // Возвращаем ошибку, если занятие не найдена.
	}

	// Обновляем поля существующего занятия.
	existingRequest.Status = status

	// Сохраняем обновленное занятие в базу данных.
	if err := r.db.Model(ds.Request{}).Where("request_id = ?", id).Updates(existingRequest).Error; err != nil {
		return err
	}
	return nil
}
