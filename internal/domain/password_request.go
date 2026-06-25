package domain

import "time"

type PasswordRequest struct {
	ID              int64     `json:"id"`
	EmployeeID      int64     `json:"employee_id"`
	SystemName      string    `json:"system_name"`
	Status          string    `json:"status"`
	PrimaryPassword string    `json:"primary_password,omitempty"`
	RejectionReason string    `json:"rejection_reason,omitempty"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`

	EmployeeName string `json:"employee_name,omitempty"`
	EmployeeIIN  string `json:"employee_iin,omitempty"`
}

type CreatePasswordRequestInput struct {
	SystemName  string `json:"system_name" binding:"required"`
	EmployeeIIN string `json:"employee_iin" binding:"required"`
}

type UpdatePasswordStatusInput struct {
	Status          string `json:"status" binding:"required"`
	PrimaryPassword string `json:"primary_password,omitempty"`
	RejectionReason string `json:"rejection_reason,omitempty"`
}

type PasswordRequestRepository interface {
	Create(req *PasswordRequest) error
	ExistsByEmployeeAndSystem(empID int64, systemName string) (bool, error)
	GetByID(id int64) (*PasswordRequest, error)
	GetByEmployeeID(empID int64) ([]PasswordRequest, error)
	GetAll() ([]PasswordRequest, error)
	UpdateStatus(id int64, status, password, reason string) error
}

type PasswordRequestUsecase interface {
	CreateRequest(empID int64, currentIIN, inputIIN, systemName string) error
	GetEmployeeRequests(empID int64) ([]PasswordRequest, error)
	GetAllRequests() ([]PasswordRequest, error)
	ReviewRequest(id int64, status, reason string) error
}
