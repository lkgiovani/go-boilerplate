package auth

import (
	"context"
	"strings"
	"unicode"

	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/lkgiovani/go-boilerplate/internal/domain/user"
	"github.com/lkgiovani/go-boilerplate/internal/errors"
	"github.com/lkgiovani/go-boilerplate/internal/security/jwt"
	"github.com/lkgiovani/go-boilerplate/pkg/encrypt"
	"github.com/lkgiovani/go-boilerplate/pkg/utils"
)

type Service struct {
	UserRepo    user.UserService
	UserService *user.Service
	AuthRepo    Repository
	JwtService  *jwt.JwtService
}

func NewService(
	userRepo user.UserService,
	userSvc *user.Service,
	authRepo Repository,
	jwtService *jwt.JwtService,
) *Service {
	return &Service{
		UserRepo:    userRepo,
		UserService: userSvc,
		AuthRepo:    authRepo,
		JwtService:  jwtService,
	}
}

func (s *Service) Login(ctx context.Context, login *Login) (*user.User, error) {
	u, err := s.UserRepo.GetByEmail(ctx, login.Email)
	if err != nil {
		return nil, errors.Errorf(errors.EUNAUTHORIZED, "invalid email or password")
	}

	if u.Password == nil {
		return nil, errors.Errorf(errors.EUNAUTHORIZED, "invalid email or password")
	}

	if err := encrypt.VerifyPassword(login.Password, *u.Password); err != nil {
		return nil, errors.Errorf(errors.EUNAUTHORIZED, "invalid email or password")
	}

	return u, nil
}

func (s *Service) CreateSession(ctx context.Context, u *user.User, userAgent, ipAddress, deviceID string) (string, string, []*http.Cookie, error) {
	accessToken, refreshToken, refreshClaims, cookies, err := s.JwtService.GenerateCookies(u)
	if err != nil {
		return "", "", nil, err
	}

	// Save refresh token to DB
	rt := &RefreshToken{
		ID:        uuid.New(),
		UserID:    u.ID,
		UserEmail: u.Email,
		DeviceID:  deviceID,
		Jti:       refreshClaims.Jti,
		FamilyID:  uuid.New(), // New family for new login
		TokenHash: utils.HashToken(refreshToken),
		ExpiresAt: time.Unix(refreshClaims.ExpiresAt, 0),
		CreatedAt: time.Now(),
		UserAgent: userAgent,
		IpAddress: ipAddress,
	}

	if err := s.AuthRepo.CreateRefreshToken(ctx, rt); err != nil {
		return "", "", nil, err
	}

	return accessToken, refreshToken, cookies, nil
}

func (s *Service) RefreshToken(ctx context.Context, token, userAgent, ipAddress, deviceID string) (string, string, []*http.Cookie, error) {
	// 1. Parse and validate the token string
	claims, err := s.JwtService.ParseToken(token)
	if err != nil {
		return "", "", nil, errors.Errorf(errors.EUNAUTHORIZED, "invalid refresh token")
	}

	if claims.Type != "refresh" {
		return "", "", nil, errors.Errorf(errors.EUNAUTHORIZED, "invalid token type")
	}

	// 2. Look up in DB
	hash := utils.HashToken(token)
	storedToken, err := s.AuthRepo.GetRefreshTokenByHash(ctx, hash)
	if err != nil {
		return "", "", nil, errors.Errorf(errors.EUNAUTHORIZED, "refresh token not found or revoked")
	}

	if storedToken.RevokedAt != nil {
		// Token reuse detection! Revoke the whole family
		_ = s.AuthRepo.RevokeAllUserRefreshTokens(ctx, storedToken.UserID)
		return "", "", nil, errors.Errorf(errors.EUNAUTHORIZED, "token revoked")
	}

	if storedToken.Used {
		// Potential reuse attack
		_ = s.AuthRepo.RevokeAllUserRefreshTokens(ctx, storedToken.UserID)
		return "", "", nil, errors.Errorf(errors.EUNAUTHORIZED, "token already used")
	}

	// 3. Get user
	u, err := s.UserRepo.GetByID(ctx, storedToken.UserID)
	if err != nil {
		return "", "", nil, errors.Errorf(errors.EUNAUTHORIZED, "user not found")
	}

	// 4. Generate new tokens (rotation)
	accessToken, refreshToken, refreshClaims, cookies, err := s.JwtService.GenerateCookies(u)
	if err != nil {
		return "", "", nil, err
	}

	// 5. Mark old as used
	if err := s.AuthRepo.MarkAsUsed(ctx, hash); err != nil {
		return "", "", nil, err
	}

	// 6. Save new token
	newRt := &RefreshToken{
		ID:          uuid.New(),
		UserID:      u.ID,
		UserEmail:   u.Email,
		DeviceID:    deviceID,
		Jti:         refreshClaims.Jti,
		FamilyID:    storedToken.FamilyID, // Keep the same family
		TokenHash:   utils.HashToken(refreshToken),
		ExpiresAt:   time.Unix(refreshClaims.ExpiresAt, 0),
		CreatedAt:   time.Now(),
		RotatedFrom: &storedToken.ID,
		UserAgent:   userAgent,
		IpAddress:   ipAddress,
	}

	if err := s.AuthRepo.CreateRefreshToken(ctx, newRt); err != nil {
		return "", "", nil, err
	}

	return accessToken, refreshToken, cookies, nil
}

func (s *Service) Register(ctx context.Context, u *user.User) error {
	exists, _ := s.UserRepo.GetByEmail(ctx, u.Email)
	if exists != nil {
		return errors.Errorf(errors.EDUPLICATION, "user already exists")
	}

	if u.Password == nil || *u.Password == "" {
		return errors.Errorf(errors.EINVALID, "password is required")
	}

	if err := PasswordRequirements(*u.Password); err != nil {
		return err
	}

	hashedPassword, err := encrypt.HashPassword(*u.Password)
	if err != nil {
		return errors.Errorf(errors.EINTERNAL, "failed to hash password")
	}
	u.Password = &hashedPassword

	if err := s.UserRepo.Create(ctx, u); err != nil {
		return errors.Errorf(errors.EINTERNAL, "failed to create user")
	}

	return nil
}

func (s *Service) RevokeRefreshToken(ctx context.Context, token string) error {
	// If token is empty, nothing to revoke
	if token == "" {
		return nil
	}
	hash := utils.HashToken(token)
	return s.AuthRepo.RevokeRefreshToken(ctx, hash)
}

func (s *Service) RevokeAllRefreshTokens(ctx context.Context, userID int64, currentToken string) error {
	if currentToken == "" {
		return s.AuthRepo.RevokeAllUserRefreshTokens(ctx, userID)
	}
	hash := utils.HashToken(currentToken)
	return s.AuthRepo.RevokeAllUserRefreshTokensExcept(ctx, userID, hash)
}

func PasswordRequirements(password string) error {
	if len(password) < 8 {
		return errors.Errorf(errors.EINVALID, "password must be at least 8 characters long")
	}

	var hasUpper, hasSpecial bool
	for _, char := range password {
		switch {
		case unicode.IsUpper(char):
			hasUpper = true
		case strings.ContainsRune("@$!%*?&", char):
			hasSpecial = true
		}
	}

	if !hasUpper {
		return errors.Errorf(errors.EINVALID, "password must contain at least one uppercase letter")
	}
	if !hasSpecial {
		return errors.Errorf(errors.EINVALID, "password must contain at least one special character")
	}

	return nil
}
