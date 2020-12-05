package main

import (
	"fmt"
	"log"

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

	srv := micro.NewService(
		micro.Name("laracom.user.service"),
		micro.Version("latest"),
	)
	srv.Init()

	pb.RegisterUserServiceHandler(srv.Server(), &handler.UserService{Repo: repo})

	if err := srv.Run(); err != nil {
		fmt.Println(err)
	}
}
