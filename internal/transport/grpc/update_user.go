package grpc

import (
	"context"

	"user-crud/internal/models"
	usermanager "user-crud/pkg/pb"

	"google.golang.org/protobuf/types/known/timestamppb"
)

func (s *server) UpdateUser(ctx context.Context, r *usermanager.UpdateUserRequest) (*usermanager.UpdateUserResponse, error) {
	user := models.User{
		ID:        r.GetId(),
		FirstName: r.GetFirstName(),
		LastName:  r.GetLastName(),
		Nickname:  r.GetNickname(),
		Email:     r.GetEmail(),
		Country:   r.GetCountry(),
	}

	updatedUser, err := s.user.Update(ctx, user)
	if err != nil {
		return nil, err
	}

	return &usermanager.UpdateUserResponse{
		Id:        updatedUser.ID,
		FirstName: updatedUser.FirstName,
		LastName:  updatedUser.LastName,
		Nickname:  updatedUser.Nickname,
		Email:     updatedUser.Email,
		Country:   updatedUser.Country,
		CreatedAt: timestamppb.New(updatedUser.CreatedAt),
		UpdatedAt: timestamppb.New(updatedUser.UpdatedAt),
	}, nil
}
