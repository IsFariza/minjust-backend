package password_request

import (
	"errors"
	"fmt"
	"minjust-website/internal/crypto"
	"minjust-website/internal/domain"
	"os"
	"regexp"
	"strings"
)

const (
	statusPending  = "pending"
	statusApproved = "approved"
	statusRejected = "rejected"
)

type passwordRequestUsecase struct {
	repo domain.PasswordRequestRepository
}

var supportedSystems = map[string]string{
	"documentolog":        "Documentolog",
	"ис documentolog":     "Documentolog",
	"документолог":        "Documentolog",
	"e-otinish":           "E-otinish",
	"e.otinish":           "E-otinish",
	"е-өтініш":            "E-otinish",
	"e-өтініш":            "E-otinish",
	"eps":                 "EPS",
	"епс":                 "EPS",
	"e-kyzmet":            "E-kyzmet",
	"e.kyzmet":            "E-kyzmet",
	"е-қызмет":            "E-kyzmet",
	"e-қызмет":            "E-kyzmet",
	"ad":                  "AD",
	"active directory":    "AD",
	"активный каталог":    "AD",
	"активная директория": "AD",
}

func NewPasswordRequestUsecase(repo domain.PasswordRequestRepository) domain.PasswordRequestUsecase {
	return &passwordRequestUsecase{repo: repo}
}

func (u *passwordRequestUsecase) CreateRequest(empID int64, currentIIN, inputIIN, systemName string) error {
	inputIIN = strings.TrimSpace(inputIIN)
	currentIIN = strings.TrimSpace(currentIIN)

	if !isValidIIN(inputIIN) {
		return errors.New("ИИН должен состоять ровно из 12 цифр")
	}
	if currentIIN != inputIIN {
		return errors.New("введенный ИИН не совпадает с ИИН текущего пользователя")
	}

	normalizedSystem, err := normalizeSystemName(systemName)
	if err != nil {
		return err
	}

	exists, err := u.repo.ExistsByEmployeeAndSystem(empID, normalizedSystem)
	if err != nil {
		return err
	}
	if exists {
		return errors.New("вы уже создали заявку на сброс пароля для этой системы")
	}

	req := &domain.PasswordRequest{
		EmployeeID: empID,
		SystemName: normalizedSystem,
		Status:     statusPending,
	}
	return u.repo.Create(req)
}

func (u *passwordRequestUsecase) GetEmployeeRequests(empID int64) ([]domain.PasswordRequest, error) {
	requests, err := u.repo.GetByEmployeeID(empID)
	if err != nil {
		return nil, err
	}
	return decryptApprovedPasswords(requests)
}

func (u *passwordRequestUsecase) GetAllRequests() ([]domain.PasswordRequest, error) {
	requests, err := u.repo.GetAll()
	if err != nil {
		return nil, err
	}
	return decryptApprovedPasswords(requests)
}

func (u *passwordRequestUsecase) ReviewRequest(id int64, status, reason string) error {
	status = strings.TrimSpace(strings.ToLower(status))
	reason = strings.TrimSpace(reason)

	if status != statusApproved && status != statusRejected {
		return errors.New("недопустимый статус заявки")
	}

	req, err := u.repo.GetByID(id)
	if err != nil {
		return err
	}
	if req.Status != statusPending {
		return errors.New("эта заявка уже обработана")
	}

	var encryptedPassword string
	if status == statusRejected {
		if reason == "" {
			return errors.New("при отклонении заявки необходимо указать причину")
		}
	} else {
		primaryPassword, err := primaryPasswordForRequest(req)
		if err != nil {
			return err
		}

		encryptedPassword, err = crypto.Encrypt(primaryPassword, os.Getenv("PASSWORD_ENCRYPTION_KEY"))
		if err != nil {
			return errors.New("ошибка шифрования пароля")
		}
		reason = ""
	}

	return u.repo.UpdateStatus(id, status, encryptedPassword, reason)
}

func normalizeSystemName(systemName string) (string, error) {
	key := strings.ToLower(strings.TrimSpace(systemName))
	key = strings.ReplaceAll(key, "ё", "е")

	if system, ok := supportedSystems[key]; ok {
		return system, nil
	}
	return "", errors.New("указанная система не поддерживается")
}

func isValidIIN(iin string) bool {
	matched, _ := regexp.MatchString(`^\d{12}$`, iin)
	return matched
}

func primaryPasswordForRequest(req *domain.PasswordRequest) (string, error) {
	switch req.SystemName {
	case "Documentolog":
		return requiredEnv("DEFAULT_PASSWORD_DOCUMENTOLOG")
	case "EPS":
		return requiredEnv("DEFAULT_PASSWORD_EPS")
	case "E-kyzmet":
		return requiredEnv("DEFAULT_PASSWORD_EKYZMET")
	case "AD":
		return requiredEnv("DEFAULT_PASSWORD_AD")
	case "E-otinish":
		if !isValidIIN(req.EmployeeIIN) {
			return "", errors.New("у сотрудника нет корректного ИИН для первичного пароля E-otinish")
		}
		return req.EmployeeIIN, nil
	default:
		return "", errors.New("указанная система не поддерживается")
	}
}

func requiredEnv(name string) (string, error) {
	value := strings.TrimSpace(os.Getenv(name))
	if value == "" {
		return "", fmt.Errorf("не настроена переменная окружения %s", name)
	}
	return value, nil
}

func decryptApprovedPasswords(requests []domain.PasswordRequest) ([]domain.PasswordRequest, error) {
	key := os.Getenv("PASSWORD_ENCRYPTION_KEY")
	for i := range requests {
		if requests[i].Status != statusApproved || requests[i].PrimaryPassword == "" {
			requests[i].PrimaryPassword = ""
			continue
		}

		password, err := crypto.Decrypt(requests[i].PrimaryPassword, key)
		if err != nil {
			return nil, errors.New("ошибка расшифровки первичного пароля")
		}
		requests[i].PrimaryPassword = password
	}
	return requests, nil
}
