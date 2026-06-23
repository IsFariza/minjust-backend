package password_request

import (
	"encoding/json"
	"errors"
	"minjust-website/internal/crypto"
	"minjust-website/internal/domain"
	"os"
	"strings"
)

type eOtinishPayload struct {
	IIN string `json:"iin"`
}

type passwordRequestUsecase struct {
	repo       domain.PasswordRequestRepository
	validators map[string]SystemValidator
	secretKey  string
}

func NewPasswordRequestUsecase(repo domain.PasswordRequestRepository, secretKey string) domain.PasswordRequestUsecase {
	return &passwordRequestUsecase{
		repo:      repo,
		secretKey: secretKey,
		validators: map[string]SystemValidator{
			"e.otinish": EOtinishValidator{},
		},
	}
}

func (u *passwordRequestUsecase) GetRequestStatus(id int64) (*domain.PasswordRequest, error) {
	req, err := u.repo.GetByID(id)
	if err != nil {
		return nil, err
	}

	if req.PrimaryPassword != "" {
		password, err := crypto.Decrypt(req.PrimaryPassword, u.secretKey)
		if err != nil {
			return nil, errors.New("failed to read primary password")
		}
		req.PrimaryPassword = password
	}

	return req, nil
}

func (u *passwordRequestUsecase) RequestResetPassword(req *domain.PasswordRequest) error {
	req.SystemName = strings.ToLower(strings.TrimSpace(req.SystemName))

	validator, exists := u.validators[req.SystemName]
	if !exists {
		return errors.New("selected system is not supported yet")
	}

	if err := validator.Validate(req.InputData); err != nil {
		return err
	}

	req.Status = domain.StatusPending
	req.PrimaryPassword = ""
	req.AdminComment = ""

	return u.repo.Create(req)
}

func (u *passwordRequestUsecase) ProcessRequest(id int64, status, password, comment string) error {
	req, err := u.repo.GetByID(id)
	if err != nil {
		return err
	}
	if req.Status != domain.StatusPending {
		return errors.New("this request has already been processed")
	}

	status = strings.ToLower(strings.TrimSpace(status))
	password = strings.TrimSpace(password)
	comment = strings.TrimSpace(comment)

	switch status {
	case domain.StatusApproved:
		if password == "" {
			password, err = defaultPrimaryPassword(req)
			if err != nil {
				return err
			}
		}

		encryptedPassword, err := crypto.Encrypt(password, u.secretKey)
		if err != nil {
			return errors.New("failed to encrypt primary password")
		}

		return u.repo.UpdateStatus(id, domain.StatusApproved, encryptedPassword, comment)

	case domain.StatusRejected:
		if comment == "" {
			return errors.New("rejection reason is required")
		}
		return u.repo.UpdateStatus(id, domain.StatusRejected, "", comment)

	default:
		return errors.New("invalid status, use approved or rejected")
	}
}

func defaultPrimaryPassword(req *domain.PasswordRequest) (string, error) {
	var payload ResetPasswordInput
	if err := json.Unmarshal(req.InputData, &payload); err != nil {
		return "", errors.New("failed to parse IIN from request data")
	}
	payload.IIN = strings.TrimSpace(payload.IIN)

	switch req.SystemName {
	case "e.otinish":
		return payload.IIN, nil

	case "documentolog":
		return os.Getenv("DEFAULT_PASSWORD_DOCUMENTOLOG"), nil

	case "eps":
		return os.Getenv("DEFAULT_PASSWORD_EPS"), nil

	case "e-kyzmet":
		return os.Getenv("DEFAULT_PASSWORD_EKYZMET"), nil

	case "ad":
		return os.Getenv("DEFAULT_PASSWORD_AD"), nil

	default:
		return os.Getenv("DEFAULT_PASSWORD"), nil
	}
}
