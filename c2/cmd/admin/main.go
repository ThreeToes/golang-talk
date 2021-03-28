package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/ThreeToes/golang-talk/c2/gen"
	"google.golang.org/grpc"
	"log"
	"os"
	"time"
)

func main() {
	addrF :=  flag.String("addr", "127.0.0.1", "Server address") // OMIT
	portF := flag.Int("port", 5555, "Server port") // OMIT
	payloadF := flag.String("payload", "", "Payload to send to client")
	outF := flag.String("out", "", "Path to write output to")
	flag.Parse()

	toDial := fmt.Sprintf("%s:%d", *addrF, *portF)
	log.Printf("Connecting to %s", toDial) // OMIT
	log.Printf("payload=%s", *payloadF) // OMIT
	conn, err := grpc.Dial(toDial, grpc.WithInsecure())
	if err != nil { // OMIT
		log.Fatal("Could not dial server") // OMIT
	}// OMIT
	client := gen.NewMemeDealerClient(conn)
	ret, err := client.DishMeme(context.TODO(), &gen.DishMemeParamaters{Payload: *payloadF})
	if err != nil {// OMIT
		log.Fatalf("Error dishing out teh meme: %v", err)// OMIT
	}// OMIT
	id := ret.Id
	for {
		clientResp, err := client.GetMemeStatus(context.TODO(), &gen.CheckMemeStatusParameters{Id: id})
		if err != nil {// OMIT
			log.Fatalf("Error trying to get the meme status: %v", err)// OMIT
		}// OMIT
		if clientResp.Status == "unready" {
			time.Sleep(500 * time.Millisecond)
			continue
		}
		if *outF == "" {
			log.Printf("Client returned '%s'", string(clientResp.Response))
		} else {
			os.WriteFile(*outF, clientResp.Response, 0644)
		}
		break
	}
}
