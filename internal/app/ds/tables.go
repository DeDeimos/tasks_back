package ds

import "time"

const TASK_STATUS_ACTIVE = "acti"
const TASK_STATUS_DELETED = "del"

type Task struct {
	Task_id         uint `gorm:"primarykey;autoIncrement"`
	Name            string
	Subject         string
	MiniDescription string
	Image           string
	Description     string
	Status          string
	Requests        []Request `gorm:"many2many:tasks_requests"`
}

type Request struct {
	ID             uint `gorm:"primarykey;autoIncrement"`
	Status         string
	start_date     time.Time `json:"start_date"`
	formation_date time.Time `json:"start_date"`
	end_date       time.Time `json:"start_date"`
	user_ID        uint
	moderator_ID   uint
	Tasks          []Task `gorm:"many2many:tasks_requests"`
}

type User struct {
	user_id      uint `gorm:"primarykey;autoIncrement"`
	name         string
	phoneNumber  string
	emailAddress string
	password     string
	role         string
	Requests     []Request
}
