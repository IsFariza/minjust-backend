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
		employee.PhonePersonal == "" || employee.Email == "" || employee.Password == "" {
		return errors.New("все обязательные поля должны быть заполнены")
	}

	matched, _ := regexp.MatchString(`^\d{12}$`, employee.IIN)
	if !matched {
		return errors.New("ИИН должен содержать ровно 12 цифр")
	}

	if _, err := mail.ParseAddress(employee.Email); err != nil {
		return errors.New("неверный формат email")
	}
	if len(employee.Password) < 6 {
		return errors.New("пароль учетной записи должен содержать минимум 6 символов")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(employee.Password), bcrypt.DefaultCost)
	if err != nil {
		return errors.New("не удалось захешировать пароль")
	}
	employee.Password = string(hashedPassword)

	return u.repo.CreateAccount(employee)
}

func (u *employeeUsecase) GetAllAccounts() ([]domain.EmployeeAccount, error) {
	return u.repo.GetAllAccounts()
}

func (u *employeeUsecase) GetProfile(id int64) (*domain.EmployeeAccount, error) {
	return u.repo.GetByID(id)
}
