package handler

import (
	"Assignment_Golang/user/entity"
	pb "Assignment_Golang/user/proto/user_service/v1"
	"Assignment_Golang/user/service"
	"context"
	"fmt"
	"log"
	"time"

	"github.com/gin-gonic/gin"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type IUserHandler interface {
	CreateUser(c *gin.Context)
	GetUserByID(c *gin.Context)
	UpdateUser(c *gin.Context)
	DeleteUser(c *gin.Context)
	GetAllUsersHandler(c *gin.Context)
}

type UserHandler struct {
	pb.UnimplementedUserServiceServer
	userService service.IUserService
}

func NewUserHandler(userService service.IUserService) *UserHandler {
	return &UserHandler{userService: userService}
}

func (h *UserHandler) Getusers(ctx context.Context, _ *emptypb.Empty) (*pb.GetUserResponse, error) {
	users, err := h.userService.GetAllUsers(ctx)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	var userProto []*pb.User
	for _, user := range users {
		userProto = append(userProto, &pb.User{
			Id:        int32(user.ID),
			Name:      user.Name,
			Email:     user.Email,
			Password:  user.Password,
			CreatedAt: timestamppb.New(user.CreatedAt),
			UpdatedAt: timestamppb.New(user.UpdatedAt),
		})
	}
	return &pb.GetUserResponse{
		Users: userProto,
	}, nil
}

func (h *UserHandler) GetUserByID(ctx context.Context, req *pb.GetUserByIDRequest) (*pb.GetUserByIDResponse, error) {
	user, err := h.userService.GetUserByID(ctx, int(req.Id))
	if err != nil {
		errres := &pb.GetUserByIDResponse{
			User: &pb.User{
				Id:        0,
				Name:      "",
				Email:     "",
				Password:  "",
				CreatedAt: timestamppb.New(time.Now()),
				UpdatedAt: timestamppb.New(time.Now()),
			},
		}
		log.Println("err")
		return errres, err
	}

	res := &pb.GetUserByIDResponse{
		User: &pb.User{
			Id:        int32(user.ID),
			Name:      user.Name,
			Email:     user.Email,
			Password:  user.Password,
			CreatedAt: timestamppb.New(user.CreatedAt),
			UpdatedAt: timestamppb.New(user.UpdatedAt),
		},
	}
	return res, nil
}

func (h *UserHandler) CreateUser(ctx context.Context, req *pb.CreateuserRequest) (*pb.MutationResponse, error) {
	createdUser, err := h.userService.CreateUser(ctx, &entity.User{
		Name:     req.Name,
		Email:    req.Email,
		Password: req.Password,
	})

	if err != nil {
		log.Println(err)
		return nil, err
	}
	return &pb.MutationResponse{
		Message: fmt.Sprintf("Success created user with ID %d", createdUser.ID),
	}, nil
}

func (h *UserHandler) UpdateUser(ctx context.Context, req *pb.UpdateUserRequest) (*pb.MutationResponse, error) {
	updatedUser, err := h.userService.UpdateUser(ctx, int(req.Id), entity.User{
		Name:     req.Name,
		Email:    req.Email,
		Password: req.Password,
	})
	if err != nil {
		log.Println(err)
		return nil, err
	}
	return &pb.MutationResponse{
		Message: fmt.Sprintf("Success updated user with ID %d", updatedUser.ID),
	}, nil
}

func (h *UserHandler) DeleteUser(ctx context.Context, req *pb.DeleteUserRequest) (*pb.MutationResponse, error) {
	if err := h.userService.DeleteUser(ctx, int(req.Id)); err != nil {
		log.Println(err)
		return nil, err
	}
	return &pb.MutationResponse{
		Message: fmt.Sprintf("Success deleted user with ID %d", req.Id),
	}, nil
}
