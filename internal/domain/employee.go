package domain

import (
	"time"
)

type Employee struct {
	ID         int64     `json:"id"`
	FullName   string    `json:"fullname"`
	Position   string    `json:"position"`
	Email      string    `json:"email"`
	Phone      string    `json:"phone"`
	AreaOfWork string    `json:"area_of_work"`
	PhotoURL   string    `json:"photo_url,omitempty"`
	CreatedAt  time.Time `json:"created_at"`
}

type EmployeeAccount struct {
	ID            int64      `json:"id"`
	FullName      string     `json:"fullname" binding:"required"`
	IIN           string     `json:"iin" binding:"required,len=12"`
	Position      string     `json:"position" binding:"required"`
	Department    string     `json:"department" binding:"required"`
	Management    string     `json:"management" binding:"required"`
	Cabinet       string     `json:"cabinet"`
	PhoneWork     string     `json:"phone_work"`
	PhonePersonal string     `json:"phone_personal" binding:"required"`
	Email         string     `json:"email" binding:"required,email"`
	Password      string     `json:"password,omitempty" binding:"required,min=6"`
	CreatedAt     *time.Time `json:"created_at,omitempty"`
}

type EmployeeRepository interface {
	GetManagement() ([]Employee, error)
	CreateAccount(employee *EmployeeAccount) error
	GetAllAccounts() ([]EmployeeAccount, error)
	GetByID(id int64) (*EmployeeAccount, error)
	GetByIIN(iin string) (*EmployeeAccount, error)
}

type EmployeeUsecase interface {
	GetManagementHandbook() ([]Employee, error)
	CreateAccount(employee *EmployeeAccount) error
	GetAllAccounts() ([]EmployeeAccount, error)
	GetProfile(id int64) (*EmployeeAccount, error)
}
