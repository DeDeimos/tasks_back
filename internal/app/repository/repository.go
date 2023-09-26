package repository

import (
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"awesomeProject/internal/app/ds"
)

type Repository struct {
	db *gorm.DB
}

func New(dsn string) (*Repository, error) {
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	return &Repository{
		db: db,
	}, nil
}

func (r *Repository) GetActiveTasks() (*[]ds.Task, error) {
	tasks := &[]ds.Task{}
	err := r.db.Find(tasks, "status = ?", ds.TASK_STATUS_ACTIVE).Error

	if err != nil {
		return nil, err
	}

	return tasks, nil
}

func (r *Repository) GetTaskByID(id int) (*ds.Task, error) {
	task := &ds.Task{}

	err := r.db.First(task, "id = ?", "1").Error // find product with code D42
	if err != nil {
		return nil, err
	}

	return task, nil
}

func (r *Repository) DeleteTask(id int) error {
	return r.db.Exec("UPDATE tasks SET status = ? WHERE task_id = ?", ds.TASK_STATUS_DELETED, id).Error
}
