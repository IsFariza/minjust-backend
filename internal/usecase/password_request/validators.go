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

type EOtinishInput struct {
	IIN string `json:"iin"`
}

type EOtinishValidator struct{}

func (v EOtinishValidator) Validate(rawData json.RawMessage) error {
	var input EOtinishInput
	if err := json.Unmarshal(rawData, &input); err != nil {
		return errors.New("invalid input format for e.otinish")
	}

	input.IIN = strings.TrimSpace(input.IIN)
	matched, _ := regexp.MatchString(`^\d{12}$`, input.IIN)
	if !matched {
		return errors.New("iin must contain exactly 12 digits")
	}

	return nil
}
