package repository

import (
	"context"
	"errors"
	"fmt"
	"final-project/models"
	"final-project/helpers"
	"time"

	nanoid "github.com/matoous/go-nanoid/v2"
	"gorm.io/gorm"
)

type userRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) *userRepository {
	return &userRepository{db}
}

func (userRepository *userRepository) Register(ctx context.Context, user *models.User) (err error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)

	defer cancel()

	ID, _ := nanoid.New(16)

	user.ID = fmt.Sprintf("user-%s", ID)

	if err = userRepository.db.WithContext(ctx).Create(&user).Error; err != nil {
		return err
	}

	return
}

func (userRepository *userRepository) Login(ctx context.Context, user *models.User) (err error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)

	defer cancel()

	password := user.Password

	if err = userRepository.db.WithContext(ctx).Where("email = ?", user.Email).Take(&user).Error; err != nil {
		return errors.New("the email you entered are not found")
	}

	if isValid := helpers.Compare([]byte(user.Password), []byte(password)); !isValid {
		return errors.New("the credential you entered are wrong")
	}

	return
}

func (userRepository *userRepository) Update(ctx context.Context, user models.User) (u models.User, err error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)

	defer cancel()

	u = models.User{}

	if err = userRepository.db.WithContext(ctx).First(&u).Error; err != nil {
		return u, err
	}

	if err = userRepository.db.WithContext(ctx).Model(&u).Updates(user).Error; err != nil {
		return u, err
	}

	return u, nil
}

func (userRepository *userRepository) Delete(ctx context.Context, id string) (err error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)

	defer cancel()

	if err = userRepository.db.WithContext(ctx).First(&models.User{}, &id).Error; err != nil {
		return err
	}

	if err = userRepository.db.WithContext(ctx).Where("user_id = ?", id).Delete(&models.SocialMedia{}).Error; err != nil {
		return err
	}

	if err = userRepository.db.WithContext(ctx).Delete(&models.User{}, &id).Error; err != nil {
		return err
	}

	return
}
