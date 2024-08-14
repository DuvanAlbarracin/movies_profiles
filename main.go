package main

import (
	"log"
	"net"

	"github.com/DuvanAlbarracin/movies_profiles/pkg/config"
	"github.com/DuvanAlbarracin/movies_profiles/pkg/db"
	"github.com/DuvanAlbarracin/movies_profiles/pkg/proto"
	"github.com/DuvanAlbarracin/movies_profiles/pkg/services"
	"google.golang.org/grpc"
)

func main() {
	c, err := config.LoadConfig()
	if err != nil {
		log.Fatalln("Failed loading config:", err)
	}

	h := db.Init(c.DBUrl)

	lis, err := net.Listen("tcp", c.Port)
	if err != nil {
		log.Fatalln("Failed to listening:", err)
	}
	log.Println("Auth service on:", c.Port)

	s := services.Server{
		H: h,
	}
	defer s.H.Conn.Close()

	grpcServer := grpc.NewServer()
	proto.RegisterProfilesServiceServer(grpcServer, &s)

	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalln("Failed to serve:", err)
	}
}
