package grpc

import (
	"context"

	usermanager "user-crud/pkg/pb"

	"google.golang.org/protobuf/types/known/timestamppb"
)

func (s *server) GetUser(ctx context.Context, r *usermanager.GetUserRequest) (*usermanager.GetUserResponse, error) {
	user, err := s.user.Get(ctx, r.GetId())
	if err != nil {
		return nil, err
	}

	return &usermanager.GetUserResponse{
		Id:        user.ID,
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Nickname:  user.Nickname,
		Email:     user.Email,
		Country:   user.Country,
		CreatedAt: timestamppb.New(user.CreatedAt),
		UpdatedAt: timestamppb.New(user.UpdatedAt),
	}, nil
}
