package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"os"
	"time"

	"github.com/lxxxxxxx/grpc-client/client"
	"github.com/lxxxxxxx/grpc-client/pb"

	"google.golang.org/grpc/metadata"
	"google.golang.org/protobuf/types/known/timestamppb"
)

const port = ":5001"

func main() {

	client.InitClient(context.Background(), port)

	getByNo(client.EmployeeService)
	getAll(client.EmployeeService)
	addPhoto(client.EmployeeService)
	saveAll(client.EmployeeService)

	getPerson(client.PersonService)
}

func getByNo(client pb.EmployeeClient) {
	res, err := client.GetByNo(context.Background(), &pb.GetByNoRequest{No: 1999})
	if err != nil {
		log.Fatalln(err.Error())
	}
	fmt.Println(res.Employee)
}

func getAll(client pb.EmployeeClient) {
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

func addPhoto(client pb.EmployeeClient) {
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

func saveAll(client pb.EmployeeClient) {
	employees := []pb.EmployeeInfo{
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

func getPerson(client pb.PersonClient) {
	person, _ := client.GetPerson(context.Background(), &pb.GetPersonReq{Name: "lixiang"})
	fmt.Println(person)
}
