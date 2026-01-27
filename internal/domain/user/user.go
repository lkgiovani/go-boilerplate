package user

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"time"
)

type AccessMode string

const (
	AccessModeReadWrite AccessMode = "READ_WRITE"
	AccessModeReadOnly  AccessMode = "READ_ONLY"
	AccessModeDisabled  AccessMode = "DISABLED"
)

type PlanType string

const (
	PlanTypeFree       PlanType = "FREE"
	PlanTypePro        PlanType = "PRO"
	PlanTypeEnterprise PlanType = "ENTERPRISE"
)

type ReputationStatus string

const (
	ReputationStatusGood       ReputationStatus = "GOOD"
	ReputationStatusSuspicious ReputationStatus = "SUSPICIOUS"
	ReputationStatusBlocked    ReputationStatus = "BLOCKED"
)

type ProSource string

const (
	ProSourceTrial        ProSource = "TRIAL"
	ProSourceSubscription ProSource = "SUBSCRIPTION"
	ProSourceAdmin        ProSource = "ADMIN_GRANTED"
)

type UserMetadata struct {
	// Subscription and Plan
	AccessMode          AccessMode `json:"access_mode"`
	PlanType            PlanType   `json:"plan_type"`
	PlanExpirationDate  *time.Time `json:"plan_expiration_date,omitempty"`
	ProSource           *ProSource `json:"pro_source,omitempty"`
	MaxResources        *int       `json:"max_resources,omitempty"`
	MaxRequestsPerMonth *int       `json:"max_requests_per_month,omitempty"`

	// Limits
	MaxAccounts             int `json:"max_accounts"`
	MaxCategoriesPerAccount int `json:"max_categories_per_account"`
	MaxTransactionsPerMonth int `json:"max_transactions_per_month"`

	// Features
	CanExportData          bool `json:"can_export_data"`
	CanUseReports          bool `json:"can_use_reports"`
	CanUseAdvancedFeatures bool `json:"can_use_advanced_features"`
	CanCreateBudgets       bool `json:"can_create_budgets"`
	CanUseGoals            bool `json:"can_use_goals"`

	// Security and Status
	EmailVerified           bool             `json:"email_verified"`
	ReputationStatus        ReputationStatus `json:"reputation_status"`
	SuspiciousActivityCount int              `json:"suspicious_activity_count"`
	LastSecurityCheck       *time.Time       `json:"last_security_check,omitempty"`
	LastPermissionCheck     *time.Time       `json:"last_permission_check,omitempty"`

	// Localization and Notes
	Notes    *string `json:"notes,omitempty"`
	Locale   string  `json:"locale"`
	Currency string  `json:"currency"`
}

// Scan implements the sql.Scanner interface for UserMetadata
func (m *UserMetadata) Scan(value interface{}) error {
	if value == nil {
		return nil
	}
	bytes, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}
	return json.Unmarshal(bytes, m)
}

// Value implements the driver.Valuer interface for UserMetadata
func (m UserMetadata) Value() (driver.Value, error) {
	return json.Marshal(m)
}

type User struct {
	ID         int64        `gorm:"primaryKey;autoIncrement"`
	Name       string       `gorm:"not null"`
	Email      string       `gorm:"uniqueIndex;not null"`
	Password   *string      `gorm:"column:password"`
	ImgURL     *string      `gorm:"column:img_url;size:500"`
	Admin      bool         `gorm:"not null;default:false"`
	Active     bool         `gorm:"not null;default:true;index"`
	Source     string       `gorm:"not null;default:'LOCAL';size:50"`
	Metadata   UserMetadata `gorm:"type:jsonb"`
	LastAccess *time.Time   `gorm:"column:last_access"`
	CreatedAt  time.Time    `gorm:"not null"`
	UpdatedAt  time.Time
}

func NewDefaultMetadata() UserMetadata {
	return UserMetadata{
		AccessMode:              AccessModeReadWrite,
		PlanType:                PlanTypeFree,
		MaxAccounts:             5,
		MaxCategoriesPerAccount: 20,
		MaxTransactionsPerMonth: 200,
		CanExportData:           false,
		CanUseReports:           false,
		CanUseAdvancedFeatures:  false,
		CanCreateBudgets:        true,
		CanUseGoals:             false,
		EmailVerified:           false,
		ReputationStatus:        ReputationStatusGood,
		SuspiciousActivityCount: 0,
		Locale:                  "pt-BR",
		Currency:                "BRL",
	}
}
