package repository

import (
	"awesomeProject/internal/app/ds"
	"log"
	"time"

	"gorm.io/gorm"
)

// func (r *Repository) GetRequestByID(id int) (*ds.Request, error) {
// 	request := &ds.Request{}

// 	err := r.db.Preload("Tasks").
// 		Joins("JOIN task_requests ON requests.request_id = task_requests.request_id").
// 		Preload("Moderator", func(db *gorm.DB) *gorm.DB {
// 			return db.Select("user_id, name, email")
// 		}).
// 		Preload("User", func(db *gorm.DB) *gorm.DB {
// 			return db.Select("user_id, name, email")
// 		}).
// 		Where("requests.request_id = ?", id).
// 		Order("task_requests.order ASC").
// 		Find(request).Error
// 	if err != nil {
// 		return nil, err
// 	}

// 	return request, nil
// }

func (r *Repository) GetRequestByID(id int) (*ds.Request, error) {
	request := &ds.Request{}

	err := r.db.Preload("Tasks", func(db *gorm.DB) *gorm.DB {
			return db.Joins("JOIN task_requests ON tasks.task_id = task_requests.task_id").
				Order("task_requests.order ASC").
				Where("task_requests.request_id = ?", id)
		}).
		Preload("Moderator", func(db *gorm.DB) *gorm.DB {
			return db.Select("user_id, name, email")
		}).
		Preload("User", func(db *gorm.DB) *gorm.DB {
			return db.Select("user_id, name, email")
		}).
		Where("requests.request_id = ?", id).
		First(request).
		Error

	if err != nil {
		return nil, err
	}

	return request, nil
}



func (r *Repository) DeleteRequest(id int) error {
	return r.db.Exec("UPDATE requests SET status = 'deleted' WHERE request_id=?", id).Error
}

func (r *Repository) GetAllRequests() ([]ds.Request, error) {
	log.Println("empty status")
	var requests []ds.Request
	err := r.db.Find(&requests, "status <> 'deleted'").Error
	if err != nil {
		return nil, err
	}
	log.Println(requests)
	return requests, nil
}

// func (r *Repository) FindAllByUserID(userID uint, status string, timeFrom *time.Time, timeTo *time.Time) ([]ds.Request, error) {
// 	log.Println("i am user")
// 	log.Println(status)
// 	requests := make([]ds.Request, 0)
// 	if timeFrom == nil && timeTo == nil {
// 		if status == "" {
// 			err := r.db.
// 				Preload("Moderator", func(db *gorm.DB) *gorm.DB {
// 					return db.Select("user_id, name, email")
// 				}).
// 				Preload("User", func(db *gorm.DB) *gorm.DB {
// 					return db.Select("user_id, name, email")
// 				}).
// 				Find(&requests, "user_id = ? and status <> 'deleted'", userID).Error
// 			if err != nil {
// 				return nil, err
// 			}
// 			return requests, nil
// 		}
// 		err := r.db.
// 			Preload("Moderator", func(db *gorm.DB) *gorm.DB {
// 				return db.Select("user_id, name, email")
// 			}).
// 			Preload("User", func(db *gorm.DB) *gorm.DB {
// 				return db.Select("user_id, name, email")
// 			}).
// 			Find(&requests, "user_id = ? AND status = ? and status <> 'deleted'", userID, status).Error
// 		if err != nil {
// 			return nil, err
// 		}
// 		return requests, nil
// 	}
// 	query := r.db.
// 		Preload("Moderator", func(db *gorm.DB) *gorm.DB {
// 			return db.Select("user_id, name, email")
// 		}).
// 		Preload("User", func(db *gorm.DB) *gorm.DB {
// 			return db.Select("user_id, name, email")
// 		}).
// 		Table("requests").Where("user_id = ?", userID).Where("status = ? and status <> 'deleted'", status).Where("formation_date >= ?", timeFrom).Where("formation_date <= ?", timeTo).Order("start_date DESC")
// 	if err := query.Find(&requests).Error; err != nil {
// 		return nil, err
// 	}
// 	return requests, nil
// }

func (r *Repository) FindAllByUserID(userID uint, status string, timeFrom *time.Time, timeTo *time.Time) ([]ds.Request, error) {
    requests := make([]ds.Request, 0)

    query := r.db.
        Preload("Moderator", func(db *gorm.DB) *gorm.DB {
            return db.Select("user_id, name, email")
        }).
        Preload("User", func(db *gorm.DB) *gorm.DB {
            return db.Select("user_id, name, email")
        }).
        Table("requests").
        Where("user_id = ?", userID).
        Where("status <> 'deleted'").
		Where("status <> 'draft'")

    if status != "" {
        query = query.Where("status = ?", status)
    }

    if timeFrom != nil {
        query = query.Where("formation_date >= ?", timeFrom)
    }

    if timeTo != nil {
        query = query.Where("formation_date <= ?", timeTo)
    }

    if err := query.Order("start_date DESC").Find(&requests).Error; err != nil {
        return nil, err
    }

    return requests, nil
}

