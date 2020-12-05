package handler

import (
	"context"

	pb "github.com/JinesD/laracom.service/user.service/proto/user"
	"github.com/JinesD/laracom.service/user.service/repo"
	"golang.org/x/crypto/bcrypt"
)

type UserService struct {
	Repo repo.Repository
}

func (srv *UserService) Get(ctx context.Context, req *pb.User, resp *pb.Response) error {
	user, err := srv.Repo.Get(req.Id)
	if err != nil {
		return err
	}

	resp.User = user

	return nil
}

func (srv *UserService) GetAll(ctx context.Context, req *pb.Request, resp *pb.Response) error {
	users, err := srv.Repo.GetAll()
	if err != nil {
		return err
	}

	resp.Users = users

	return nil
}

func (srv *UserService) Create(ctx context.Context, req *pb.User, resp *pb.Response) error {
	hashed, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	req.Password = string(hashed)
	if err := srv.Repo.Create(req); err != nil {
		return err
	}

	resp.User = req

	return nil
}
