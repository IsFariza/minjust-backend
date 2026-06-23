package postgres

import (
	"database/sql"
	"encoding/json"
	"errors"
	"minjust-website/internal/domain"
)

type postgresPasswordRequestRepository struct {
	db *sql.DB
}

func NewPostgresPasswordRequestRepository(db *sql.DB) domain.PasswordRequestRepository {
	return &postgresPasswordRequestRepository{db: db}
}

func (r *postgresPasswordRequestRepository) Create(req *domain.PasswordRequest) error {
	var input struct {
		IIN string `json:"iin"`
	}
	if err := json.Unmarshal(req.InputData, &input); err != nil {
		return errors.New("failed to parse IIN from request data")
	}
	var employeeExists bool
	checkQuery := `SELECT EXISTS(SELECT 1 FROM employee_accounts WHERE iin = $1)`
	err := r.db.QueryRow(checkQuery, input.IIN).Scan(&employeeExists)
	if err != nil {
		return err
	}
	if !employeeExists {
		return errors.New("employee with this IIN does not exist")
	}
	query := `
		INSERT INTO password_requests (system_name, input_data, status)
		VALUES ($1, $2, $3)
		RETURNING  id, created_at, updated_at`

	return r.db.QueryRow(query, req.EmployeeID, req.SystemName, req.InputData, req.Status).
		Scan(&req.ID, &req.CreatedAt, &req.UpdatedAt)
}

func (r *postgresPasswordRequestRepository) GetByID(id int64) (*domain.PasswordRequest, error) {
	query := `
		SELECT id, employee_id, system_name, input_data, status, primary_password, admin_comment, created_at, updated_at
		FROM password_requests
		WHERE id = $1`

	var req domain.PasswordRequest
	var employeeID sql.NullInt64
	var primaryPassword sql.NullString
	var adminComment sql.NullString

	err := r.db.QueryRow(query, id).Scan(
		&req.ID,
		&employeeID,
		&req.SystemName,
		&req.InputData,
		&req.Status,
		&primaryPassword,
		&adminComment,
		&req.CreatedAt,
		&req.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, errors.New("request not found")
	}
	if err != nil {
		return nil, err
	}

	if employeeID.Valid {
		req.EmployeeID = &employeeID.Int64
	}
	if primaryPassword.Valid {
		req.PrimaryPassword = primaryPassword.String
	}
	if adminComment.Valid {
		req.AdminComment = adminComment.String
	}

	return &req, nil
}

func (r *postgresPasswordRequestRepository) UpdateStatus(id int64, status, encryptedPassword, comment string) error {
	query := `
		UPDATE password_requests
		SET status = $1, primary_password = NULLIF($2, ''), admin_comment = NULLIF($3, ''), updated_at = CURRENT_TIMESTAMP
		WHERE id = $4`

	result, err := r.db.Exec(query, status, encryptedPassword, comment, id)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return errors.New("request not found")
	}

	return nil
}
