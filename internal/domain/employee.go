package domain

import "time"

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
	ID         int64     `json:"id"`
	IIN        string    `json:"iin"`
	FullName   string    `json:"fullname"`
	Email      string    `json:"email"`
	Phone      string    `json:"phone"`
	Department string    `json:"department"`
	Position   string    `json:"position"`
	CreatedAt  time.Time `json:"created_at"`
}

type EmployeeRepository interface {
	GetManagement() ([]Employee, error)
	CreateAccount(employee *EmployeeAccount) error
	GetAllAccounts() ([]EmployeeAccount, error)
}

type EmployeeUsecase interface {
	GetManagementHandbook() ([]Employee, error)
	CreateAccount(employee *EmployeeAccount) error
	GetAllAccounts() ([]EmployeeAccount, error)
}
