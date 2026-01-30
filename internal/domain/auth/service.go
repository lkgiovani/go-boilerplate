package auth

import (
	"context"
	"strings"
	"time"
	"unicode"

	"net/http"

	"github.com/google/uuid"
	"github.com/lkgiovani/go-boilerplate/internal/domain/emailverification"
	"github.com/lkgiovani/go-boilerplate/internal/domain/user"
	"github.com/lkgiovani/go-boilerplate/internal/errors"
	"github.com/lkgiovani/go-boilerplate/internal/security/jwt"
	"github.com/lkgiovani/go-boilerplate/pkg/encrypt"
	"github.com/lkgiovani/go-boilerplate/pkg/utils"
)

type Service struct {
	UserRepo                 user.UserService
	UserService              *user.Service
	AuthRepo                 Repository
	JwtService               *jwt.JwtService
	EmailVerificationService *emailverification.Service
	GoogleTokenGateway       GoogleTokenGateway
}

func NewService(
	userRepo user.UserService,
	userSvc *user.Service,
	authRepo Repository,
	jwtService *jwt.JwtService,
	emailVerSvc *emailverification.Service,
	googleTokenGateway GoogleTokenGateway,
) *Service {
	return &Service{
		UserRepo:                 userRepo,
		UserService:              userSvc,
		AuthRepo:                 authRepo,
		JwtService:               jwtService,
		EmailVerificationService: emailVerSvc,
		GoogleTokenGateway:       googleTokenGateway,
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

	if !u.Admin {
		if !u.Active {
			return nil, errors.Errorf(errors.EUNAUTHORIZED, "Sua conta está inativa. Entre em contato com o suporte.")
		}
		if !u.Metadata.EmailVerified {
			return nil, errors.Errorf(errors.EUNAUTHORIZED, "Email não verificado. Verifique seu email para acessar a conta.")
		}
	}

	return u, nil
}

func (s *Service) CreateSession(ctx context.Context, u *user.User, userAgent, ipAddress, deviceID string) (string, string, []*http.Cookie, error) {
	accessToken, refreshToken, refreshClaims, cookies, err := s.JwtService.GenerateCookies(u)
	if err != nil {
		return "", "", nil, err
	}

	rt := &RefreshToken{
		ID:        uuid.New(),
		UserID:    u.ID,
		UserEmail: u.Email,
		DeviceID:  deviceID,
		Jti:       refreshClaims.Jti,
		FamilyID:  uuid.New(),
		TokenHash: utils.HashToken(refreshToken),
		ExpiresAt: time.Unix(refreshClaims.ExpiresAt, 0),
		CreatedAt: utils.Now(),
		UserAgent: userAgent,
		IpAddress: ipAddress,
	}

	if err := s.AuthRepo.CreateRefreshToken(ctx, rt); err != nil {
		return "", "", nil, err
	}

	return accessToken, refreshToken, cookies, nil
}

func (s *Service) RefreshToken(ctx context.Context, token, userAgent, ipAddress, deviceID string) (string, string, []*http.Cookie, error) {

	claims, err := s.JwtService.ParseToken(token)
	if err != nil {
		return "", "", nil, errors.Errorf(errors.EUNAUTHORIZED, "invalid refresh token")
	}

	if claims.Type != "refresh" {
		return "", "", nil, errors.Errorf(errors.EUNAUTHORIZED, "invalid token type")
	}

	hash := utils.HashToken(token)
	storedToken, err := s.AuthRepo.GetRefreshTokenByHash(ctx, hash)
	if err != nil {
		return "", "", nil, errors.Errorf(errors.EUNAUTHORIZED, "refresh token not found or revoked")
	}

	if storedToken.RevokedAt != nil {

		_ = s.AuthRepo.RevokeAllUserRefreshTokens(ctx, storedToken.UserID)
		return "", "", nil, errors.Errorf(errors.EUNAUTHORIZED, "token revoked")
	}

	if storedToken.Used {

		_ = s.AuthRepo.RevokeAllUserRefreshTokens(ctx, storedToken.UserID)
		return "", "", nil, errors.Errorf(errors.EUNAUTHORIZED, "token already used")
	}

	u, err := s.UserRepo.GetByID(ctx, storedToken.UserID)
	if err != nil || u == nil {
		return "", "", nil, errors.Errorf(errors.EUNAUTHORIZED, "user not found")
	}

	if !u.Admin {
		if !u.Active {
			return "", "", nil, errors.Errorf(errors.EUNAUTHORIZED, "Sua conta está inativa.")
		}
		if !u.Metadata.EmailVerified {
			return "", "", nil, errors.Errorf(errors.EUNAUTHORIZED, "Email não verificado.")
		}
	}

	accessToken, refreshToken, refreshClaims, cookies, err := s.JwtService.GenerateCookies(u)
	if err != nil {
		return "", "", nil, err
	}

	if err := s.AuthRepo.MarkAsUsed(ctx, hash); err != nil {
		return "", "", nil, err
	}

	newRt := &RefreshToken{
		ID:          uuid.New(),
		UserID:      u.ID,
		UserEmail:   u.Email,
		DeviceID:    deviceID,
		Jti:         refreshClaims.Jti,
		FamilyID:    storedToken.FamilyID,
		TokenHash:   utils.HashToken(refreshToken),
		ExpiresAt:   time.Unix(refreshClaims.ExpiresAt, 0),
		CreatedAt:   utils.Now(),
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

	if _, err := s.EmailVerificationService.CreateAndSendVerificationToken(ctx, u); err != nil {

		return nil
	}

	return nil
}

func (s *Service) RevokeRefreshToken(ctx context.Context, token string) error {

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

func (s *Service) AuthenticateWithGoogleMobile(ctx context.Context, idToken, deviceID, userAgent, ipAddress string) (*MobileAuthResult, error) {
	googleUser, err := s.GoogleTokenGateway.VerifyAndExtract(ctx, idToken)
	if err != nil {
		return nil, errors.Errorf(errors.EUNAUTHORIZED, "falha ao verificar token do Google: %v", err)
	}

	userEntity, isNewUser, err := s.findOrCreateGoogleUser(ctx, googleUser)
	if err != nil {
		return nil, err
	}

	now := time.Now()
	userEntity.LastAccess = &now
	if err := s.UserRepo.Update(ctx, userEntity); err != nil {
		return nil, errors.Errorf(errors.EINTERNAL, "falha ao atualizar último acesso")
	}

	accessToken, refreshToken, _, err := s.CreateSession(ctx, userEntity, userAgent, ipAddress, deviceID)
	if err != nil {
		return nil, err
	}

	return &MobileAuthResult{
		UserID:       userEntity.ID,
		Email:        userEntity.Email,
		Name:         userEntity.Name,
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresIn:    s.JwtService.GetAccessTokenExpirationSeconds(),
		IsNewUser:    isNewUser,
	}, nil
}

func (s *Service) findOrCreateGoogleUser(ctx context.Context, googleUser *GoogleUserInfo) (*user.User, bool, error) {
	existingUser, err := s.UserRepo.GetByEmail(ctx, googleUser.Email)
	if err == nil && existingUser != nil {
		if existingUser.Source != "GOOGLE" {
			existingUser.Source = "GOOGLE"
			if googleUser.PictureURL != "" {
				existingUser.ImgURL = &googleUser.PictureURL
			}

			existingUser.Metadata.EmailVerified = true
			if err := s.UserRepo.Update(ctx, existingUser); err != nil {
				return nil, false, errors.Errorf(errors.EINTERNAL, "falha ao atualizar usuário")
			}
		}
		return existingUser, false, nil
	}

	newUser := &user.User{
		Name:     googleUser.Name,
		Email:    googleUser.Email,
		Source:   "GOOGLE",
		ImgURL:   &googleUser.PictureURL,
		Active:   true,
		Admin:    false,
		Metadata: user.NewDefaultMetadata(),
	}
	newUser.Metadata.EmailVerified = true

	if createErr := s.UserRepo.Create(ctx, newUser); createErr != nil {
		return nil, false, errors.Errorf(errors.EINTERNAL, "falha ao criar usuário")
	}

	return newUser, true, nil
}

func (s *Service) RefreshMobileToken(ctx context.Context, refreshToken, userAgent, ipAddress, deviceID string) (*MobileRefreshResult, error) {
	accessToken, newRefreshToken, _, err := s.RefreshToken(ctx, refreshToken, userAgent, ipAddress, deviceID)
	if err != nil {
		return nil, err
	}

	return &MobileRefreshResult{
		AccessToken:  accessToken,
		RefreshToken: newRefreshToken,
		ExpiresIn:    s.JwtService.GetAccessTokenExpirationSeconds(),
	}, nil
}
