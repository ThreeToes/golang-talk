package pnghider

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"hash/crc32"
	"strings"
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

func (chunk *pngChunk) marshalData(buff *bytes.Buffer) error {
	if err := binary.Write(buff, binary.BigEndian, chunk.size); err != nil {
		return fmt.Errorf("error writing chunk size to buffer: %v", err)
	}
	if err := binary.Write(buff, binary.BigEndian, chunk.chunkType); err != nil {
		return fmt.Errorf("error writing chunk type '%s' to buffer: %v", string(chunk.chunkType), err)
	}
	if err := binary.Write(buff, binary.BigEndian, chunk.data); err != nil {
		return fmt.Errorf("error writing chunk data to buffer: %v", err)
	}
	if err := binary.Write(buff, binary.BigEndian, chunk.crc); err != nil {
		return fmt.Errorf("error writing chunk CRC to buffer: %v", err)
	}
	return nil
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

// calculateCrc calculates the cyclic redundancy check
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

// loadChunks parses the binary data into the pngChunk type.
func loadChunks(b *bytes.Reader) ([]*pngChunk, error) {
	if n, _ := b.Seek(0, 1); n != 8 {
		return nil, fmt.Errorf("stream must be at the first chunk after the initial 8 size bytes")
	}
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

// createChunk creates a sneaky chunk for the payload
func createChunk(typeBytes, payload []byte) (*pngChunk, error) {
	var chunk pngChunk
	typeString := string(typeBytes)
	if len(typeBytes) != 4 {
		// Can we pad it if it's less than 4 bytes? Probably
		return nil, fmt.Errorf("type string must be 4 bytes long")
	} else if typeString[0] == strings.ToUpper(typeString)[0] {
		return nil, fmt.Errorf("cannot use chunk type %s, leading capital letter marks it as critical", typeString)
	}
	chunk.chunkType = typeBytes
	chunk.data = payload
	chunk.size = uint32(len(payload))
	var err error
	chunk.crc, err = calculateCrc(typeBytes, payload)
	return &chunk, err
}

// chunksToPngBytes generates a byte array from a list of PNG chunks
func chunksToPngBytes(chunks []*pngChunk) ([]byte, error) {
	var buff bytes.Buffer
	header := make([]byte, 8)
	binary.BigEndian.PutUint64(header, pngHeaderHex)
	if err := binary.Write(&buff, binary.BigEndian, header); err != nil {
		return nil, err
	}
	for _, chunk := range chunks {
		if err := chunk.marshalData(&buff); err != nil {
			return nil, err
		}
	}

	return buff.Bytes(), nil
}


func unmarshalData(pic []byte) ([]*pngChunk, error) {
	picReader := bytes.NewReader(pic)
	err := confirmPng(picReader)
	if err != nil {
		return nil, err
	}
	chunks, err := loadChunks(picReader)
	if err != nil {
		return nil, err
	}
	return chunks, nil
}

// HidePayload Places a payload into a provided PNG file
func HidePayload(typeString, payload, pic []byte) ([]byte, error) {
	chunks, err := unmarshalData(pic)
	if err != nil {
		return nil, err
	}
	newChunk, err := createChunk(typeString, payload)
	if err != nil {
		return nil, err
	}
	// Put our chunk second to last
	chunks = append(chunks[0:len(chunks)-1], newChunk, chunks[len(chunks)-1])
	return chunksToPngBytes(chunks)
}

func RecoverPayload(typeBytes, pic []byte) ([]byte, error) {
	chunks, err := unmarshalData(pic)
	if err != nil {
		return nil, err
	}
	sneakyChunk := chunks[len(chunks)-2]
	chunkType := string(sneakyChunk.chunkType)
	typeString := string(typeBytes)
	if chunkType != typeString {
		return nil, fmt.Errorf(
			"type string did not match. expected: %s, actual: %s", typeString, chunkType)
	}
	return sneakyChunk.data, nil
}
