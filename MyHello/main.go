package main

import (
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"log"
	"net"
	"net/http"
	"strings"
	edu_pb "MyHello/protos"
	services "MyHello/services"
)
type server struct{}
func main() {
	lis, err := net.Listen("tcp", ":10003")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	//pb.RegisterGreeterServer(s, &server{})

	var con services.ConfigService
	edu_pb.RegisterGreeterServer(s,con)
	// Register reflection service on gRPC server.
	reflection.Register(s)
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}

}


func grpcHandleFunc(grpcServer *grpc.Server, otherHander http.Handler) http.Handler {
	return h2c.NewHandler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.ProtoMajor == 2 && strings.Contains(r.Header.Get("Content-Type"), "application/grpc") {
			grpcServer.ServeHTTP(w, r)
		} else {
			otherHander.ServeHTTP(w, r)
		}
	}), &http2.Server{})
}
