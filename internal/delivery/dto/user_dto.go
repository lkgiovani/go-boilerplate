package dto

import (
	"time"
)

// Request DTOs

type SignupUserRequestDTO struct {
	Name     string `json:"name" validate:"required,min=3,max=255"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=6"`
}

type UserPostRequestDTO struct {
	Name     string  `json:"name" validate:"required,min=3,max=255"`
	Email    string  `json:"email" validate:"required,email"`
	Password *string `json:"password,omitempty" validate:"omitempty,min=6"`
	Admin    *bool   `json:"admin,omitempty"`
	Active   *bool   `json:"active,omitempty"`
	Source   *string `json:"source,omitempty"`
}

type UserPutRequestDTO struct {
	Name   string  `json:"name" validate:"required,min=3,max=255"`
	Email  string  `json:"email" validate:"required,email"`
	ImgURL *string `json:"imgUrl,omitempty"`
}

type UserPutPasswordRequestDTO struct {
	CurrentPassword *string `json:"currentPassword,omitempty"`
	Password        string  `json:"password" validate:"required,min=6"`
}

type UploadImageRequestDTO struct {
	FileName      string `json:"fileName" validate:"required"`
	ContentType   string `json:"contentType" validate:"required"`
	ContentLength int64  `json:"contentLength" validate:"required,min=1"`
}

type UserUpdateAccessModeDTO struct {
	AccessMode string `json:"accessMode" validate:"required,oneof=FREE BASIC PRO LIFETIME_PRO"`
}

type UserUpdateFeaturesDTO struct {
	CanCreateBudgets *bool `json:"canCreateBudgets,omitempty"`
	CanExportData    *bool `json:"canExportData,omitempty"`
	CanUseReports    *bool `json:"canUseReports,omitempty"`
	CanUseGoals      *bool `json:"canUseGoals,omitempty"`
}

type UserUpdateLimitsDTO struct {
	MaxAccounts             *int `json:"maxAccounts,omitempty"`
	MaxTransactionsPerMonth *int `json:"maxTransactionsPerMonth,omitempty"`
	MaxCategoriesPerAccount *int `json:"maxCategoriesPerAccount,omitempty"`
}

type UserGrantLifetimeProDTO struct {
	Reason string `json:"reason" validate:"required"`
}

type ResendVerificationRequest struct {
	Email string `json:"email" validate:"required,email"`
}

// Response DTOs

type UserResponseDTO struct {
	ID         int64      `json:"id"`
	Name       string     `json:"name"`
	Email      string     `json:"email"`
	ImgURL     *string    `json:"imgUrl,omitempty"`
	Admin      bool       `json:"admin"`
	Active     bool       `json:"active"`
	Source     string     `json:"source"`
	Metadata   *string    `json:"metadata,omitempty"`
	LastAccess *time.Time `json:"lastAccess,omitempty"`
	CreatedAt  time.Time  `json:"createdAt"`
	UpdatedAt  time.Time  `json:"updatedAt"`
}

type UploadResponseDTO struct {
	UploadSignedURL string `json:"uploadSignedUrl"`
	PublicURL       string `json:"publicUrl"`
}

type MessageResponse struct {
	Message string `json:"message"`
}

type PageResponse struct {
	Content       []UserResponseDTO `json:"content"`
	TotalElements int64             `json:"totalElements"`
	TotalPages    int               `json:"totalPages"`
	Size          int               `json:"size"`
	Number        int               `json:"number"`
	First         bool              `json:"first"`
	Last          bool              `json:"last"`
}
