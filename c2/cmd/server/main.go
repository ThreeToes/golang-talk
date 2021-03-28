package main

import (
	"flag"
	"fmt"
	"github.com/ThreeToes/golang-talk/c2"
	"github.com/ThreeToes/golang-talk/c2/gen"
	"google.golang.org/grpc"
	"log"
	"net"
)

const CryptoKey = "SuperSecret"

func main() {
	imageFolder := flag.String("images", "", "Folder with images")
	bindAddr := flag.String("addr", "0.0.0.0", "Interface to bind on")
	portF := flag.Int("port", 5555, "Port to bind on")
	flag.Parse()
	if *imageFolder == "" { // OMIT
		log.Fatal("Must set -images") // OMIT
	} // OMIT
	log.Printf("Using images in %s", *imageFolder)
	svr := c2.NewPictureServer(*imageFolder, CryptoKey)

	clientSocket, err := net.Listen("tcp", fmt.Sprintf("%s:%d",*bindAddr, *portF))
	if err != nil { // OMIT
		log.Fatalf("failed to listen: %v", err) // OMIT
	} // OMIT
	var opts []grpc.ServerOption
	grpcServer := grpc.NewServer(opts...)
	gen.RegisterPictureSharingServer(grpcServer, svr)
	gen.RegisterMemeDealerServer(grpcServer, svr)
	log.Printf("Starting server on %s:%d", *bindAddr, *portF)
	grpcServer.Serve(clientSocket)
}