package main

import (
	"context"
	"log"
	"net"
	"net/http"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	gw "com.example.test/pb"
)

type myServiceImpl struct {
	gw.UnimplementedMyServiceServer
}

func (*myServiceImpl) Echo(ctx context.Context, req *gw.StrMsg) (*gw.StrMsg, error) {
	log.Println("recv msg :" + req.Msg + ",resp : hello " + req.Msg)
	return &gw.StrMsg{Msg: "hello " + req.Msg}, nil
}

func main() {
	s := grpc.NewServer()
	gw.RegisterMyServiceServer(s, &myServiceImpl{})
	lis, err := net.Listen("tcp", ":9099")
	if err != nil {
		panic(err)
	}

	log.Println("Serving gRPC on 0.0.0.0:9099...")
	go func() {
		if err := s.Serve(lis); err != nil {
			log.Fatalln("Failed to serve grpc, err:", err.Error())
		}
	}()

	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)

	defer cancel()

	mux := runtime.NewServeMux()
	opts := []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}
	err = gw.RegisterMyServiceHandlerFromEndpoint(ctx, mux, "127.0.0.1:9099", opts)
	if err != nil {
		panic(err)
	}

	log.Println("Serving http on 0.0.0.0:8088...")
	http.ListenAndServe(":8088", mux)
}
