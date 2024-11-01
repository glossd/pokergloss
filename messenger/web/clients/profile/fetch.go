package profile

import (
	"context"
	"fmt"
	"github.com/glossd/pokergloss/gogrpc/grpcprofile"
	"github.com/glossd/pokergloss/messenger/domain"
	"github.com/glossd/pokergloss/profile/web/grpc"
	log "github.com/sirupsen/logrus"
)

func FetchUserMap(ctx context.Context, ulc *domain.UserChatList) map[string]*grpcprofile.Identity {
	var userMap map[string]*grpcprofile.Identity
	var userIDs []string
	for _, chat := range ulc.U2UChats {
		userIDs = append(userIDs, chat.OtherUserID)
	}

	res, err := grpc.GetUsers(ctx, &grpcprofile.GetUsersRequest{UserIds: userIDs})
	if err != nil {
		log.Errorf("Failed to get users from profile: %s", err)
	} else {
		userMap = res.GetUsers()
	}

	if userMap == nil {
		userMap = make(map[string]*grpcprofile.Identity)
	}
	return userMap
}

func FetchUser(ctx context.Context, userID string) (*grpcprofile.Identity, error) {
	res, err := grpc.GetUsers(ctx, &grpcprofile.GetUsersRequest{UserIds: []string{userID}})
	if err != nil {
		log.Errorf("Failed to get users from profile: %s", err)
		return nil, err
	}
	if len(res.Users) == 0 {
		log.Errorf("No profiles found with userId %s", userID)
		return nil, fmt.Errorf("no user found %s", err)
	}
	return res.Users[userID], nil
}
