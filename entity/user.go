package entity

import (
	"time"
)

type Account struct {
	ID             string     `json:"id"`
	PublicKey      string     `json:"public_key"`
	Name           string     `json:"name"`
	Email          string     `json:"email"`
	PhoneNumber    string     `json:"phone_number"`
	WebUrl         string     `json:"web_url"`
	FrontendUrl    string     `json:"frontend_url"`
	DisplayName    string     `json:"display_name"`
	Status         int        `json:"status"`
	Settings       string     `json:"settings"`
	Platform       string     `json:"platform"`
	AssetDirectory string     `json:"asset_directory"`
	BillingPlanId  int        `json:"billing_plan_id,omitempty"`
	TrialEndsAt    time.Time  `json:"trial_ends_at,omitempty"`
	CreatedAt      time.Time  `json:"created_at,omitempty"`
	UpdatedAt      time.Time  `json:"updated_at,omitempty"`
	DeletedAt      *time.Time `json:"deleted_at,omitempty"`
}
type UserDetails struct {
	ID          string     `json:"id"`
	AccountId   string     `json:"account_id"`
	Name        string     `json:"name"`
	Token       string     `json:"token"`
	Password    string     `json:"password"`
	Email       string     `json:"email"`
	Status      int        `json:"status"`
	CreatedAt   time.Time  `json:"created_at,omitempty"`
	LastLoginAt time.Time  `json:"last_login_at,omitempty"`
	UpdatedAt   time.Time  `json:"updated_at,omitempty"`
	DeletedAt   *time.Time `json:"deleted_at,omitempty"`
}

type UserLogin struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type ResetPasswordRequest struct {
	Email       string `json:"email"`
	NewPassword string `json:"new_password"`
}

type UserResponse struct {
	Message    string      `json:"message,omitempty"`
	Data       interface{} `json:"data,omitempty"`
	Status     bool        `json:"status,omitempty"`
	StatusCode int         `json:"statuscode,omitempty"`
	Error      interface{} `json:"error,omitempty"`
	Timestamp  time.Time   `json:"timestamp,omitempty"`
}
type Settings struct {
	AccessToken string `json:"access_token,omitempty"`
}
