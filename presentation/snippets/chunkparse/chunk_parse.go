package main

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"os"
)

func main() {
	ReadChunk()
}

func ReadChunk() {
	b, _ := os.ReadFile("./images/png_hex.png")
	r := bytes.NewReader(b)
	var header uint64
	binary.Read(r, binary.BigEndian, &header)
	var chunkSize uint32
	binary.Read(r, binary.BigEndian, &chunkSize)
	var chunkType uint32
	binary.Read(r, binary.BigEndian, &chunkType)
	typeString := make([]byte, 4)
	binary.BigEndian.PutUint32(typeString, chunkType)
	r.Seek(int64(chunkSize), 1)
	var crc uint32
	binary.Read(r, binary.BigEndian, &crc)
	fmt.Printf("Header: %x\nSize: %d\nType: %s\nCRC: %x",
		header, chunkSize, string(typeString), crc)
}
