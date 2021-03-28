package pnghider

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"hash/crc32"
	"io/ioutil"
)

// PNG looks like:
// 8 bytes for PNG header
// successive chunks

const (
	pngHeaderHex = 0x89504e470d0a1a0a
	endChunkType = "IEND"
)

type pngChunk struct {
	size      uint32
	chunkType []byte
	data      []byte
	crc       uint32
	offset    int64
}

// confirmPng confirms that a file is a png. Note this advances the stream reader
func confirmPng(b *bytes.Reader) error {
	var header uint64

	// Try and read the first 8 bytes into our header struct
	if err := binary.Read(b, binary.BigEndian, &header); err != nil {
		return err
	}

	// Check to make sure the header is a PNG
	if header != pngHeaderHex {
		return fmt.Errorf("file is not in png format")
	}
	return nil
}

// Chunk structure
//  - 4 bytes for chunk size in bytes
//  - 4 bytes for chunk type
//  - chunk data, size is determined by 4 size bytes
//  - 4 bytes for Cyclic Redundancy Check (CRC)
func loadChunk(b *bytes.Reader) (*pngChunk, error) {
	var chunk pngChunk
	var err error
	chunk.offset, err = b.Seek(0, 1)
	if err != nil {
		return nil, fmt.Errorf("error reading chunk offset: %v", err)
	}
	err = binary.Read(b, binary.BigEndian, &chunk.size)
	if err != nil {
		return nil, fmt.Errorf("error reading chunk at offset %d's size: %v", chunk.offset, err)
	}
	var chunkType uint32
	err = binary.Read(b, binary.BigEndian, &chunkType)
	if err != nil {
		return nil, fmt.Errorf("error reading chunk type: %v", err)
	}
	chunk.chunkType = make([]byte, 4)
	binary.BigEndian.PutUint32(chunk.chunkType, chunkType)
	chunk.data = make([]byte, chunk.size)
	bytesRead, err := b.Read(chunk.data)
	if err != nil {
		return nil, fmt.Errorf("error reading chunk data: %v", err)
	} else if uint32(bytesRead) != chunk.size {
		return nil, fmt.Errorf(
			"chunk bytes read was different from declared size. expected: %d, actual: %d", chunk.size, bytesRead)
	}
	err = binary.Read(b, binary.BigEndian, &chunk.crc)
	return &chunk, nil
}

func calculateCrc(chunkType, data []byte) (uint32, error) {
	var buffer bytes.Buffer
	if err := binary.Write(&buffer, binary.BigEndian, chunkType); err != nil {
		return 0, err
	}
	if err := binary.Write(&buffer, binary.BigEndian, data); err != nil {
		return 0, err
	}
	return crc32.ChecksumIEEE(buffer.Bytes()), nil
}

func loadChunks(b *bytes.Reader) ([]*pngChunk, error) {
	var chunks []*pngChunk
	for {
		chunk, err := loadChunk(b)
		if err != nil {
			return chunks, fmt.Errorf("error loading chunks: %v", err)
		}
		chunks = append(chunks, chunk)
		if string(chunk.chunkType) == endChunkType {
			break
		}
	}
	return chunks, nil
}

func createChunk(typeString, payload []byte) (*pngChunk, error) {
	var chunk pngChunk
	if len(typeString) != 4 {
		// Can we pad it if it's less than 4 bytes? Probably
		return nil, fmt.Errorf("type string must be 4 bytes long")
	}
	chunk.chunkType = typeString
	chunk.data = payload
	chunk.size = uint32(len(payload))
	var err error
	chunk.crc, err = calculateCrc(typeString, payload)
	return &chunk, err
}

func chunksToPngBytes(chunks []*pngChunk) ([]byte, error) {
	var buff bytes.Buffer
	header := make([]byte, 8)
	binary.BigEndian.PutUint64(header, pngHeaderHex)
	if err := binary.Write(&buff, binary.BigEndian, header); err != nil {
		return nil, err
	}
	for _, chunk := range chunks {
		if err := binary.Write(&buff, binary.BigEndian, chunk.size); err != nil {
			return nil, fmt.Errorf("error writing chunk size to buffer: %v", err)
		}
		if err := binary.Write(&buff, binary.BigEndian, chunk.chunkType); err != nil {
			return nil, fmt.Errorf("error writing chunk type '%s' to buffer: %v", string(chunk.chunkType), err)
		}
		if err := binary.Write(&buff, binary.BigEndian, chunk.data); err != nil {
			return nil, fmt.Errorf("error writing chunk data to buffer: %v", err)
		}
		if err := binary.Write(&buff, binary.BigEndian, chunk.crc); err != nil {
			return nil, fmt.Errorf("error writing chunk CRC to buffer: %v", err)
		}
	}

	return buff.Bytes(), nil
}

// HidePayload Places a payload into a provided PNG file
func HidePayload(typeString, payload []byte, pngPath string) ([]byte, error) {
	pic, err := ioutil.ReadFile(pngPath)
	if err != nil {
		return nil, err
	}
	picReader := bytes.NewReader(pic)
	err = confirmPng(picReader)
	if err != nil {
		return nil, err
	}
	chunks, err := loadChunks(picReader)
	if err != nil {
		return nil, err
	}
	newChunk, err := createChunk(typeString, payload)
	if err != nil {
		return nil, err
	}
	// TODO: do something trickier here maybe?
	chunks = append([]*pngChunk{newChunk}, chunks...)
	return chunksToPngBytes(chunks)
}

func ReturnPayload(pngPath string) ([]byte, error) {
	return nil, nil
}
