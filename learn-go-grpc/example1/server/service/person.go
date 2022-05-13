package service

import (
	"context"

	"github.com/lxxxxxxxx/grpc-server/pb"
)

type PersonService struct {
	pb.UnimplementedPersonServer
}

func (PersonService) GetPerson(ctx context.Context, req *pb.GetPersonReq) (*pb.GetPersonRes, error) {
	return &pb.GetPersonRes{Person: &pb.PersonInfo{Name: req.Name, Age: 18}}, nil
}
