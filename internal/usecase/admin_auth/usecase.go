package admin_auth

import (
	"errors"
	"minjust-website/internal/auth"
	"minjust-website/internal/domain"
	"time"

	"golang.org/x/crypto/bcrypt"
)

type authUsecase struct {
	repo      domain.AuthRepository
	jwtSecret string
}

func NewAuthUsecase(repo domain.AuthRepository, jwtSecret string) domain.AuthUsecase {
	return &authUsecase{repo: repo, jwtSecret: jwtSecret}
}

func (u *authUsecase) Login(username, password string) (string, error) {
	admin, err := u.repo.GetByUsername(username)
	if err != nil {
		return "", errors.New("invalid username or password")
	}

	err = bcrypt.CompareHashAndPassword([]byte(admin.PasswordHash), []byte(password))
	if err != nil {
		return "", errors.New("invalid username or password")
	}

	token, err := auth.GenerateToken(admin.ID, "admin", u.jwtSecret, 12*time.Hour)
	if err != nil {
		return "", errors.New("could not generate access token")
	}

	return token, nil
}
