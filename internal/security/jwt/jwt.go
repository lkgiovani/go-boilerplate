package jwt

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
	"github.com/lkgiovani/go-boilerplate/infra/config"
	"github.com/lkgiovani/go-boilerplate/internal/domain/user"
)

const (
	CookieName = "token"
)

type CustomClaims struct {
	ID    string   `json:"id"`
	Name  string   `json:"name"`
	Email string   `json:"email"`
	Roles []string `json:"roles,omitempty"`
	jwt.StandardClaims
}

type JwtService struct {
	secretKey    string
	issuer       string
	audience     string
	cookieDomain string
	tokenTTL     int64
	parser       *jwt.Parser
	userService  *user.Service
}

func NewJwtService(settings config.JWTConfig, userService *user.Service) (*JwtService, error) {
	if settings.ExpiresIn <= 0 {
		return nil, errors.New("JWT_EXPIRES_IN inválido")
	}

	audience := settings.Audience
	if audience == "" {
		audience = "boilerplate-api"
	}

	parser := &jwt.Parser{ValidMethods: []string{jwt.SigningMethodHS256.Name}}
	return &JwtService{
		secretKey:    settings.SecretKey,
		issuer:       settings.Issuer,
		audience:     audience,
		cookieDomain: settings.CookieDomain,
		tokenTTL:     int64(settings.ExpiresIn.Seconds()),
		parser:       parser,
		userService:  userService,
	}, nil
}

func (s *JwtService) GetTokenFromCookie(r *http.Request) (string, error) {
	cookie, err := r.Cookie(CookieName)
	if err != nil {
		return "", err
	}
	return cookie.Value, nil
}

func (s *JwtService) GetUserName(tokenString string) (string, error) {
	claims, err := s.parseCustomClaims(tokenString)
	if err != nil {
		return "", err
	}
	return claims.Subject, nil
}

func (s *JwtService) GenerateTokenFromUser(ctx context.Context, u *user.User) (string, error) {
	now := time.Now().Unix()

	// Set roles based on user's admin status
	roles := []string{"USER"}
	if u.Admin {
		roles = append(roles, "ADMIN")
	}

	claims := CustomClaims{
		ID:    strconv.FormatInt(u.ID, 10),
		Name:  u.Name,
		Email: u.Email,
		Roles: roles,
		StandardClaims: jwt.StandardClaims{
			Subject:   u.Email,
			Issuer:    s.issuer,
			Audience:  s.audience,
			IssuedAt:  now,
			ExpiresAt: now + s.tokenTTL,
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(s.secretKey))
}

func (s *JwtService) GenerateTokenFromEmail(email string) (string, error) {
	now := time.Now().Unix()

	claims := jwt.StandardClaims{
		Subject:   email,
		Issuer:    s.issuer,
		Audience:  s.audience,
		IssuedAt:  now,
		ExpiresAt: now + s.tokenTTL,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(s.secretKey))
}

func (s *JwtService) GenerateToken(ctx context.Context, id string) (string, error) {

	userID, err := uuid.Parse(id)
	if err != nil {
		return "", fmt.Errorf("invalid user ID: %w", err)
	}

	now := time.Now().Unix()
	claims := CustomClaims{
		ID: userID.String(),
		StandardClaims: jwt.StandardClaims{
			Subject:   id,
			Issuer:    s.issuer,
			Audience:  s.audience,
			IssuedAt:  now,
			ExpiresAt: now + s.tokenTTL,
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(s.secretKey))
}

func (s *JwtService) GenerateCookie(u *user.User, r *http.Request) (*http.Cookie, error) {
	tokenString, err := s.GenerateTokenFromUser(context.Background(), u)
	if err != nil {
		return nil, err
	}

	cookie := &http.Cookie{
		Name:     CookieName,
		Value:    tokenString,
		Path:     "/",
		MaxAge:   216000,
		Secure:   true,
		HttpOnly: true,
		SameSite: http.SameSiteNoneMode,
	}

	if s.cookieDomain != "" {
		cookie.Domain = s.cookieDomain
	}

	return cookie, nil
}

func (s *JwtService) CleanCookie() *http.Cookie {
	cookie := &http.Cookie{
		Name:     CookieName,
		Value:    "",
		Path:     "/",
		MaxAge:   -1,
		Secure:   true,
		HttpOnly: true,
		SameSite: http.SameSiteNoneMode,
	}

	if s.cookieDomain != "" {
		cookie.Domain = s.cookieDomain
	}

	return cookie
}

func (s *JwtService) ValidateToken(tokenString string) bool {
	claims, err := s.parseCustomClaims(tokenString)
	if err != nil {
		s.logValidationError(err)
		return false
	}

	now := time.Now().Unix()

	if !claims.VerifyIssuer(s.issuer, true) {
		log.Printf("Invalid JWT issuer: expected %s", s.issuer)
		return false
	}

	if !claims.VerifyAudience(s.audience, true) {
		log.Printf("Invalid JWT audience: expected %s", s.audience)
		return false
	}

	if !claims.VerifyExpiresAt(now, true) {
		log.Printf("JWT token is expired")
		return false
	}

	return true
}

func (s *JwtService) ParseToken(tokenString string) (*CustomClaims, error) {
	claims, err := s.parseCustomClaims(tokenString)
	if err != nil {
		s.logValidationError(err)
		return nil, err
	}

	now := time.Now().Unix()

	if !claims.VerifyIssuer(s.issuer, true) {
		return nil, fmt.Errorf("invalid issuer")
	}

	if !claims.VerifyAudience(s.audience, true) {
		return nil, fmt.Errorf("invalid audience")
	}

	if !claims.VerifyExpiresAt(now, true) {
		return nil, fmt.Errorf("token expired")
	}

	return claims, nil
}

func (s *JwtService) GetAccessTokenExpirationSeconds() int64 {
	return s.tokenTTL
}

func (s *JwtService) parseCustomClaims(tokenString string) (*CustomClaims, error) {
	token, err := s.parser.ParseWithClaims(tokenString, &CustomClaims{}, func(token *jwt.Token) (interface{}, error) {

		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(s.secretKey), nil
	})

	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*CustomClaims)
	if !ok || !token.Valid {
		return nil, fmt.Errorf("invalid token claims")
	}

	return claims, nil
}

func (s *JwtService) logValidationError(err error) {
	if err == nil {
		return
	}

	var ve *jwt.ValidationError
	if errors.As(err, &ve) {
		if ve.Errors&jwt.ValidationErrorMalformed != 0 {
			log.Printf("Invalid JWT token: malformed token")
		} else if ve.Errors&jwt.ValidationErrorExpired != 0 {
			log.Printf("JWT token is expired: %v", err)
		} else if ve.Errors&jwt.ValidationErrorNotValidYet != 0 {
			log.Printf("JWT token is not valid yet: %v", err)
		} else if ve.Errors&jwt.ValidationErrorSignatureInvalid != 0 {
			log.Printf("Invalid JWT signature: %v", err)
		} else {
			log.Printf("JWT validation error: %v", err)
		}
	} else {
		log.Printf("JWT validation error: %v", err)
	}
}