func (r *Repository) FindAllByModeratorID(moderatorID uint, status string, timeFrom *time.Time, timeTo *time.Time) ([]ds.Request, error) {
	log.Println("i am admin")
	requests := make([]ds.Request, 0)
	
    query := r.db.
        Preload("Moderator", func(db *gorm.DB) *gorm.DB {
            return db.Select("user_id, name, email")
        }).
        Preload("User", func(db *gorm.DB) *gorm.DB {
            return db.Select("user_id, name, email")
        }).
        Table("requests").
        Where("status <> 'deleted'").
		Where("status <> 'draft'")

    if status != "" {
        query = query.Where("status = ?", status)
    }

    if timeFrom != nil {
        query = query.Where("formation_date >= ?", timeFrom)
    }

    if timeTo != nil {
        query = query.Where("formation_date <= ?", timeTo)
    }

    if err := query.Order("start_date DESC").Find(&requests).Error; err != nil {
        return nil, err
    }

    return requests, nil
	// if timeFrom == nil && timeTo == nil {
	// 	err := r.db.
	// 		Preload("Moderator", func(db *gorm.DB) *gorm.DB {
	// 			return db.Select("user_id, name, email")
	// 		}).
	// 		Preload("User", func(db *gorm.DB) *gorm.DB {
	// 			return db.Select("user_id, name, email")
	// 		}).
	// 		Find(&requests, "status <> 'deleted'").Error
	// 	// Table("requests").Where("? = '' OR status = ?", status, status).Error
	// 	if err != nil {
	// 		return nil, err
	// 	}
	// 	return requests, nil
	// }
	// query := r.db.
	// 	Preload("Moderator", func(db *gorm.DB) *gorm.DB {
	// 		return db.Select("user_id, name, email")
	// 	}).
	// 	Preload("User", func(db *gorm.DB) *gorm.DB {
	// 		return db.Select("user_id, name, email")
	// 	}).
	// 	Table("requests").Where("? = '' OR status = ?", status, status).Where("formation_date >= ?", timeFrom).Where("formation_date <= ?", timeTo).Order("created_at DESC")
	// if err := query.Find(&requests).Error; err != nil {
	// 	return nil, err
	// }
	// return requests, nil
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

func (r *Repository) GetDraftUser(id int) int {
	var request ds.Request

	err := r.db.Where("status = ? AND user_id = ?", "draft", id).First(&request).Error
	if err != nil {
		return -1
	}
	// log.Println(request)

	// Если request был создан, установите дополнительные поля
	// if request.Status != "draft" {
	// 	request.Status = "draft"
	// 	request.StartDate = time.Now()
	// 	request.UserID = uint(id)
	// 	request.ModeratorID = uint(1)
	// 	err = r.db.Save(&request).Error
	// }

	return int(request.Request_id)

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

func (r *Repository) UpdateUserRequestStatus(id int, status string) error {
	// Проверяем, существует ли занятие с указанным ID.
	existingRequest, err := r.GetRequestByID(id)
	if err != nil {
		return err // Возвращаем ошибку, если занятие не найдена.
	}
	// Обновляем поля существующего занятия.
	existingRequest.Status = status
	existingRequest.FormationDate = time.Now()

	// Сохраняем обновленное занятие в базу данных.
	if err := r.db.Model(ds.Request{}).Where("request_id = ?", id).Updates(existingRequest).Error; err != nil {
		return err
	}
	return nil
}

func (r *Repository) UpdateAdminRequestStatus(id int, status string, userID uint) error {
	// Проверяем, существует ли занятие с указанным ID.
	existingRequest, err := r.GetRequestByID(id)
	if err != nil {
		return err // Возвращаем ошибку, если занятие не найдена.
	}
	// Обновляем поля существующего занятия.
	existingRequest.Status = status
	existingRequest.EndDate = time.Now()
	existingRequest.ModeratorID = userID

	// Сохраняем обновленное занятие в базу данных.
	if err := r.db.Model(ds.Request{}).Where("request_id = ?", id).Updates(existingRequest).Error; err != nil {
		return err
	}
	return nil
}
