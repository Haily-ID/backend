package tasks

import (
	"encoding/json"
	"fmt"

	"github.com/hibiken/asynq"
)

const (
	TypeSendOTPEmail = "email:send_otp"
)

type SendOTPEmailPayload struct {
	To      string `json:"to"`
	Name    string `json:"name"`
	OTP     string `json:"otp"`
	Purpose string `json:"purpose"`
}

func NewSendOTPEmailTask(to, name, otp, purpose string) (*asynq.Task, error) {
	payload, err := json.Marshal(SendOTPEmailPayload{
		To:      to,
		Name:    name,
		OTP:     otp,
		Purpose: purpose,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to marshal payload: %w", err)
	}
	return asynq.NewTask(TypeSendOTPEmail, payload), nil
}
