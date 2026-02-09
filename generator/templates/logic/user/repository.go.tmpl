package user

import (
	"context"

	"github.com/lpphub/goweb/ext/dbx"
	"gorm.io/gorm"
)

type UserRepo struct {
	*dbx.BaseRepo[User]
}

func NewUserRepo(db *gorm.DB) *UserRepo {
	return &UserRepo{
		BaseRepo: dbx.NewBaseRepo[User](db),
	}
}

func (r *UserRepo) GetByEmail(ctx context.Context, email string) (*User, error) {
	var user User
	if err := r.DB().WithContext(ctx).Where("email = ?", email).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *UserRepo) ExistsByEmail(ctx context.Context, email string) (bool, error) {
	var count int64
	if err := r.DB().WithContext(ctx).Model(&User{}).Where("email = ?", email).Count(&count).Error; err != nil {
		return false, err
	}
	return count > 0, nil
}

func (r *UserRepo) Update(ctx context.Context, userID uint, updates map[string]interface{}) error {
	return dbx.TxAwareDB(ctx, r.DB()).Model(&User{}).Where("id = ?", userID).Updates(updates).Error
}

func (r *UserRepo) GetByEmails(ctx context.Context, emails []string) ([]*User, error) {
	var users []*User
	if err := r.DB().WithContext(ctx).
		Where("email IN ?", emails).
		Find(&users).Error; err != nil {
		return nil, err
	}
	return users, nil
}
