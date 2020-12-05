package main

import (
	"fmt"
	"log"

	"github.com/JinesD/laracom.service/user.service/service"

	database "github.com/JinesD/laracom.service/user.service/db"
	"github.com/JinesD/laracom.service/user.service/handler"
	pb "github.com/JinesD/laracom.service/user.service/proto/user"
	repository "github.com/JinesD/laracom.service/user.service/repo"
	"github.com/micro/go-micro/v2"
)

func main() {
	db, err := database.CreateConnection()
	defer db.Close()

	if err != nil {
		log.Fatalf("Could not connect to db: %v", err)
	}

	db.AutoMigrate(&pb.User{})

	repo := &repository.UserRepository{Db: db}
	token := &service.TokenService{Repo: repo}

	srv := micro.NewService(
		micro.Name("laracom.user.service"),
		micro.Version("latest"),
	)
	srv.Init()

	if err := pb.RegisterUserServiceHandler(srv.Server(), &handler.UserService{Repo: repo, Token: token}); err != nil {
		fmt.Println(err)
	}

	if err := srv.Run(); err != nil {
		fmt.Println(err)
	}
}
