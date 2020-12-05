package main

import (
	"context"
	"log"
	"os"

	pb "github.com/JinesD/laracom.service/user.service/proto/user"
	"github.com/micro/cli/v2"
	"github.com/micro/go-micro/v2"
)

func main() {
	srv := micro.NewService(
		micro.Flags(
			&cli.StringFlag{
				Name:  "name",
				Usage: "Your Name",
			},
			&cli.StringFlag{
				Name:  "email",
				Usage: "Your Email",
			},
			&cli.StringFlag{
				Name:  "password",
				Usage: "Your Password",
			},
		),
	)

	client := pb.NewUserService("laracom.user.service", srv.Client())

	srv.Init(
		micro.Action(func(c *cli.Context) error {
			name := c.String("name")
			email := c.String("email")
			password := c.String("password")

			log.Println("参数：", name, email, password)

			r, err := client.Create(context.TODO(), &pb.User{
				Name:     name,
				Email:    email,
				Password: password,
			})
			if err != nil {
				log.Fatalf("创建用户失败：%v", err)
			}
			log.Printf("创建用户成功：%v", r.User.Id)

			token, err := client.Auth(context.TODO(), &pb.User{
				Email:    email,
				Password: password,
			})
			if err != nil {
				log.Fatalf("failed to logging in: %v", err)
			}
			log.Printf("logging in success: %v", token.Token)

			token, err = client.ValidateToken(context.TODO(), token)
			if err != nil {
				log.Fatalf("failed to validate token: %v", err)
			}
			log.Printf("validate token success: %v", token.Valid)

			getAll, err := client.GetAll(context.Background(), &pb.Request{})
			if err != nil {
				log.Fatalf("获取所有用户失败：%v", err)
			}
			for _, v := range getAll.Users {
				log.Println(v)
			}

			os.Exit(0)

			return nil
		}),
	)

	if err := srv.Run(); err != nil {
		log.Fatalf("用户客户端启动失败：%v", err)
	}
}
