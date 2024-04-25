package store

import (
	"errors"
	"fmt"

	"github.com/voidmaindev/doctra_lis_middleware/model"
	"gorm.io/gorm"
)

type UserStore struct {
	db *gorm.DB
}

func NewUserStore(db *gorm.DB) (*UserStore, error) {
	store := &UserStore{db: db}
	err := store.db.AutoMigrate(&model.User{})
	if err != nil {
		return nil, errors.New("failed to migrate User model")
	}

	return store, nil
}

func (s *UserStore) Create(user *model.User) error {
	err := s.db.Create(user).Error
	if err != nil {
		return fmt.Errorf("failed to create user: %v", user.ID)
	}

	return nil
}

func (s *UserStore) GetByID(id uint) (*model.User, error) {
	user := &model.User{}
	err := s.db.First(user, id).Error
	if err != nil {
		return nil, fmt.Errorf("failed to get user by ID: %v", id)
	}

	return user, nil
}

func (s *UserStore) GetByUsername(username string) (*model.User, error) {
	user := &model.User{}
	err := s.db.Where("username = ?", username).First(user).Error
	if err != nil {
		return nil, fmt.Errorf("failed to get user by username: %v", username)
	}

	return user, nil
}

func (s *UserStore) GetByRole(role string) ([]*model.User, error) {
	users := []*model.User{}
	err := s.db.Where("role = ?", role).Find(&users).Error
	if err != nil {
		return nil, fmt.Errorf("failed to get users by role: %v", role)
	}

	return users, nil
}

func (s *UserStore) GetAll() ([]model.User, error) {
	users := []model.User{}
	err := s.db.Find(&users).Error
	if err != nil {
		return nil, errors.New("failed to get all users")
	}

	return users, nil
}

func (s *UserStore) Update(user *model.User) error {
	err := s.db.Save(user).Error
	if err != nil {
		return fmt.Errorf("failed to update user: %v", user.ID)
	}

	return nil
}

func (s *UserStore) Delete(user *model.User) error {
	err := s.db.Delete(user).Error
	if err != nil {
		return fmt.Errorf("failed to delete user: %v", user.ID)
	}

	return nil
}
