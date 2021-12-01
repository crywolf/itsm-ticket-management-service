package usersvc

import (
	"context"
	"net/http"

	"github.com/KompiTech/itsm-ticket-management-service/internal/domain/ref"
	"github.com/KompiTech/itsm-ticket-management-service/internal/domain/user"
	"github.com/KompiTech/itsm-ticket-management-service/internal/domain/user/actor"
	usermanagement "github.com/KompiTech/itsm-user-service/api/userservice"
	"github.com/spf13/viper"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"google.golang.org/protobuf/types/known/emptypb"
)

// Service provides basic info about user
type Service interface {
	// ActorFromRequest calls external use service and returns Actor object that represents a user who initiated the request
	// or about user this request is made on behalf of
	ActorFromRequest(r *http.Request) (actor.Actor, error)
}

// ServiceCloser provides Service functionality plus allows to close connection to external service
type ServiceCloser interface {
	Service

	// Close tears down connection to external user service
	Close() error
}

// NewService creates new user service with initialized client for connection to external user service
func NewService() (ServiceCloser, error) {
	conn, err := grpc.Dial(
		viper.GetString("UserServiceGRPCDialTarget"),
		grpc.WithInsecure(),
	)
	if err != nil {
		return nil, err
	}

	return &userService{
		conn:   conn,
		client: usermanagement.NewUserManagementServiceClient(conn),
	}, nil
}

type userService struct {
	conn   *grpc.ClientConn
	client usermanagement.UserManagementServiceClient
}

func (s userService) Close() error {
	return s.conn.Close()
}

func (s userService) ActorFromRequest(r *http.Request) (actor.Actor, error) {
	basicUser, err := s.basicUserFromRequest(r)
	if err != nil {
		return actor.Actor{}, err
	}

	actorUser := actor.Actor{
		BasicUser: basicUser,
	}

	// TODO - try to find field engineer with this basicUser in repository and assign it
	//fieldEngineer := &fieldengineer.FieldEngineer{}
	//actorUser.SetFieldEngineer(fieldEngineer)

	return actorUser, nil
}

func (s userService) basicUserFromRequest(r *http.Request) (user.BasicUser, error) {
	md := metadata.New(map[string]string{
		"grpc-metadata-space": r.Header.Get("channel-id"),
		"authorization":       r.Header.Get("authorization"),
	})

	ctx := metadata.NewOutgoingContext(context.Background(), md)

	var resp *usermanagement.UserPersonalDetailsResponse
	var err error

	if onBehalf := r.Header.Get("on_behalf"); onBehalf != "" {
		resp, err = s.client.UserGet(ctx, &usermanagement.UserRequest{Uuid: onBehalf})
		if err != nil {
			return user.BasicUser{}, err
		}
	} else {
		resp, err = s.client.UserGetMyPersonalDetails(ctx, &emptypb.Empty{})
		if err != nil {
			return user.BasicUser{}, err
		}
	}

	u := resp.GetResult()

	// TODO - take returned ExternalUserUUID and get BasicUser from repository

	userData := user.BasicUser{
		ExternalUserUUID: ref.ExternalUserUUID(u.Uuid),
		Name:             u.Name,
		Surname:          u.Surname,
		OrgName:          u.OrgName,
		OrgDisplayName:   u.OrgDisplayName,
	}

	return userData, nil
}
