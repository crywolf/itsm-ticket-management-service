package externalusersvc

import (
	"context"
	"fmt"

	usermanagement "github.com/crywolf/itsm-ticket-management-service/external/itsm-user-service/api"
	"github.com/crywolf/itsm-ticket-management-service/internal/domain"
	"github.com/crywolf/itsm-ticket-management-service/internal/domain/ref"
	"github.com/crywolf/itsm-ticket-management-service/internal/domain/user"
	"github.com/crywolf/itsm-ticket-management-service/internal/domain/user/actor"
	"github.com/crywolf/itsm-ticket-management-service/internal/repository"
	"github.com/spf13/viper"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"google.golang.org/protobuf/types/known/emptypb"
)

// Service provides basic info about user
type Service interface {
	// ActorFromRequest calls external use service and returns an Actor object that represents a user who initiated the request
	// or about user this request is made on behalf of
	ActorFromRequest(ctx context.Context, authToken string, channelID ref.ChannelID, onBehalf string) (actor.Actor, error)
}

// ServiceCloser provides Service functionality plus allows to close connection to external service
type ServiceCloser interface {
	Service

	// Close tears down connection to external user service
	Close() error
}

// NewService creates new user service with initialized client for connection to external user service
func NewService(basicUserRepository repository.BasicUserRepository) (ServiceCloser, error) {
	conn, err := grpc.Dial(
		viper.GetString("UserServiceGRPCDialTarget"),
		grpc.WithInsecure(),
	)
	if err != nil {
		return nil, err
	}

	return &userService{
		conn:                conn,
		client:              usermanagement.NewUserManagementServiceClient(conn),
		basicUserRepository: basicUserRepository,
	}, nil
}

type userService struct {
	conn                *grpc.ClientConn
	client              usermanagement.UserManagementServiceClient
	basicUserRepository repository.BasicUserRepository
}

func (s userService) Close() error {
	return s.conn.Close()
}

func (s userService) ActorFromRequest(ctx context.Context, authToken string, channelID ref.ChannelID, onBehalf string) (actor.Actor, error) {
	basicUser, err := s.basicUserFromRequest(ctx, authToken, channelID, onBehalf)
	if err != nil {
		return actor.Actor{}, err
	}

	fmt.Println(basicUser)
	actorUser := actor.Actor{
		BasicUser: basicUser,
	}

	// TODO - try to find field engineer with this basicUser in repository and assign it
	//fieldEngineer := &fieldengineer.FieldEngineer{}
	//actorUser.SetFieldEngineerID(fieldEngineer)

	return actorUser, nil
}

func (s userService) basicUserFromRequest(ctx context.Context, authToken string, channelID ref.ChannelID, onBehalf string) (user.BasicUser, error) {
	md := metadata.New(map[string]string{
		"grpc-metadata-space": channelID.String(),
		"authorization":       authToken,
	})

	grpcCtx := metadata.NewOutgoingContext(context.Background(), md)

	var resp *usermanagement.UserPersonalDetailsResponse
	var err error

	if onBehalf != "" {
		resp, err = s.client.UserGet(grpcCtx, &usermanagement.UserRequest{Uuid: onBehalf})
		if err != nil {
			err = domain.WrapErrorf(err, domain.ErrorCodeUnknown, "authorization failed")
			return user.BasicUser{}, err
		}
	} else {
		resp, err = s.client.UserGetMyPersonalDetails(grpcCtx, &emptypb.Empty{})
		if err != nil {
			err = domain.WrapErrorf(err, domain.ErrorCodeUnknown, "authorization failed")
			return user.BasicUser{}, err
		}
	}

	u := resp.GetResult()

	//take returned ExternalUserUUID and get BasicUser from repository
	externalID := ref.ExternalUserUUID(u.Uuid)
	basicUser, err := s.basicUserRepository.GetBasicUserByExternalID(ctx, channelID, externalID)
	if err != nil {
		err = domain.WrapErrorf(err, domain.ErrorCodeUserNotAuthorized, "user could not be authorized")
	}

	return basicUser, err
}
