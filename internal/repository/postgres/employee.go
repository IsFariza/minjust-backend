package postgres

import (
	"database/sql"
	"errors"
	"minjust-website/internal/domain"
	"strings"
)

type postgresEmployeeRepository struct {
	db *sql.DB
}

func NewPostgresEmployeeRepository(db *sql.DB) domain.EmployeeRepository {
	return &postgresEmployeeRepository{db: db}
}

func (r *postgresEmployeeRepository) GetManagement() ([]domain.Employee, error) {
	query := `SELECT id, fullname, position, email, phone, area_of_work, photo_url, created_at FROM management ORDER BY id ASC`

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var managers []domain.Employee
	for rows.Next() {
		var emp domain.Employee
		err := rows.Scan(&emp.ID, &emp.FullName, &emp.Position, &emp.Email, &emp.Phone, &emp.AreaOfWork, &emp.PhotoURL, &emp.CreatedAt)
		if err != nil {
			return nil, err
		}
		managers = append(managers, emp)
	}

	return managers, rows.Err()
}

func (r *postgresEmployeeRepository) CreateAccount(emp *domain.EmployeeAccount) error {
	query := `
        INSERT INTO employee_accounts (iin, fullname, email, phone, department, position)
        VALUES ($1, $2, $3, $4, $5, $6)
        RETURNING id, created_at`

	err := r.db.QueryRow(
		query,
		emp.IIN,
		emp.FullName,
		emp.Email,
		emp.Phone,
		emp.Department,
		emp.Position,
	).Scan(&emp.ID, &emp.CreatedAt)

	if err != nil {
		if strings.Contains(err.Error(), "duplicate key") || strings.Contains(err.Error(), "unique_iin") {
			return errors.New("employee with this IIN already exists")
		}
		return err
	}

	return nil
}

func (r *postgresEmployeeRepository) GetAllAccounts() ([]domain.EmployeeAccount, error) {
	query := `SELECT id, iin, fullname, email, phone, department, position, created_at FROM employee_accounts ORDER BY id DESC`

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var accounts []domain.EmployeeAccount
	for rows.Next() {
		var acc domain.EmployeeAccount
		err := rows.Scan(&acc.ID, &acc.IIN, &acc.FullName, &acc.Email, &acc.Phone, &acc.Department, &acc.Position, &acc.CreatedAt)
		if err != nil {
			return nil, err
		}
		accounts = append(accounts, acc)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return accounts, nil
}
