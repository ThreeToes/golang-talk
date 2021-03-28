package main

import (
	"flag"
	"fmt"
	"github.com/ThreeToes/golang-talk/pnghider"
	"os"
	"unicode"
)

func main() {
	decodeF := flag.Bool("decode", false, "Decode payload from file")
	inputF := flag.String("in", "", "Input file")
	outputF := flag.String("out", "", "Output file")
	typeF := flag.String("type", "sNKY", "4 byte type string")
	payloadF := flag.String("payload", "", "Payload to hide")
	flag.Parse()

	if *inputF == "" {
		fmt.Printf("Must set the -in flag\n")
		flag.Usage()
		return
	}
	if !*decodeF {
		if *outputF == "" {
			fmt.Printf("When encoding, you must set the -out flag")
			flag.Usage()
			return
		} else if *payloadF == "" {
			fmt.Printf("When encoding, you must set the -payload flag")
			flag.Usage()
			return
		}
	}
	if len(*typeF) != 4 || unicode.IsUpper([]rune(*typeF)[0]) {
		fmt.Printf("Type string must be of length 4 and start with a lowercase letter\n")
		flag.Usage()
		return
	}
	if *decodeF {
		pic, err := os.ReadFile(*inputF)
		if err != nil {
			fmt.Printf("Error: Could not read input file %s: %v\n", *inputF, err)
			return
		}
		out, err := pnghider.RecoverPayload([]byte(*typeF), pic)
		if err != nil {
			fmt.Printf("Error: Could not write output file %s: %v\n", *outputF, err)
			return
		}
		fmt.Printf("Recovered payload: %s\n", string(out))
	} else {
		pic, err := os.ReadFile(*inputF)
		if err != nil {
			fmt.Printf("Error: Could not read input file %s: %v\n", *inputF, err)
			return
		}
		out, err := pnghider.HidePayload([]byte(*typeF), []byte(*payloadF), pic)
		err = os.WriteFile(*outputF, out, 0644)
		if err != nil {
			fmt.Printf("Error: Could not write output file %s: %v\n", *outputF, err)
			return
		}
		fmt.Println("Encoded file successfully")
	}
}
