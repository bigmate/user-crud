package grpc

import (
	"context"

	"user-crud/internal/models"
	usermanager "user-crud/pkg/pb"

	"google.golang.org/protobuf/types/known/timestamppb"
)

func (s *server) CreateUser(ctx context.Context, r *usermanager.CreateUserRequest) (*usermanager.CreateUserResponse, error) {
	user := models.User{
		FirstName: r.GetFirstName(),
		LastName:  r.GetLastName(),
		Nickname:  r.GetNickname(),
		Password:  r.GetPassword(),
		Email:     r.GetEmail(),
		Country:   r.GetCountry(),
	}

	createdUser, err := s.user.Create(ctx, user)
	if err != nil {
		return nil, err
	}

	return &usermanager.CreateUserResponse{
		Id:        createdUser.ID,
		FirstName: createdUser.FirstName,
		LastName:  createdUser.LastName,
		Nickname:  createdUser.Nickname,
		Email:     createdUser.Email,
		Country:   createdUser.Country,
		CreatedAt: timestamppb.New(createdUser.CreatedAt),
		UpdatedAt: timestamppb.New(createdUser.UpdatedAt),
	}, nil
}
