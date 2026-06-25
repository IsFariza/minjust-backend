package postgres

import (
	"database/sql"
	"errors"
	"minjust-website/internal/domain"
	"strings"

	"github.com/lib/pq"
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
        INSERT INTO employee_accounts (fullname, iin, position, department, management, cabinet, phone_work, phone_personal, email, password)
        VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
        RETURNING id, created_at`

	err := r.db.QueryRow(
		query,
		emp.FullName,
		emp.IIN,
		emp.Position,
		emp.Department,
		emp.Management,
		emp.Cabinet,
		emp.PhoneWork,
		emp.PhonePersonal,
		emp.Email,
		emp.Password,
	).Scan(&emp.ID, &emp.CreatedAt)

	if err != nil {
		if pgErr, ok := err.(*pq.Error); ok {
			if pgErr.Code == "23505" { // UniqueViolation
				if strings.Contains(pgErr.Message, "iin") {
					return errors.New("этот ИИН уже зарегистрирован")
				}
				if strings.Contains(pgErr.Message, "email") {
					return errors.New("этот email уже зарегистрирован")
				}
			}
		}
		return err
	}

	return nil
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

func (r *postgresEmployeeRepository) GetByID(id int64) (*domain.EmployeeAccount, error) {
	query := `
        SELECT 
            id, iin, fullname, department, management, position, 
            COALESCE(cabinet, ''), 
            COALESCE(phone_personal, ''), 
            COALESCE(phone_work, ''), 
            email,
            created_at 
        FROM employee_accounts 
        WHERE id = $1`

	row := r.db.QueryRow(query, id)

	var emp domain.EmployeeAccount
	err := row.Scan(
		&emp.ID,
		&emp.IIN,
		&emp.FullName,
		&emp.Department,
		&emp.Management,
		&emp.Position,
		&emp.Cabinet,
		&emp.PhonePersonal,
		&emp.PhoneWork,
		&emp.Email,
		&emp.CreatedAt,
	)
	if err != nil {
		return nil, err
	}
	return &emp, nil
}

func (r *postgresEmployeeRepository) GetByIIN(iin string) (*domain.EmployeeAccount, error) {
	query := `
        SELECT id, iin, fullname, department, management, position,
            COALESCE(cabinet, ''),
            COALESCE(phone_personal, ''),
            COALESCE(phone_work, ''),
            email,
            password,
            created_at
        FROM employee_accounts 
        WHERE iin = $1`

	row := r.db.QueryRow(query, iin)

	var emp domain.EmployeeAccount
	err := row.Scan(
		&emp.ID,
		&emp.IIN,
		&emp.FullName,
		&emp.Department,
		&emp.Management,
		&emp.Position,
		&emp.Cabinet,
		&emp.PhonePersonal,
		&emp.PhoneWork,
		&emp.Email,
		&emp.Password,
		&emp.CreatedAt,
	)
	if err != nil {
		return nil, err
	}
	return &emp, nil
}
