package config

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
	go_boilerplate "github.com/lkgiovani/go-boilerplate"

	"github.com/lkgiovani/go-boilerplate/pkg/utils"
)

var embeddedEnv = go_boilerplate.EnvFile

type Config struct {
	Database DatabaseConfig
	Server   ServerConfig
	JWT      JWTConfig
	Admin    AdminConfig
	Email    EmailConfig
	Storage  StorageConfig
	OAuth2   OAuth2Config
}

type StorageConfig struct {
	Provider             string
	LocalDir             string
	PresignedUrlDuration int
	S3AccessKey          string
	S3SecretKey          string
	S3Region             string
	S3BucketName         string
	S3Endpoint           string
	R2AccountID          string
	R2AccessKey          string
	R2SecretKey          string
	R2BucketName         string
	R2PublicURL          string
	PublicBaseURL        string
}

type EmailConfig struct {
	Provider     string
	FromEmail    string
	FromName     string
	SMTPHost     string
	SMTPPort     int
	SMTPUser     string
	SMTPPassword string
	APIKey       string
	SESAccessKey string
	SESSecretKey string
	SESRegion    string
	SESEndpoint  string
	FrontendURL  string
}

type DatabaseConfig struct {
	Host     string
	Port     int
	User     string
	Password string
	DBName   string
}

type ServerConfig struct {
	Port           int
	LogLevel       string
	Mode           string
	AllowedOrigins string
}

type JWTConfig struct {
	SecretKey    string
	Issuer       string
	Audience     string
	CookieDomain string
	ExpiresIn    time.Duration
}

type AdminConfig struct {
	Email    string
	Password string
}

type OAuth2Config struct {
	GoogleAndroidClientID string
	GoogleIosClientID     string
}

func LoadEnvironment() {

	err := godotenv.Load()
	if err == nil {
		fmt.Println("✓ .env file loaded from current directory")
		return
	}

	fmt.Println("⚠ .env file not found. Trying embedded configuration...")
	envMap, parseErr := godotenv.Unmarshal(embeddedEnv)
	if parseErr != nil {
		fmt.Printf("❌ Error: Failed to parse embedded configuration: %v\n", parseErr)
		return
	}

	for key, value := range envMap {
		if os.Getenv(key) == "" {
			os.Setenv(key, value)
		}
	}
	fmt.Println("✓ Embedded configuration applied")
}

func loadServerConfig() ServerConfig {
	port, err := utils.GetInt("PORT")
	if err != nil {
		log.Fatalf("Failed to get PORT from environment: %v", err)
	}

	logLevel, err := utils.GetString("LOG_LEVEL")
	if err != nil {
		log.Fatalf("Failed to get LOG_LEVEL from environment: %v", err)
	}

	mode, err := utils.GetString("APP_MODE")
	if err != nil {
		log.Fatalf("Failed to get APP_MODE from environment: %v", err)
	}

	allowedOrigins, _ := utils.GetString("ALLOWED_ORIGINS")
	if allowedOrigins == "" {
		allowedOrigins = ""
	}

	return ServerConfig{
		Port:           port,
		LogLevel:       logLevel,
		Mode:           mode,
		AllowedOrigins: allowedOrigins,
	}
}

func loadDatabaseConfig() DatabaseConfig {
	host, err := utils.GetString("DB_HOST")
	if err != nil {
		log.Fatalf("Failed to get DB_HOST from environment: %v", err)
	}

	port, err := utils.GetInt("DB_PORT")
	if err != nil {
		log.Fatalf("Failed to get DB_PORT from environment: %v", err)
	}

	user, err := utils.GetString("DB_USER")
	if err != nil {
		log.Fatalf("Failed to get DB_USER from environment: %v", err)
	}

	password, err := utils.GetString("DB_PASSWORD")
	if err != nil {
		log.Fatalf("Failed to get DB_PASSWORD from environment: %v", err)
	}

	dbName, err := utils.GetString("DB_NAME")
	if err != nil {
		log.Fatalf("Failed to get DB_NAME from environment: %v", err)
	}

	return DatabaseConfig{
		Host:     host,
		Port:     port,
		User:     user,
		Password: password,
		DBName:   dbName,
	}
}

func loadJWTConfig() JWTConfig {
	secretKey, err := utils.GetString("JWT_SECRET_KEY")
	if err != nil {
		log.Fatalf("Failed to get JWT_SECRET_KEY from environment: %v", err)
	}

	issuer, err := utils.GetString("JWT_ISSUER")
	if err != nil {
		log.Fatalf("Failed to get JWT_ISSUER from environment: %v", err)
	}

	expiresIn, err := utils.GetDuration("JWT_EXPIRES_IN")
	if err != nil {
		log.Fatalf("Failed to get JWT_EXPIRES_IN from environment: %v", err)
	}

	audience, _ := utils.GetString("JWT_AUDIENCE")
	if audience == "" {
		audience = "boilerplate-api"
	}

	cookieDomain, _ := utils.GetString("COOKIE_DOMAIN")

	return JWTConfig{
		SecretKey:    secretKey,
		Issuer:       issuer,
		Audience:     audience,
		CookieDomain: cookieDomain,
		ExpiresIn:    expiresIn,
	}
}

