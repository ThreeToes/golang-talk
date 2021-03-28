package main

import (
	"flag"
	"fmt"
	"github.com/ThreeToes/golang-talk/encrypter"
	"os"
)

func main() {
	encryptF := flag.Bool("encrypt", false, "Encrypt data")
	decryptF := flag.Bool("decrypt", false, "Decrypt data")
	keyF := flag.String("key", "", "Encryption key to use")
	inF := flag.String("file", "", "File to operate on")

	flag.Parse()
	if *inF == "" {
		fmt.Println("-file flag must be set")
		flag.Usage()
		return
	}
	if *keyF == "" {
		fmt.Println("-key flag must be set")
		flag.Usage()
		return
	}
	if (!*encryptF && !*decryptF) || (*encryptF && *decryptF) {
		fmt.Println("Must specify one of -encrypt or -decrypt")
		flag.Usage()
		return
	}

	if *decryptF {
		payload, err := os.ReadFile(*inF)
		if err != nil {
			fmt.Printf("Could not read file %s: %v\n", *inF, err)
			return
		}
		plainText, err := encrypter.DecryptData(*keyF, payload)
		if err != nil {
			fmt.Printf("Could not decrypt file %s: %v\n", *inF, err)
			return
		}
		err = os.WriteFile(*inF, plainText, 0644)
		if err != nil {
			fmt.Printf("Could not overwrite file %s: %v\n", *inF, err)
			return
		}
	}else if *encryptF {
		payload, err := os.ReadFile(*inF)
		if err != nil {
			fmt.Printf("Could not read file %s: %v\n", *inF, err)
			return
		}
		plainText, err := encrypter.EncryptData(*keyF, payload)
		if err != nil {
			fmt.Printf("Could not encrypt file %s: %v\n", *inF, err)
			return
		}
		err = os.WriteFile(*inF, plainText, 0644)
		if err != nil {
			fmt.Printf("Could not overwrite file %s: %v\n", *inF, err)
			return
		}
	}
	fmt.Printf("Wrote out file %s\n", *inF)
}
