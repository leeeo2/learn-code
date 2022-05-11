package main

// https://www.bilibili.com/video/BV1eE411T7GC?p=1

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log"
	"net"
	"time"

	"github.com/lxxxxxxxx/grpc-server/pb"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/metadata"
	"google.golang.org/protobuf/types/known/timestamppb"
)

const port = ":5001"

var employees = []pb.Employee{
	{
		Id:   1,
		No:   1999,
		Name: "Dave",
		MonthSalary: &pb.MonthSalary{
			Basic: 5000,
			Bonus: 1000,
		},
		Status: pb.EmployeeStatus_NORMAL,
		LastModified: &timestamppb.Timestamp{
			Seconds: time.Now().Unix(),
		},
	},
	{
		Id:   2,
		No:   1996,
		Name: "Lili",
		MonthSalary: &pb.MonthSalary{
			Basic: 6000,
			Bonus: 500,
		},
		Status: pb.EmployeeStatus_NORMAL,
		LastModified: &timestamppb.Timestamp{
			Seconds: time.Now().Unix(),
		},
	},
}

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
	pb.RegisterEmployeeServiceServer(server, new(employeeService))

	log.Println("grpc server listening " + port)

	server.Serve(listen)

}

type employeeService struct {
	pb.UnimplementedEmployeeServiceServer
}

func (*employeeService) GetByNo(ctx context.Context, req *pb.GetByNoRequest) (*pb.EmployeeResponse, error) {
	for _, e := range employees {
		if e.No == req.No {
			fmt.Println("find employee no:", e.No)
			return &pb.EmployeeResponse{Employee: &e}, nil
		}
	}
	return nil, errors.New("employee not found")
}

func (*employeeService) GetAll(req *pb.GetAllRequest, stream pb.EmployeeService_GetAllServer) error {
	for _, e := range employees {
		stream.Send(&pb.EmployeeResponse{
			Employee: &e,
		})
		time.Sleep(time.Second)
	}
	return nil
}

func (*employeeService) AllPhoto(stream pb.EmployeeService_AllPhotoServer) error {
	md, ok := metadata.FromIncomingContext(stream.Context())
	if ok {
		fmt.Printf("Employee no:%s \n", md["no"][0])
	}

	img := []byte{}
	for {
		data, err := stream.Recv()
		if err == io.EOF {
			fmt.Printf("img size: %d\n", len(img))
			stream.SendAndClose(&pb.AddPhotoResponse{IsOk: true})
			return nil
		}
		if err != nil {
			log.Fatalln(err.Error())
			return err
		}

		fmt.Printf("recv size:%d \n", len(data.Data))
		img = append(img, data.Data...)
		time.Sleep(time.Millisecond * 500)
	}
	return errors.New("add photo failed.")
}

func (*employeeService) Save(ctx context.Context, req *pb.SaveEmployeeRequest) (*pb.EmployeeResponse, error) {
	for _, e := range employees {
		if e.No == req.Employee.No {
			fmt.Println("employee exist,no:", e.No)
			return &pb.EmployeeResponse{
				Employee: req.Employee,
			}, nil
		}
	}
	employees = append(employees, *req.Employee)

	return &pb.EmployeeResponse{
		Employee: req.Employee,
	}, nil
}
func (*employeeService) SaveAll(stream pb.EmployeeService_SaveAllServer) error {
	for {
		empReq, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatalln(err.Error())
			return err
		}
		employees = append(employees, *empReq.Employee)
		stream.Send(&pb.EmployeeResponse{Employee: empReq.Employee})
	}

	for _, e := range employees {
		fmt.Println(e)
	}
	return nil
}
