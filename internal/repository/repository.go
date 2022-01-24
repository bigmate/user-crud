package repository

import (
	"context"

	"user-crud/internal/models"
	"user-crud/pkg/filter"
	"user-crud/pkg/paginator"
)

//User is an interface to talk to persistent storage
type User interface {
	Get(ctx context.Context, id string) (*models.User, error)
	Insert(ctx context.Context, user models.User) (*models.User, error)
	Update(ctx context.Context, user models.User) (updatedUser *models.User, updated bool, err error)
	List(ctx context.Context, p paginator.Paginator, filters ...filter.Filter) (*models.UsersList, error)
	Delete(ctx context.Context, id string) error
}
