package c2

import (
	"context"
	"github.com/ThreeToes/golang-talk/c2/gen"
	"github.com/ThreeToes/golang-talk/encrypter"
	"github.com/ThreeToes/golang-talk/pnghider"
	"github.com/google/uuid"
	"io/ioutil"
	"log"
	"math/rand"
	"os"
	"path"
	"sync"
)

const payloadType = "sNKY"

type PictureServer struct {
	gen.PictureSharingServer
	gen.MemeDealerServer
	toDishOut []*gen.InnocentPicture
	dishLock *sync.Mutex
	returned map[string] []byte
	imagesFolder string
	cryptoKey string
}

func (p *PictureServer) DishMeme(_ context.Context,
		parameters *gen.DishMemeParamaters) (*gen.DishMemeResponse, error) {
	// Pick the file to hide our message
	fs, err := ioutil.ReadDir(p.imagesFolder)
	if err != nil { // OMIT
		return nil, err // OMIT
	}// OMIT
	toOpen := fs[rand.Intn(len(fs))]
	pic, err := os.ReadFile(path.Join(p.imagesFolder, toOpen.Name()))
	if err != nil {// OMIT
		return nil, err// OMIT
	}// OMIT
	log.Printf("Hiding payload %s", parameters.Payload)

	// Encrypt the payload
	encryptedPayload, err := encrypter.EncryptData(p.cryptoKey, []byte(parameters.Payload))
	if err != nil {// OMIT
		return nil, err// OMIT
	}// OMIT
	// Hide the payload
	hidden, err := pnghider.HidePayload([]byte(payloadType), encryptedPayload, pic)
	if err != nil {// OMIT
		return nil, err// OMIT
	}// OMIT
	id := uuid.New().String()
	sneakyPayload := &gen.InnocentPicture{
		Id:   id,
		Data: hidden,
	}

	// Put it aside until the client is ready to pick it up
	p.dishLock.Lock()
	p.toDishOut = append(p.toDishOut, sneakyPayload)
	p.dishLock.Unlock()
	return &gen.DishMemeResponse{Id: id}, nil
}

func (p *PictureServer) GetMemeStatus(ctx context.Context, parameters *gen.CheckMemeStatusParameters) (*gen.CheckMemeStatusResponse, error) {
	v, ok := p.returned[parameters.Id]
	if !ok {
		return &gen.CheckMemeStatusResponse{Status: "unready"}, nil
	}
	return &gen.CheckMemeStatusResponse{Status: "ready", Response: v}, nil
}

func (p *PictureServer) mustEmbedUnimplementedMemeDealerServer() {
}

func (p *PictureServer) GetPicture(ctx context.Context, parameters *gen.GetPictureParameters) (*gen.InnocentPicture, error) {
	p.dishLock.Lock()
	if len(p.toDishOut) == 0 {
		p.dishLock.Unlock()
		return &gen.InnocentPicture{Id: "", Data: nil}, nil
	}
	pic := p.toDishOut[0]
	if len(p.toDishOut) == 1 {
		p.toDishOut = nil
	} else {
		p.toDishOut = p.toDishOut[1:]
	}
	p.dishLock.Unlock()
	return pic, nil
}

func (p *PictureServer) SayThankyou(ctx context.Context, thankyou *gen.Thankyou) (*gen.ThankyouOutput, error) {
	payload, err := pnghider.RecoverPayload([]byte(payloadType), thankyou.AnotherPicture)
	if err != nil {
		return nil, err
	}
	p.returned[thankyou.Id] = payload
	return &gen.ThankyouOutput{}, nil
}

func (p *PictureServer) mustEmbedUnimplementedPictureSharingServer() {
}

func NewPictureServer(imageFolder string, cryptoKey string) *PictureServer {
	return &PictureServer{
		toDishOut:            nil,
		dishLock:             &sync.Mutex{},
		returned: 			  map[string][]byte{},
		imagesFolder:         imageFolder,
		cryptoKey: cryptoKey,
	}
}