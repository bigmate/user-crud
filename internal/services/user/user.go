package user

import (
	"context"
	"log"

	"user-crud/internal/models"
	"user-crud/internal/repository"
	"user-crud/internal/services/notifier"
	"user-crud/pkg/filter"
	"user-crud/pkg/paginator"

	"golang.org/x/crypto/bcrypt"
)

const (
	encryptionCost = 12
)

type Service interface {
	Create(ctx context.Context, user models.User) (*models.User, error)
	Update(ctx context.Context, user models.User) (*models.User, error)
	Get(ctx context.Context, id string) (*models.User, error)
	List(ctx context.Context, paginator paginator.Paginator, filters ...filter.Filter) (*models.UsersList, error)
	Delete(ctx context.Context, id string) error
	Filter() FilterFactory
}

type service struct {
	repo          repository.User
	notifier      notifier.Service
	filterFactory filterFactory
}

func (s *service) Filter() FilterFactory {
	return s.filterFactory
}

func (s *service) List(ctx context.Context, paginator paginator.Paginator, filters ...filter.Filter) (*models.UsersList, error) {
	return s.repo.List(ctx, paginator, filters...)
}

func (s *service) Create(ctx context.Context, user models.User) (*models.User, error) {
	if user.Password != "" {
		pw, err := encryptPassword(user.Password)
		if err != nil {
			return nil, err
		}

		user.Password = string(pw)
	}

	return s.repo.Insert(ctx, user)
}

func (s *service) Update(ctx context.Context, user models.User) (*models.User, error) {
	if user.Password != "" {
		pw, err := encryptPassword(user.Password)
		if err != nil {
			return nil, err
		}

		user.Password = string(pw)
	}

	updatedUser, updated, err := s.repo.Update(ctx, user)
	if err != nil {
		return nil, err
	}

	if updated {
		if err = s.notifier.Publish("user-change", NewModelEncoder(updatedUser)); err != nil {
			log.Printf("user: failed to publish message: %v", err)
		}
	}

	return updatedUser, nil
}

func (s service) Get(ctx context.Context, id string) (*models.User, error) {
	return s.repo.Get(ctx, id)
}

func (s *service) Delete(ctx context.Context, id string) error {
	return s.repo.Delete(ctx, id)
}

func encryptPassword(password string) ([]byte, error) {
	return bcrypt.GenerateFromPassword([]byte(password), encryptionCost)
}

func NewService(repo repository.User, notifier notifier.Service) Service {
	return &service{
		repo:          repo,
		notifier:      notifier,
		filterFactory: filterFactory{},
	}
}
