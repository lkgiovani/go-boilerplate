package dto

// Auth Request DTOs

type LoginRequestDTO struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

// Auth Response DTOs

type LoginResponseDTO struct {
	UserID      int64  `json:"userId"`
	Email       string `json:"email"`
	AccessToken string `json:"accessToken"`
	ExpiresIn   int64  `json:"expiresIn"`
}

type RefreshResponseDTO struct {
	AccessToken string `json:"accessToken"`
	ExpiresIn   int    `json:"expiresIn"`
}