func loadAdminConfig() AdminConfig {
	email, err := utils.GetString("ADMIN_EMAIL")
	if err != nil {
		log.Println("ADMIN_EMAIL not set, using default: admin@boilerplate.com")
		email = "admin@boilerplate.com"
	}

	password, err := utils.GetString("ADMIN_PASSWORD")
	if err != nil {
		log.Println("ADMIN_PASSWORD not set, using default: admin123")
		password = "admin123"
	}

	return AdminConfig{
		Email:    email,
		Password: password,
	}
}

func loadEmailConfig() EmailConfig {
	provider, _ := utils.GetString("EMAIL_PROVIDER")
	fromEmail, _ := utils.GetString("EMAIL_FROM")
	fromName, _ := utils.GetString("EMAIL_FROM_NAME")
	smtpHost, _ := utils.GetString("EMAIL_SMTP_HOST")
	smtpPort, _ := utils.GetInt("EMAIL_SMTP_PORT")
	smtpUser, _ := utils.GetString("EMAIL_SMTP_USER")
	smtpPassword, _ := utils.GetString("EMAIL_SMTP_PASSWORD")
	apiKey, _ := utils.GetString("EMAIL_API_KEY")
	sesAccessKey, _ := utils.GetString("AWS_SES_ACCESS_KEY_ID")
	sesSecretKey, _ := utils.GetString("AWS_SES_SECRET_ACCESS_KEY")
	sesRegion, _ := utils.GetString("AWS_SES_REGION")
	sesEndpoint, _ := utils.GetString("AWS_SES_ENDPOINT")
	frontendURL, _ := utils.GetString("FRONTEND_URL")

	return EmailConfig{
		Provider:     provider,
		FromEmail:    fromEmail,
		FromName:     fromName,
		SMTPHost:     smtpHost,
		SMTPPort:     smtpPort,
		SMTPUser:     smtpUser,
		SMTPPassword: smtpPassword,
		APIKey:       apiKey,
		SESAccessKey: sesAccessKey,
		SESSecretKey: sesSecretKey,
		SESRegion:    sesRegion,
		SESEndpoint:  sesEndpoint,
		FrontendURL:  frontendURL,
	}
}

func loadStorageConfig() StorageConfig {
	provider, _ := utils.GetString("STORAGE_PROVIDER")
	localDir, _ := utils.GetString("STORAGE_LOCAL_DIR")
	duration, _ := utils.GetInt("STORAGE_PRESIGNED_URL_DURATION")
	s3AccessKey, _ := utils.GetString("AWS_S3_ACCESS_KEY_ID")
	s3SecretKey, _ := utils.GetString("AWS_S3_SECRET_ACCESS_KEY")
	s3Region, _ := utils.GetString("AWS_S3_REGION")
	s3BucketName, _ := utils.GetString("AWS_S3_BUCKET_NAME")
	s3Endpoint, _ := utils.GetString("AWS_S3_ENDPOINT")
	r2AccountID, _ := utils.GetString("R2_ACCOUNT_ID")
	r2AccessKey, _ := utils.GetString("R2_ACCESS_KEY_ID")
	r2SecretKey, _ := utils.GetString("R2_SECRET_ACCESS_KEY")
	r2BucketName, _ := utils.GetString("R2_BUCKET_NAME")
	r2PublicURL, _ := utils.GetString("R2_PUBLIC_URL")
	publicBaseURL, _ := utils.GetString("STORAGE_PUBLIC_BASE_URL")

	if duration == 0 {
		duration = 60
	}

	return StorageConfig{
		Provider:             provider,
		LocalDir:             localDir,
		PresignedUrlDuration: duration,
		S3AccessKey:          s3AccessKey,
		S3SecretKey:          s3SecretKey,
		S3Region:             s3Region,
		S3BucketName:         s3BucketName,
		S3Endpoint:           s3Endpoint,
		R2AccountID:          r2AccountID,
		R2AccessKey:          r2AccessKey,
		R2SecretKey:          r2SecretKey,
		R2BucketName:         r2BucketName,
		R2PublicURL:          r2PublicURL,
		PublicBaseURL:        publicBaseURL,
	}
}

func loadOAuth2Config() OAuth2Config {
	androidID, err := utils.GetString("GOOGLE_ANDROID_CLIENT_ID")
	if androidID == "" {
		log.Fatalf("Failed to get GOOGLE_ANDROID_CLIENT_ID from environment: %v", err)
	}

	iosID, err := utils.GetString("GOOGLE_IOS_CLIENT_ID")
	if iosID == "" {
		log.Fatalf("Failed to get GOOGLE_IOS_CLIENT_ID from environment: %v", err)
	}

	return OAuth2Config{
		GoogleAndroidClientID: androidID,
		GoogleIosClientID:     iosID,
	}
}

func LoadConfig() *Config {
	LoadEnvironment()
	return &Config{
		Database: loadDatabaseConfig(),
		Server:   loadServerConfig(),
		JWT:      loadJWTConfig(),
		Admin:    loadAdminConfig(),
		Email:    loadEmailConfig(),
		Storage:  loadStorageConfig(),
		OAuth2:   loadOAuth2Config(),
	}
}
