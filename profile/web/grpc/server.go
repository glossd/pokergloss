package grpc

import (
	"context"
	"github.com/glossd/pokergloss/gogrpc/grpcprofile"
	"github.com/glossd/pokergloss/profile/db"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func GetUsers(ctx context.Context, r *grpcprofile.GetUsersRequest) (*grpcprofile.GetUsersResponse, error) {
	if len(r.UserIds) == 0 {
		return &grpcprofile.GetUsersResponse{}, nil
	}
	profiles, err := db.FindAllProfilesByUserIDs(ctx, r.UserIds)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	userMap := make(map[string]*grpcprofile.Identity)
	for _, profile := range profiles {
		userMap[profile.UserID] = &grpcprofile.Identity{
			UserId:   profile.UserID,
			Username: profile.Username,
			Picture:  profile.Picture,
		}
	}

	return &grpcprofile.GetUsersResponse{Users: userMap}, nil
}
