package ds

import (
	"awesomeProject/internal/app/role"
	"time"

	"github.com/google/uuid"
)

const TASK_STATUS_ACTIVE = "active"
const TASK_STATUS_DELETED = "deleted"
const USER_ROLE_MODERATOR = "admin"

type Task struct {
	Task_id         uint `gorm:"primarykey;autoIncrement"`
	Name            string
	Subject         string
	Minidescription string
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

type User struct {
	User_id  uint `gorm:"primarykey;"`
	Name     string
	Phone    string
	Email    string
	Password string
	// Role         role.Role `sql:"type:string;"`
	Role string
}

type Status struct {
	Status string
}

type UUser struct {
	UUID uuid.UUID `gorm:"type:uuid"`
	Name string    `json:"name"`
	Role role.Role `sql:"type:string;"`
	Pass string
}

type TaskRequest struct {
	Task_id    int `gorm:"primarykey"`
	Request_id int `gorm:"primarykey"`
	Order int
}