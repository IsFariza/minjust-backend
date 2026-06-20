package employee

import (
	"errors"
	"minjust-website/internal/domain"
	"net/mail"
	"regexp"
	"strings"
)

type employeeUsecase struct {
	repo domain.EmployeeRepository
}

func NewEmployeeUsecase(repo domain.EmployeeRepository) domain.EmployeeUsecase {
	return &employeeUsecase{repo: repo}
}

func (u *employeeUsecase) GetManagementHandbook() ([]domain.Employee, error) {
	return u.repo.GetManagement()
}

func (u *employeeUsecase) CreateAccount(employee *domain.EmployeeAccount) error {
	employee.IIN = strings.TrimSpace(employee.IIN)
	employee.FullName = strings.TrimSpace(employee.FullName)
	employee.Email = strings.TrimSpace(strings.ToLower(employee.Email))
	employee.Phone = strings.TrimSpace(employee.Phone)
	employee.Department = strings.TrimSpace(employee.Department)
	employee.Position = strings.TrimSpace(employee.Position)

	if employee.FullName == "" || employee.Email == "" || employee.Phone == "" || employee.Department == "" || employee.Position == "" {
		return errors.New("all required fields must be filled")
	}

	matched, _ := regexp.MatchString(`^\d{12}$`, employee.IIN)
	if !matched {
		return errors.New("iin must contain exactly 12 digits")
	}

	if _, err := mail.ParseAddress(employee.Email); err != nil {
		return errors.New("email has invalid format")
	}

	return u.repo.CreateAccount(employee)
}
func (u *employeeUsecase) GetAllAccounts() ([]domain.EmployeeAccount, error) {
	return u.repo.GetAllAccounts()
}
