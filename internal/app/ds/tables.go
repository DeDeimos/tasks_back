package ds

import "time"

const TASK_STATUS_ACTIVE = "active"
const TASK_STATUS_DELETED = "delete"

type Task struct {
	Task_id         uint `gorm:"primarykey;autoIncrement"`
	Name            string
	Subject         string
	MiniDescription string
	Image           string
	Description     string
	Status          string
}

type Request struct {
	Request_id    uint `gorm:"primarykey;autoIncrement"`
	Status        string
	StartDate     time.Time `json:"start_date"`
	FormationDate time.Time `json:"formation_date"`
	EndDate       time.Time `json:"end_date"`
	UserID        uint
	ModeratorID   uint
	Tasks         []Task `gorm:"many2many:task_requests;foreignKey:request_id;joinForeignKey:request_id;References:task_id;JoinReferences:task_id"`
}

type TaskRequest struct {
	Task_id    int `gorm:"primarykey"`
	Request_id int `gorm:"primarykey"`
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

type Status struct {
	Status string
}
