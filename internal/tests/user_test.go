package tests

import (
	"context"
	"testing"
	"time"

	"user-crud/internal/models"
	"user-crud/pkg/paginator"

	"github.com/stretchr/testify/assert"
	"syreclabs.com/go/faker"
)

func Test_Create(t *testing.T) {
	user := mockUser()
	created, err := userService.Create(context.Background(), user)
	if err != nil {
		t.Fatalf("error creating user: %v", err)
	}
	if created.ID == "" {
		t.Errorf("id is not set")
	}
	if created.CreatedAt.IsZero() {
		t.Errorf("created_at is not set")
	}
	if created.UpdatedAt.IsZero() {
		t.Errorf("updated_at is not set")
	}

	created.ID = ""
	created.CreatedAt = time.Time{}
	created.UpdatedAt = time.Time{}
	created.Password = ""
	user.Password = ""

	assert.Equal(t, user, *created)
}

func Test_Delete(t *testing.T) {
	ctx := context.Background()
	user := mockUser()
	created, err := userService.Create(ctx, user)
	if err != nil {
		t.Fatalf("error creating user: %v", err)
	}
	if err = userService.Delete(ctx, created.ID); err != nil {
		t.Errorf("error deleting user: %v", err)
	}
}

func Test_Get(t *testing.T) {
	ctx := context.Background()
	user := mockUser()
	created, err := userService.Create(ctx, user)
	if err != nil {
		t.Fatalf("error creating user: %v", err)
	}
	created.Password = ""
	fetched, err := userService.Get(ctx, created.ID)

	assert.Nil(t, err)
	assert.Equal(t, *created, *fetched)
}

func Test_List(t *testing.T) {
	ctx := context.Background()
	user := mockUser()
	created, err := userService.Create(ctx, user)
	if err != nil {
		t.Fatalf("error creating user: %v", err)
	}

	created.Password = ""
	users, err := userService.List(ctx, paginator.New(1, 1), userService.Filter().ByEmail(created.Email))

	assert.Nil(t, err)

	if len(users.Users) == 0 {
		t.Errorf("list did not get any results")
	}
}

func Test_Update(t *testing.T) {
	ctx := context.Background()
	user := mockUser()
	created, err := userService.Create(ctx, user)
	if err != nil {
		t.Fatalf("error creating user: %v", err)
	}

	createdEmail := created.Email
	created.Email = "random@example.com"

	updated, err := userService.Update(ctx, *created)
	if err != nil {
		t.Fatalf("error updating user: %v", err)
	}

	if created.UpdatedAt.Equal(updated.UpdatedAt) {
		t.Errorf("updated_at field value is not updated")
	}

	if createdEmail == updated.Email {
		t.Errorf("email is not updated")
	}
}

func mockUser() models.User {
	return models.User{
		FirstName: faker.RandomString(8),
		LastName:  faker.RandomString(8),
		Nickname:  faker.RandomString(10),
		Password:  faker.RandomString(10),
		Email:     faker.RandomString(10),
		Country:   faker.RandomChoice([]string{"UK", "RU", "US"}),
	}
}
