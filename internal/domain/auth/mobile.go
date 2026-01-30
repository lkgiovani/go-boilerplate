package auth

import (
	"context"
)

type GoogleUserInfo struct {
	Email      string
	Name       string
	PictureURL string
}

type MobileAuthResult struct {
	UserID       int64
	Email        string
	Name         string
	AccessToken  string
	RefreshToken string
	ExpiresIn    int64
	IsNewUser    bool
}

type MobileRefreshResult struct {
	AccessToken  string
	RefreshToken string
	ExpiresIn    int64
}

type GoogleTokenGateway interface {
	VerifyAndExtract(ctx context.Context, idToken string) (*GoogleUserInfo, error)
}
