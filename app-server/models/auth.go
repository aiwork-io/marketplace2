package models

import "time"

const (
	AUTH_SCENARIO_ACCOUNT_CONFIRMATION string = "account_confirmation"
	AUTH_SCENARIO_PASSWORD_RESET       string = "password_reset"
)

type AuthState struct {
	UserId    string `json:"user_id"`
	Scenario  string `json:"scenario"`
	ExpiredAt int64  `json:"expired_at"`
}

func (state *AuthState) Valid() bool {
	return state.ExpiredAt >= time.Now().UTC().Unix()
}

func NewAuthState(id, scenario string) *AuthState {
	return &AuthState{
		UserId:    id,
		Scenario:  scenario,
		ExpiredAt: time.Now().UTC().Add(24 * time.Hour).Unix(),
	}
}
