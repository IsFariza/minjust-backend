package postgres

import (
	"database/sql"
	"errors"
	"minjust-website/internal/domain"
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

	var existsIIN bool
	err := r.db.QueryRow("SELECT EXISTS(SELECT 1 FROM employee_accounts WHERE iin=$1)", emp.IIN).Scan(&existsIIN)
	if existsIIN {
		return errors.New("this IIN already exists")
	}
	var existsEmail bool
	err = r.db.QueryRow("SELECT EXISTS(SELECT 1 FROM employee_accounts WHERE email=$1)", emp.Email).Scan(&existsEmail)
	if existsEmail {
		return errors.New("this email already exists")
	}

	query := `
        INSERT INTO employee_accounts (fullname, iin, position, department, management, cabinet, phone_work, phone_personal, email, password)
        VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
        RETURNING id, created_at`

	err = r.db.QueryRow(query, emp.FullName, emp.IIN, emp.Position, emp.Department, emp.Management, emp.Cabinet, emp.PhoneWork, emp.PhonePersonal, emp.Email, emp.Password).Scan(&emp.ID, &emp.CreatedAt)
	return err
}

func (r *postgresEmployeeRepository) GetAllAccounts() ([]domain.EmployeeAccount, error) {
	query := `SELECT id, fullname, iin, position, department, management, cabinet, phone_work, phone_personal, email, created_at FROM employee_accounts ORDER BY id DESC`

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var accounts []domain.EmployeeAccount
	for rows.Next() {
		var acc domain.EmployeeAccount
		err := rows.Scan(
			&acc.ID,
			&acc.FullName,
			&acc.IIN,
			&acc.Position,
			&acc.Department,
			&acc.Management,
			&acc.Cabinet,
			&acc.PhoneWork,
			&acc.PhonePersonal,
			&acc.Email,
			&acc.CreatedAt,
		)
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
