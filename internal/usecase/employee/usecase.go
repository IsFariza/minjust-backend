package employee

import (
	"errors"
	"minjust-website/internal/domain"
	"net/mail"
	"regexp"
	"strings"

	"golang.org/x/crypto/bcrypt"
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
	employee.Position = strings.TrimSpace(employee.Position)
	employee.Department = strings.TrimSpace(employee.Department)
	employee.Management = strings.TrimSpace(employee.Management)
	employee.Cabinet = strings.TrimSpace(employee.Cabinet)
	employee.PhoneWork = strings.TrimSpace(employee.PhoneWork)
	employee.PhonePersonal = strings.TrimSpace(employee.PhonePersonal)
	employee.Email = strings.TrimSpace(strings.ToLower(employee.Email))

	if employee.IIN == "" || employee.FullName == "" || employee.Position == "" || employee.Department == "" || employee.Management == "" ||
		employee.PhonePersonal == "" || employee.Email == "" {
		return errors.New("all required fields must be filled")
	}

	matched, _ := regexp.MatchString(`^\d{12}$`, employee.IIN)
	if !matched {
		return errors.New("iin must contain exactly 12 digits")
	}

	if _, err := mail.ParseAddress(employee.Email); err != nil {
		return errors.New("email has invalid format")
	}
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(employee.Password), bcrypt.DefaultCost)
	if err != nil {
		return errors.New("failed to hash the password")
	}
	employee.Password = string(hashedPassword)

	return u.repo.CreateAccount(employee)
}
func (u *employeeUsecase) GetAllAccounts() ([]domain.EmployeeAccount, error) {
	return u.repo.GetAllAccounts()
}
