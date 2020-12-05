package handler

import (
	"context"
	"errors"
	"log"

	"github.com/JinesD/laracom.service/user.service/service"

	pb "github.com/JinesD/laracom.service/user.service/proto/user"
	"github.com/JinesD/laracom.service/user.service/repo"
	"golang.org/x/crypto/bcrypt"
)

type UserService struct {
	Repo  repo.Repository
	Token service.Authable
}

func (srv *UserService) Auth(ctx context.Context, req *pb.User, resp *pb.Token) error {
	log.Println("Logging in with: ", req.Email, req.Password)

	user, err := srv.Repo.GetByEmail(req.Email)
	log.Println(user)
	if err != nil {
		return err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		return err
	}

	token, err := srv.Token.Encode(user)
	if err != nil {
		return err
	}

	resp.Token = token

	return nil
}

func (srv *UserService) ValidateToken(ctx context.Context, req *pb.Token, resp *pb.Token) error {
	claims, err := srv.Token.Decode(req.Token)
	if err != nil {
		return err
	}

	if claims.User.Id == "" {
		return errors.New("invalid user")
	}

	resp.Valid = true

	return nil
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
