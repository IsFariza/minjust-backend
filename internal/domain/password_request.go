package domain

import (
	"encoding/json"
	"time"
)

const (
	StatusPending  = "pending"
	StatusApproved = "approved"
	StatusRejected = "rejected"
)

type PasswordRequest struct {
	ID              int64           `json:"id"`
	EmployeeID      *int64          `json:"employee_id,omitempty"`
	SystemName      string          `json:"system_name"`
	InputData       json.RawMessage `json:"input_data"`
	Status          string          `json:"status"`
	PrimaryPassword string          `json:"primary_password,omitempty"`
	AdminComment    string          `json:"admin_comment,omitempty"`
	CreatedAt       time.Time       `json:"created_at"`
	UpdatedAt       time.Time       `json:"updated_at"`
}

type PasswordRequestRepository interface {
	Create(req *PasswordRequest) error
	GetByID(id int64) (*PasswordRequest, error)
	UpdateStatus(id int64, status, encryptedPassword, comment string) error
}

type PasswordRequestUsecase interface {
	RequestResetPassword(req *PasswordRequest) error
	GetRequestStatus(id int64) (*PasswordRequest, error)
	ProcessRequest(id int64, status, password, comment string) error
}
