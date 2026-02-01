package jwt

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
	"github.com/lkgiovani/go-boilerplate/infra/config"
	"github.com/lkgiovani/go-boilerplate/internal/domain/user"
	"github.com/lkgiovani/go-boilerplate/pkg/utils"
)

const (
	AccessTokenCookieName  = "access_token"
	RefreshTokenCookieName = "refresh_token"
)

type CustomClaims struct {
	ID         string   `json:"id"`
	Name       string   `json:"name"`
	Email      string   `json:"email"`
	Roles      []string `json:"roles,omitempty"`
	Plan       string   `json:"plan,omitempty"`
	AccessMode string   `json:"access_mode,omitempty"`
	Jti        string   `json:"jti,omitempty"`
	Type       string   `json:"type,omitempty"`
	jwt.StandardClaims
}

type JwtService struct {
	secretKey                string
	issuer                   string
	audience                 string
	cookieDomain             string
	tokenTTL                 int64
	accessTokenCookieMaxAge  int
	refreshTokenCookieMaxAge int
	parser                   *jwt.Parser
	userService              *user.Service
}

func NewJwtService(settings config.JWTConfig, userService *user.Service) (*JwtService, error) {
	if settings.ExpirationMs <= 0 {
		return nil, errors.New("JWT_EXPIRATION_MS inválido")
	}

	audience := settings.Audience
	if audience == "" {
		audience = "boilerplate-api"
	}

	parser := &jwt.Parser{ValidMethods: []string{jwt.SigningMethodHS256.Name}}
	return &JwtService{
		secretKey:                settings.SecretKey,
		issuer:                   settings.Issuer,
		audience:                 audience,
		cookieDomain:             settings.CookieDomain,
		tokenTTL:                 int64(settings.ExpirationMs / 1000),
		accessTokenCookieMaxAge:  settings.AccessTokenCookieMaxAge,
		refreshTokenCookieMaxAge: settings.RefreshTokenCookieMaxAge,
		parser:                   parser,
		userService:              userService,
	}, nil
}

func (s *JwtService) GetTokenFromCookie(r *http.Request) (string, error) {
	cookie, err := r.Cookie(AccessTokenCookieName)
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
	token, _, err := s.GenerateAccessToken(u)
	return token, err
}

func (s *JwtService) GenerateAccessToken(u *user.User) (string, *CustomClaims, error) {
	return s.generateToken(u, s.tokenTTL, "access")
}

func (s *JwtService) GenerateRefreshToken(u *user.User) (string, *CustomClaims, error) {
	return s.generateToken(u, s.tokenTTL*7, "refresh")
}

func (s *JwtService) generateToken(u *user.User, ttl int64, tokenType string) (string, *CustomClaims, error) {
	now := utils.Now().Unix()

	roles := []string{"USER"}
	if u.Admin {
		roles = append(roles, "ADMIN")
	}

	claims := CustomClaims{
		ID:         strconv.FormatInt(u.ID, 10),
		Name:       u.Name,
		Email:      u.Email,
		Roles:      roles,
		Plan:       string(u.Metadata.PlanType),
		AccessMode: string(u.Metadata.AccessMode),
		Jti:        uuid.New().String(),
		Type:       tokenType,
		StandardClaims: jwt.StandardClaims{
			Subject:   u.Email,
			Issuer:    s.issuer,
			Audience:  s.audience,
			IssuedAt:  now,
			ExpiresAt: now + ttl,
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signed, err := token.SignedString([]byte(s.secretKey))
	if err != nil {
		return "", nil, err
	}
	return signed, &claims, nil
}

func (s *JwtService) GenerateTokenFromEmail(email string) (string, error) {
	now := utils.Now().Unix()

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

	now := utils.Now().Unix()
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

func (s *JwtService) GenerateCookies(u *user.User) (string, string, *CustomClaims, []*http.Cookie, error) {
	accessToken, _, err := s.GenerateAccessToken(u)
	if err != nil {
		return "", "", nil, nil, err
	}

	refreshToken, refreshClaims, err := s.GenerateRefreshToken(u)
	if err != nil {
		return "", "", nil, nil, err
	}

	isSecure := s.cookieDomain != ""

	cookies := []*http.Cookie{
		{
			Name:     AccessTokenCookieName,
			Value:    accessToken,
			Path:     "/",
			MaxAge:   s.accessTokenCookieMaxAge,
			Secure:   isSecure,
			HttpOnly: true,
			SameSite: http.SameSiteLaxMode,
		},
		{
			Name:     RefreshTokenCookieName,
			Value:    refreshToken,
			Path:     "/v1/auth/refresh",
			MaxAge:   s.refreshTokenCookieMaxAge,
			Secure:   isSecure,
			HttpOnly: true,
			SameSite: http.SameSiteLaxMode,
		},
	}

	if s.cookieDomain != "" {
		for _, c := range cookies {
			c.Domain = s.cookieDomain
		}
	}

	return accessToken, refreshToken, refreshClaims, cookies, nil
}

func (s *JwtService) GenerateCookie(u *user.User, r *http.Request) (*http.Cookie, error) {
	_, _, _, cookies, err := s.GenerateCookies(u)
	if err != nil {
		return nil, err
	}
	return cookies[0], nil
}

func (s *JwtService) CleanCookies() []*http.Cookie {
	names := []string{AccessTokenCookieName, RefreshTokenCookieName}
	return s.CleanAll(names)
}

func (s *JwtService) CleanAll(cookieNames []string) []*http.Cookie {
	cookies := make([]*http.Cookie, 0, len(cookieNames))
	for _, name := range cookieNames {
		if name == RefreshTokenCookieName {

			cookies = append(cookies, s.makeCleanCookie(name, "/v1/auth/refresh"))
		} else {

			cookies = append(cookies, s.makeCleanCookie(name, "/"))
		}
	}
	return cookies
}

func (s *JwtService) makeCleanCookie(name, path string) *http.Cookie {
	isSecure := s.cookieDomain != ""
	cookie := &http.Cookie{
		Name:     name,
		Value:    "",
		Path:     path,
		MaxAge:   -1,
		Secure:   isSecure,
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
	}
	if s.cookieDomain != "" {
		cookie.Domain = s.cookieDomain
	}
	return cookie
}

func (s *JwtService) CleanAllFromHeader(cookieHeader string) []*http.Cookie {
	var names []string
	if cookieHeader != "" {
		cookies := strings.Split(cookieHeader, ";")
		for _, cookie := range cookies {
			parts := strings.SplitN(cookie, "=", 2)
			if len(parts) > 0 {
				name := strings.TrimSpace(parts[0])
				if name != "" {
					names = append(names, name)
				}
			}
		}
	}

	names = append(names, AccessTokenCookieName, RefreshTokenCookieName)
	return s.CleanAll(names)
}

func (s *JwtService) ValidateToken(tokenString string) bool {
	claims, err := s.parseCustomClaims(tokenString)
	if err != nil {
		s.logValidationError(err)
		return false
	}

	now := utils.Now().Unix()

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

	now := utils.Now().Unix()

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

func (s *JwtService) GetRefreshTokenExpirationSeconds() int64 {
	return s.tokenTTL * 7
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
