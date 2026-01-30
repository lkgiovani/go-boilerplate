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
type MobileOAuth2RequestDTO struct {
	IdToken  string `json:"idToken" validate:"required"`
	DeviceId string `json:"deviceId"`
}

type MobileLoginResponseDTO struct {
	UserID       int64  `json:"userId"`
	Email        string `json:"email"`
	Name         string `json:"name"`
	AccessToken  string `json:"accessToken"`
	RefreshToken string `json:"refreshToken"`
	ExpiresIn    int64  `json:"expiresIn"`
	IsNewUser    bool   `json:"isNewUser"`
}

type MobileRefreshRequestDTO struct {
	RefreshToken string `json:"refreshToken" validate:"required"`
}

type MobileRefreshResponseDTO struct {
	AccessToken  string `json:"accessToken"`
	RefreshToken string `json:"refreshToken"`
	ExpiresIn    int64  `json:"expiresIn"`
}
