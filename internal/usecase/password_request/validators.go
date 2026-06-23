package password_request

import (
	"encoding/json"
	"errors"
	"regexp"
	"strings"
)

type SystemValidator interface {
	Validate(rawData json.RawMessage) error
}

type ResetPasswordInput struct {
	IIN string `json:"iin"`
}

func validateIIN(rawData json.RawMessage, systemName string) error {
	var input ResetPasswordInput
	if err := json.Unmarshal(rawData, &input); err != nil {
		return errors.New("invalid input format for system")
	}

	input.IIN = strings.TrimSpace(input.IIN)
	matched, _ := regexp.MatchString(`^\d{12}$`, input.IIN)
	if !matched {
		return errors.New("iin must contain exactly 12 digits")
	}

	return nil
}

type EOtinishValidator struct{}

func (v EOtinishValidator) Validate(rawData json.RawMessage) error {
	return validateIIN(rawData, "E-otinish")
}

type DocumentologValidator struct{}

func (v DocumentologValidator) Validate(rawData json.RawMessage) error {
	return validateIIN(rawData, "ИС Documentolog")
}

type EpsMailValidator struct{}

func (v EpsMailValidator) Validate(rawData json.RawMessage) error {
	return validateIIN(rawData, "EPS")
}

type EKyzmetValidator struct{}

func (v EKyzmetValidator) Validate(rawData json.RawMessage) error {
	return validateIIN(rawData, "E-kyzmet")
}

type ActiveDirectoryValidator struct{}

func (v ActiveDirectoryValidator) Validate(rawData json.RawMessage) error {
	return validateIIN(rawData, "AD")
}
