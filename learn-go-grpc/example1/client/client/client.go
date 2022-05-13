package client

import (
	"context"
	"log"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"

	"github.com/lxxxxxxx/grpc-client/pb"
)

var (
	EmployeeService pb.EmployeeClient
	PersonService   pb.PersonClient
)

func InitClient(ctx context.Context, port string) error {
	creds, err := credentials.NewClientTLSFromFile("./keys/server.pem", "127.0.0.1")
	if err != nil {
		log.Fatalln(err.Error())
	}

	options := []grpc.DialOption{grpc.WithTransportCredentials(creds)}
	conn, err := grpc.Dial("localhost"+port, options...)
	if err != nil {
		log.Fatalln(err.Error())
	}

	EmployeeService = pb.NewEmployeeClient(conn)
	PersonService = pb.NewPersonClient(conn)
	return nil
}
