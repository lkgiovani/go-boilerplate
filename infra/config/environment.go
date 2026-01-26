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
}

type DatabaseConfig struct {
	Host     string
	Port     int
	User     string
	Password string
	DBName   string
}

type ServerConfig struct {
	Port     int
	LogLevel string
	Mode     string
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

	return ServerConfig{
		Port:     port,
		LogLevel: logLevel,
		Mode:     mode,
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

	// Optional: Audience (default: "boilerplate-api")
	audience, _ := utils.GetString("JWT_AUDIENCE")
	if audience == "" {
		audience = "boilerplate-api"
	}

	// Optional: Cookie Domain (default: empty)
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

func LoadConfig() *Config {
	LoadEnvironment()
	return &Config{
		Database: loadDatabaseConfig(),
		Server:   loadServerConfig(),
		JWT:      loadJWTConfig(),
		Admin:    loadAdminConfig(),
	}
}
