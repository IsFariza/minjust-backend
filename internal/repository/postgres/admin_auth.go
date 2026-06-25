package postgres

import (
	"database/sql"
	"errors"
	"minjust-website/internal/domain"
)

type postgresAuthRepository struct {
	db *sql.DB
}

func NewPostgresAuthRepository(db *sql.DB) domain.AuthRepository {
	return &postgresAuthRepository{db: db}
}

func (r *postgresAuthRepository) GetByUsername(username string) (*domain.Admin, error) {
	query := `SELECT id, username, password_hash FROM admin_accounts WHERE username = $1`

	var admin domain.Admin
	err := r.db.QueryRow(query, username).Scan(&admin.ID, &admin.Username, &admin.PasswordHash)
	if err == sql.ErrNoRows {
		return nil, errors.New("админ не найден")
	}
	return &admin, err
}
