package main

// https://www.bilibili.com/video/BV1eE411T7GC?p=1

import (
	"log"
	"net"

	"github.com/lxxxxxxxx/grpc-server/pb"
	"github.com/lxxxxxxxx/grpc-server/service"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

const port = ":5001"

func main() {
	listen, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalln(err.Error())
	}
	creds, err := credentials.NewServerTLSFromFile("./keys/server.pem", "./keys/ca.key")
	if err != nil {
		log.Fatalln(err.Error())
	}
	options := []grpc.ServerOption{grpc.Creds(creds)}
	server := grpc.NewServer(options...)
	pb.RegisterEmployeeServer(server, new(service.EmployeeService))
	pb.RegisterPersonServer(server, new(service.PersonService))
	log.Println("grpc server listening " + port)

	server.Serve(listen)

}
