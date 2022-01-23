package grpc

import (
	"context"

	usermanager "user-crud/pkg/pb"
)

func (s *server) DeleteUser(ctx context.Context, r *usermanager.DeleteUserRequest) (*usermanager.DeleteUserResponse, error) {

	if err := s.user.Delete(ctx, r.GetId()); err != nil {
		return nil, err
	}

	return &usermanager.DeleteUserResponse{}, nil
}
