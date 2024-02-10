package main

import (
	"context"
	"fmt"
	"github.com/brianvoe/gofakeit/v6"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"google.golang.org/protobuf/types/known/timestamppb"
	desc "grpc/pkg/note_v1"
	"log"
	"net"
)

const grpcPort = 50051

type server struct {
	desc.UnimplementedNoteV1Server
}

func main() {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", grpcPort))
	if err != nil {
		log.Fatalf("failed to lesten:%v", err)
	}
	s := grpc.NewServer()
	reflection.Register(s)
	desc.RegisterNoteV1Server(s, &server{})

	log.Println("server listening at ")
	if err := s.Serve(lis); err != nil {
		log.Fatal("Failed to server")
	}
}

func (s *server) Get(ctx context.Context, req *desc.GetRequest) (*desc.GetResponse, error) {
	log.Printf("note id: %d", req.GetId())

	return &desc.GetResponse{
		Note: &desc.Note{
			Id: req.GetId(),
			Info: &desc.NoteInfo{
				Title:    gofakeit.BeerName(),
				Content:  gofakeit.IPv6Address(),
				Author:   gofakeit.Name(),
				IsPublic: gofakeit.Bool(),
			},
			CreatedAt: timestamppb.New(gofakeit.Date()),
			UpdatedAt: timestamppb.New(gofakeit.Date()),
		},
	}, nil
}
