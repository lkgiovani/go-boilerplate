package user

import (
	"context"

	"github.com/lkgiovani/go-boilerplate/pkg/encrypt"
	"gorm.io/gorm"
)

type GormRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) UserService {
	return &GormRepository{db: db}
}

func (r *GormRepository) Create(ctx context.Context, user *User) error {
	return r.db.WithContext(ctx).Create(user).Error
}

func (r *GormRepository) GetByEmail(ctx context.Context, email string) (*User, error) {
	var u User
	if err := r.db.WithContext(ctx).Where("email = ?", email).First(&u).Error; err != nil {
		return nil, err
	}
	return &u, nil
}

func (r *GormRepository) GetByID(ctx context.Context, id int64) (*User, error) {
	var u User
	if err := r.db.WithContext(ctx).Where("id = ?", id).First(&u).Error; err != nil {
		return nil, err
	}
	return &u, nil
}

func (r *GormRepository) Update(ctx context.Context, user *User) error {
	return r.db.WithContext(ctx).Save(user).Error
}

func (r *GormRepository) Delete(ctx context.Context, id int64) error {
	return r.db.WithContext(ctx).Delete(&User{}, id).Error
}

func (r *GormRepository) DeleteByIDs(ctx context.Context, ids []int64) error {
	return r.db.WithContext(ctx).Delete(&User{}, ids).Error
}

func (r *GormRepository) FindAll(ctx context.Context, page, size int) ([]User, int64, error) {
	var users []User
	var total int64

	offset := (page - 1) * size

	if err := r.db.WithContext(ctx).Model(&User{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if err := r.db.WithContext(ctx).Offset(offset).Limit(size).Find(&users).Error; err != nil {
		return nil, 0, err
	}

	return users, total, nil
}

func (r *GormRepository) FindAllWithFilter(ctx context.Context, keyword string, page, size int) ([]User, int64, error) {
	var users []User
	var total int64

	offset := (page - 1) * size
	query := r.db.WithContext(ctx).Model(&User{})

	// Filtrar por nome ou email
	if keyword != "" {
		query = query.Where("name ILIKE ? OR email ILIKE ?", "%"+keyword+"%", "%"+keyword+"%")
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if err := query.Offset(offset).Limit(size).Find(&users).Error; err != nil {
		return nil, 0, err
	}

	return users, total, nil
}

func (r *GormRepository) ToggleStatus(ctx context.Context, id int64, active bool) error {
	return r.db.WithContext(ctx).Model(&User{}).Where("id = ?", id).Update("active", active).Error
}

func (r *GormRepository) RequestPasswordReset(ctx context.Context, email string) error {
	// TODO: Implementar lógica de reset de senha
	return nil
}

func (r *GormRepository) ResetPassword(ctx context.Context, token, newPassword string) error {
	// TODO: Implementar lógica de reset de senha com token
	return nil
}

func (r *GormRepository) ChangePassword(ctx context.Context, id int64, currentPassword, newPassword string) error {
	// Buscar usuário
	user, err := r.GetByID(ctx, id)
	if err != nil {
		return err
	}

	// Verificar senha atual
	if user.Password != nil {
		if err := encrypt.VerifyPassword(currentPassword, *user.Password); err != nil {
			return err
		}
	}

	// Gerar hash da nova senha
	hashedPassword, err := encrypt.HashPassword(newPassword)
	if err != nil {
		return err
	}

	user.Password = &hashedPassword
	return r.Update(ctx, user)
}

func (r *GormRepository) ResetUserPassword(ctx context.Context, id int64, newPassword string) error {
	hashedPassword, err := encrypt.HashPassword(newPassword)
	if err != nil {
		return err
	}
	return r.db.WithContext(ctx).Model(&User{}).Where("id = ?", id).Update("password", hashedPassword).Error
}

func (r *GormRepository) UpdateAccessMode(ctx context.Context, id int64, accessMode string) (*User, error) {
	user, err := r.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	user.Metadata.AccessMode = AccessMode(accessMode)
	if err := r.Update(ctx, user); err != nil {
		return nil, err
	}
	return user, nil
}

func (r *GormRepository) UpdateFeatures(ctx context.Context, id int64, canCreateBudgets, canExportData, canUseReports, canUseGoals *bool) (*User, error) {
	user, err := r.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if canCreateBudgets != nil {
		user.Metadata.CanCreateBudgets = *canCreateBudgets
	}
	if canExportData != nil {
		user.Metadata.CanExportData = *canExportData
	}
	if canUseReports != nil {
		user.Metadata.CanUseReports = *canUseReports
	}
	if canUseGoals != nil {
		user.Metadata.CanUseGoals = *canUseGoals
	}

	if err := r.Update(ctx, user); err != nil {
		return nil, err
	}
	return user, nil
}

func (r *GormRepository) UpdateLimits(ctx context.Context, id int64, maxAccounts, maxTransactionsPerMonth, maxCategoriesPerAccount *int) (*User, error) {
	user, err := r.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if maxAccounts != nil {
		user.Metadata.MaxAccounts = *maxAccounts
	}
	if maxTransactionsPerMonth != nil {
		user.Metadata.MaxTransactionsPerMonth = *maxTransactionsPerMonth
	}
	if maxCategoriesPerAccount != nil {
		user.Metadata.MaxCategoriesPerAccount = *maxCategoriesPerAccount
	}

	if err := r.Update(ctx, user); err != nil {
		return nil, err
	}
	return user, nil
}

func (r *GormRepository) GrantLifetimePro(ctx context.Context, id int64, reason string) (*User, error) {
	user, err := r.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	user.Metadata.PlanType = PlanTypePro
	user.Metadata.Notes = &reason
	// Para Lifetime Pro, podemos definir uma data muito distante ou deixar nulo
	// Aqui vamos deixar nulo para representar "vitalício"
	user.Metadata.PlanExpirationDate = nil

	if err := r.Update(ctx, user); err != nil {
		return nil, err
	}
	return user, nil
}

func (r *GormRepository) RevokeLifetimePro(ctx context.Context, id int64) (*User, error) {
	user, err := r.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	user.Metadata.PlanType = PlanTypeFree
	if err := r.Update(ctx, user); err != nil {
		return nil, err
	}
	return user, nil
}

func (r *GormRepository) EnsureMetadata(ctx context.Context, id int64) (*User, error) {
	user, err := r.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	// Se o campo Locale (que é obrigatório no default) estiver vazio, assumimos que precisa de inicialização
	if user.Metadata.Locale == "" {
		user.Metadata = NewDefaultMetadata()
		if err := r.Update(ctx, user); err != nil {
			return nil, err
		}
	}

	return user, nil
}
