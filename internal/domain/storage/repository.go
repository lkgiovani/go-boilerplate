package storage

import (
	"context"

	"gorm.io/gorm"
)

type GormFileRepository struct {
	db *gorm.DB
}

func NewFileRepository(db *gorm.DB) FileRepository {
	return &GormFileRepository{db: db}
}

func (r *GormFileRepository) Save(ctx context.Context, file *FileReference) error {
	return r.db.WithContext(ctx).Create(file).Error
}

func (r *GormFileRepository) GetByID(ctx context.Context, id int64) (*FileReference, error) {
	var file FileReference
	if err := r.db.WithContext(ctx).First(&file, id).Error; err != nil {
		return nil, err
	}
	return &file, nil
}

func (r *GormFileRepository) GetByStorageKey(ctx context.Context, key string) (*FileReference, error) {
	var file FileReference
	if err := r.db.WithContext(ctx).Where("storage_key = ?", key).First(&file).Error; err != nil {
		return nil, err
	}
	return &file, nil
}

func (r *GormFileRepository) Delete(ctx context.Context, id int64) error {
	return r.db.WithContext(ctx).Delete(&FileReference{}, id).Error
}

func (r *GormFileRepository) FindByUserID(ctx context.Context, userID int64) ([]FileReference, error) {
	var files []FileReference
	if err := r.db.WithContext(ctx).Where("user_id = ?", userID).Find(&files).Error; err != nil {
		return nil, err
	}
	return files, nil
}
