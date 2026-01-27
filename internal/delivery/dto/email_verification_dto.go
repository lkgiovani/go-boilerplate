package dto

// Email Verification Request DTOs

type VerifyEmailRequest struct {
	Token string `json:"token" validate:"required"`
}

type ResendEmailVerificationRequest struct {
	Email string `json:"email" validate:"required,email"`
}

// Email Verification Response DTOs

type EmailVerificationResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
	UserID  int64  `json:"userId,omitempty"`
	Email   string `json:"email,omitempty"`
}
