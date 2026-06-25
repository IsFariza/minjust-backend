package postgres

import (
	"database/sql"
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
	query := `INSERT INTO password_requests (employee_id, system_name, status, input_data) 
	          VALUES ($1, $2, 'pending', '{}'::jsonb)`
	_, err := r.db.Exec(query, req.EmployeeID, req.SystemName)
	return err
}

func (r *postgresPasswordRequestRepository) ExistsByEmployeeAndSystem(empID int64, systemName string) (bool, error) {
	query := `SELECT EXISTS(
		SELECT 1 FROM password_requests
		WHERE employee_id = $1 AND system_name = $2
	)`

	var exists bool
	err := r.db.QueryRow(query, empID, systemName).Scan(&exists)
	return exists, err
}

func (r *postgresPasswordRequestRepository) GetByID(id int64) (*domain.PasswordRequest, error) {
	query := `SELECT r.id, r.employee_id, r.system_name, r.status, COALESCE(r.primary_password, ''),
	          COALESCE(r.admin_comment, ''), r.created_at, COALESCE(r.updated_at, r.created_at),
	          e.fullname, e.iin
	          FROM password_requests r
	          JOIN employee_accounts e ON r.employee_id = e.id
	          WHERE r.id = $1`

	var req domain.PasswordRequest
	err := r.db.QueryRow(query, id).Scan(
		&req.ID,
		&req.EmployeeID,
		&req.SystemName,
		&req.Status,
		&req.PrimaryPassword,
		&req.RejectionReason,
		&req.CreatedAt,
		&req.UpdatedAt,
		&req.EmployeeName,
		&req.EmployeeIIN,
	)
	if err == sql.ErrNoRows {
		return nil, errors.New("заявка не найдена")
	}
	if err != nil {
		return nil, err
	}
	return &req, nil
}

func (r *postgresPasswordRequestRepository) GetByEmployeeID(empID int64) ([]domain.PasswordRequest, error) {
	query := `SELECT id, employee_id, system_name, status, COALESCE(primary_password, ''), COALESCE(admin_comment, ''), created_at, 
	          COALESCE(updated_at, created_at) FROM password_requests WHERE employee_id = $1 ORDER BY created_at DESC`

	rows, err := r.db.Query(query, empID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var requests []domain.PasswordRequest
	for rows.Next() {
		var req domain.PasswordRequest
		err := rows.Scan(&req.ID, &req.EmployeeID, &req.SystemName, &req.Status, &req.PrimaryPassword, &req.RejectionReason, &req.CreatedAt, &req.UpdatedAt)
		if err != nil {
			return nil, err
		}
		requests = append(requests, req)
	}
	return requests, rows.Err()
}

func (r *postgresPasswordRequestRepository) GetAll() ([]domain.PasswordRequest, error) {
	query := `SELECT r.id, r.employee_id, r.system_name, r.status, COALESCE(r.primary_password, ''), COALESCE(r.admin_comment, ''), r.created_at, COALESCE(r.updated_at, r.created_at), e.fullname, e.iin 
	          FROM password_requests r
	          JOIN employee_accounts e ON r.employee_id = e.id
	          ORDER BY r.created_at DESC`

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var requests []domain.PasswordRequest
	for rows.Next() {
		var req domain.PasswordRequest
		err := rows.Scan(
			&req.ID,
			&req.EmployeeID,
			&req.SystemName,
			&req.Status,
			&req.PrimaryPassword,
			&req.RejectionReason,
			&req.CreatedAt,
			&req.UpdatedAt,
			&req.EmployeeName,
			&req.EmployeeIIN,
		)
		if err != nil {
			return nil, err
		}
		requests = append(requests, req)
	}
	return requests, rows.Err()
}
func (r *postgresPasswordRequestRepository) UpdateStatus(id int64, status, password, reason string) error {
	query := `UPDATE password_requests
	          SET status = $1, primary_password = NULLIF($2, ''), admin_comment = NULLIF($3, ''), updated_at = NOW()
	          WHERE id = $4`
	result, err := r.db.Exec(query, status, password, reason, id)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return errors.New("заявка не найдена")
	}
	return err
}
