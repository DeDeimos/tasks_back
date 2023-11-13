package repository

import (
	"awesomeProject/internal/app/ds"
)

func (r *Repository) FindByID(id int) (*ds.User, error) {
	user := &ds.User{}

	err := r.db.First(user, "user_id = ?", id).Error
	if err != nil {
		return nil, err
	}

	return user, nil
}
