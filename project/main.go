package main

import (
	"fmt"
	"log"
	"net"

	"github.com/Mensurui/clerkService/project/middleware"
	"github.com/Mensurui/clerkService/project/service"
	protos "github.com/Mensurui/clerkService/proto/golang/service"
	"github.com/clerk/clerk-sdk-go/v2"
	"github.com/hashicorp/go-hclog"
	"github.com/joho/godotenv"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	clerk.SetKey("CLERK_SECRET_KEY")

	l := hclog.Default()
	servInterceptor := middleware.ServInterceptor{}
	gs := grpc.NewServer(
		grpc.UnaryInterceptor(servInterceptor.Unary()),
	)
	ss := service.NewService(l)
	protos.RegisterServiceServer(gs, ss)
	reflection.Register(gs)
	l.Info("Starting server")

	servicePort := "8080"
	listen, err := net.Listen("tcp", fmt.Sprintf("0.0.0.0:%s", servicePort))
	if err != nil {
		l.Error("Unable to listen", "error", err)
	}
	err = gs.Serve(listen)
	if err != nil {
		l.Error("Unable to serve", "error", err)
	}

	l.Info("Server started")
}
