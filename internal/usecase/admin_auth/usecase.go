package admin_auth

import (
	"errors"
	"minjust-website/internal/auth"
	"minjust-website/internal/domain"
	"time"

	"golang.org/x/crypto/bcrypt"
)

type authUsecase struct {
	adminRepo    domain.AuthRepository
	employeeRepo domain.EmployeeRepository
	jwtSecret    string
}

func NewAuthUsecase(adminRepo domain.AuthRepository, employeeRepo domain.EmployeeRepository, jwtSecret string) domain.AuthUsecase {
	return &authUsecase{
		adminRepo:    adminRepo,
		employeeRepo: employeeRepo,
		jwtSecret:    jwtSecret}
}

func (u *authUsecase) Login(username, password, role string) (string, error) {
	var userID int64
	var hashedPassword string
	var userIIN string

	if role == "admin" {
		admin, err := u.adminRepo.GetByUsername(username)
		if err != nil {
			return "", errors.New("неверный пароль или логин")
		}
		userID = admin.ID
		hashedPassword = admin.PasswordHash
		userIIN = ""
	} else {
		emp, err := u.employeeRepo.GetByIIN(username)
		if err != nil {
			return "", errors.New("неверный пароль или логин")
		}
		userID = emp.ID
		hashedPassword = emp.Password
		userIIN = emp.IIN
	}

	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	if err != nil {
		return "", errors.New("неверный пароль или логин")
	}

	token, err := auth.GenerateToken(userID, role, userIIN, u.jwtSecret, 12*time.Hour)
	if err != nil {
		return "", errors.New("could not generate access token")
	}

	return token, nil
}
