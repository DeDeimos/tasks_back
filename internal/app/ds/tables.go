package ds

import "time"

const TASK_STATUS_ACTIVE = "active"
const TASK_STATUS_DELETED = "delete"
const USER_ROLE_MODERATOR = "admin"

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
	User          *User  `gorm:"foreignkey:UserID;references:user_id"`
	Moderator     *User  `gorm:"foreignkey:ModeratorID;references:user_id"`
	Tasks         []Task `gorm:"many2many:task_requests;foreignKey:request_id;joinForeignKey:request_id;References:task_id;JoinReferences:task_id"`
}

type TaskRequest struct {
	Task_id    int `gorm:"primarykey"`
	Request_id int `gorm:"primarykey"`
}

type User struct {
	User_id      uint `gorm:"primarykey;"`
	Name         string
	PhoneNumber  string
	EmailAddress string
	Password     string
	Role         string
}

type Status struct {
	Status string
}
