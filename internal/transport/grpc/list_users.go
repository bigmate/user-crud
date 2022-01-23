package grpc

import (
	"context"

	"user-crud/pkg/paginator"
	usermanager "user-crud/pkg/pb"

	"google.golang.org/protobuf/types/known/timestamppb"
)

func (s *server) ListUsers(ctx context.Context, r *usermanager.ListUsersRequest) (*usermanager.ListUsersResponse, error) {
	list, err := s.user.List(ctx,
		paginator.New(r.GetPageNumber(), r.GetPageSize()),
		s.user.Filter().ByFirstName(r.GetFirstName()),
		s.user.Filter().ByLastName(r.GetLastName()),
		s.user.Filter().ByCountry(r.GetCountry()),
		s.user.Filter().ByEmail(r.GetEmail()),
	)
	if err != nil {
		return nil, err
	}

	lu := &usermanager.ListUsersResponse{
		Users:      make([]*usermanager.ListUsersResponse_User, 0, len(list.Users)),
		TotalCount: list.TotalCount,
	}

	for _, u := range list.Users {
		user := u
		lu.Users = append(lu.Users, &usermanager.ListUsersResponse_User{
			Id:        user.ID,
			FirstName: user.FirstName,
			LastName:  user.LastName,
			Nickname:  user.Nickname,
			Email:     user.Email,
			Country:   user.Country,
			CreatedAt: timestamppb.New(user.CreatedAt),
			UpdatedAt: timestamppb.New(user.UpdatedAt),
		})
	}

	return lu, nil
}
