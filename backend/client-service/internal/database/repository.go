package database

import (
	"context"

	"github.com/shanto-323/Chat-Server-1/client-service/internal/database/model"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type UserRepository interface {
	InsertUser(ctx context.Context, user *model.User) error
	GetUser(ctx context.Context, username string) (*model.User, error)
	DeleteUser(ctx context.Context, username string) error
}

type gormRepository struct {
	db *gorm.DB
}

func NewUserRepository(url string) (UserRepository, error) {
	db, err := gorm.Open(postgres.Open(url), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	if err := db.AutoMigrate(&model.User{}); err != nil {
		return nil, err
	}

	return &gormRepository{
		db: db,
	}, nil
}

func (g *gormRepository) InsertUser(ctx context.Context, user *model.User) error {
	return g.db.WithContext(ctx).Create(user).Error
}

func (g *gormRepository) GetUser(ctx context.Context, username string) (*model.User, error) {
	user := model.User{}
	if err := g.db.WithContext(ctx).Model(&model.User{}).Where("username = ?", username).Find(&user).Error; err != nil {
		return nil, err
	}

	return &user, nil
}

func (g *gormRepository) DeleteUser(ctx context.Context, username string) error {
	return g.db.WithContext(ctx).Where("username = ?", username).Delete(&model.User{}).Error
}
