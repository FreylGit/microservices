package main

import (
	"context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	desc "grpc/pkg/note_v1"
	"log"
	"time"
)

const address = "localhost:50051"
const noteID = 12

func main() {
	conn, err := grpc.Dial(address, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatal("failed to connect to server")
	}
	defer conn.Close()

	c := desc.NewNoteV1Client(conn)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	r, err := c.Get(ctx, &desc.GetRequest{Id: noteID})
	if err != nil {
		log.Fatal("failed to get note by id")
	}
	log.Printf("note info:\n %+v", r.GetNote())
}
