package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"os"
	"time"

	"github.com/lxxxxxxx/grpc-client/pb"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/metadata"
	"google.golang.org/protobuf/types/known/timestamppb"
)

const port = ":5001"

func main() {
	creds, err := credentials.NewClientTLSFromFile("./keys/server.pem", "127.0.0.1")
	if err != nil {
		log.Fatalln(err.Error())
	}

	options := []grpc.DialOption{grpc.WithTransportCredentials(creds)}
	conn, err := grpc.Dial("localhost"+port, options...)
	if err != nil {
		log.Fatalln(err.Error())
	}
	defer conn.Close()

	// 和server不同的是，grpc会实现client的接口
	// 我们只需要调用这些接口取数据就可以
	client := pb.NewEmployeeServiceClient(conn)
	getByNo(client)
	getAll(client)
	addPhoto(client)
	saveAll(client)
}

func getByNo(client pb.EmployeeServiceClient) {
	res, err := client.GetByNo(context.Background(), &pb.GetByNoRequest{No: 1999})
	if err != nil {
		log.Fatalln(err.Error())
	}
	fmt.Println(res.Employee)
}

func getAll(client pb.EmployeeServiceClient) {
	stream, err := client.GetAll(context.Background(), &pb.GetAllRequest{})

	if err != nil {
		log.Fatalln(err.Error())
	}
	for {
		res, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatalln(err.Error())
		}
		fmt.Println(res.Employee)
	}
}

func addPhoto(client pb.EmployeeServiceClient) {
	imgFile, err := os.Open("avatar.jpeg")
	if err != nil {
		log.Fatalln(err.Error())
	}
	defer imgFile.Close()

	md := metadata.New(map[string]string{"no": "1996"})
	context := context.Background()
	context = metadata.NewOutgoingContext(context, md)

	stream, err := client.AllPhoto(context)
	if err != nil {
		log.Fatalln(err.Error())
	}

	for {
		chunk := make([]byte, 32*1024)
		chunkSize, err := imgFile.Read(chunk)
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatalln("--- " + err.Error())
		}

		if chunkSize < len(chunk) {
			chunk = chunk[:chunkSize]
		}

		stream.Send(&pb.AddPhotoRequest{Data: chunk})

	}

	res, err := stream.CloseAndRecv()
	if err != nil {
		log.Fatalln("--- 1" + err.Error())
	}

	fmt.Println("add photo isok:", res.IsOk)
}

func saveAll(client pb.EmployeeServiceClient) {
	employees := []pb.Employee{
		{
			Id:   1,
			No:   23,
			Name: "lx",
			MonthSalary: &pb.MonthSalary{
				Basic: 200.0,
				Bonus: 23233.0,
			},
			Status: pb.EmployeeStatus_NORMAL,
			LastModified: &timestamppb.Timestamp{
				Seconds: time.Now().Unix(),
			},
		},
		{
			Id:   2,
			No:   23,
			Name: "wzz",
			MonthSalary: &pb.MonthSalary{
				Basic: 200.0,
				Bonus: 23233.0,
			},
			Status: pb.EmployeeStatus_NORMAL,
			LastModified: &timestamppb.Timestamp{
				Seconds: time.Now().Unix(),
			},
		},
	}

	stream, err := client.SaveAll(context.Background())
	if err != nil {
		log.Fatalln(err.Error())
	}

	finishChannel := make(chan struct{})
	go func() {
		for {
			res, err := stream.Recv()
			if err == io.EOF {
				finishChannel <- struct{}{}
				break
			}
			if err != nil {
				log.Fatalln(err.Error())
			}
			fmt.Println(res.Employee)
		}
	}()

	for _, e := range employees {
		err := stream.Send(&pb.SaveEmployeeRequest{Employee: &e})
		if err != nil {
			log.Fatalln(err.Error())
		}
	}
	stream.CloseSend()
	<-finishChannel
}
