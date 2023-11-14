package main

import (
	"log"
	"net"
	red "service/internal/db/redis"
	"service/internal/service"
	"service/pkg/serv"

	"google.golang.org/grpc"
)

func main() {
	if err := run(":8001", "localhost:6379"); err != nil {
		log.Fatal(err)
	}
}

func run(address_serv, address_db string) error {
	list, err := net.Listen("tcp", address_serv)
	if err != nil {
		return err
	}
	serverRegistration := grpc.NewServer()
	nw := red.MyNewRedis(address_db)
	service := &service.MyAuthServer{Rb: nw}
	serv.RegisterAuthServer(serverRegistration, service)
	log.Printf("Running server: %s", address_serv)
	if err := serverRegistration.Serve(list); err != nil {
		return err
	}
	return nil
}
