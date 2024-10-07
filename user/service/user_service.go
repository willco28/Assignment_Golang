package service

import (
	"Assignment_Golang/user/entity"
	"Assignment_Golang/user/repository"
	"context"
	"fmt"
)

// service ==> buat interface --> struct --> constructor nya --> bikin function implement method dari interface nya
type IUserService interface { //interface dengan method getallusers (get all user untuk service)
	CreateUser(ctx context.Context, user *entity.User) (entity.User, error)
	GetUserByID(ctx context.Context, id int) (entity.User, error)
	UpdateUser(ctx context.Context, id int, user entity.User) (entity.User, error)
	DeleteUser(ctx context.Context, id int) error
	GetAllUsers(ctx context.Context) ([]entity.User, error)
}

type userService struct {
	userRepo repository.IUserRepository //bikin struct dengan balikan interface repo
}

func NewUserService(userRepo repository.IUserRepository) IUserService {
	//bikin constructor dengan balikan user service yang isinya user repo
	return &userService{userRepo: userRepo}
}

func (s *userService) CreateUser(ctx context.Context, user *entity.User) (entity.User, error) {
	//kalau ada validasi ditambahkan disini
	createduser, err := s.userRepo.CreateUser(ctx, user)
	if err != nil {
		return entity.User{}, fmt.Errorf("error created user: %v", err)
	}
	return createduser, nil
}

func (s *userService) GetUserByID(ctx context.Context, id int) (entity.User, error) {
	user, err := s.userRepo.GetUserByID(ctx, id)
	if err != nil {
		return entity.User{}, fmt.Errorf("user not found: %v", err)
	}
	return user, nil
}

func (s *userService) UpdateUser(ctx context.Context, id int, user entity.User) (entity.User, error) {
	//kalau ada validasi ditambahkan disini
	updateduser, err := s.userRepo.UpdateUser(ctx, id, user)
	if err != nil {
		return entity.User{}, fmt.Errorf("failed to update user: %v", err)
	}
	return updateduser, nil
}

func (s *userService) DeleteUser(ctx context.Context, id int) error {
	if err := s.userRepo.DeleteUser(ctx, id); err != nil {
		return fmt.Errorf("failed to delete user: %v", err)
	}
	return nil
}

func (s *userService) GetAllUsers(ctx context.Context) ([]entity.User, error) {
	//bikin function untuk mengembalikan user repo
	users, err := s.userRepo.GetAllUser(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to retreive all users: %v", err)
	}
	return users, nil
}
