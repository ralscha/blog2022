package main

import (
	"crypto/sha1"
	"github.com/bits-and-blooms/bloom/v3"
	"io/ioutil"
	"log"
	"os"
)

func main() {
	file, err := os.Open("./pwned-passwords-sha1-ordered-by-hash-v8_1.gob")
	if err != nil {
		log.Fatal("Can't open pwned file", err)
	}
	defer file.Close()

	const nPasswords = 847_223_402
	const falsePositiveRate = 0.1
	filter := bloom.NewWithEstimates(nPasswords, falsePositiveRate)

	bytes, err := ioutil.ReadAll(file)
	if err != nil {
		log.Fatal("Can't read pwned file", err)
	}

	err = filter.GobDecode(bytes)
	if err != nil {
		log.Fatal("Can't decode bytes", err)
	}

	const password = "123456"

	h := sha1.New()
	h.Write([]byte(password))
	sha1hash := h.Sum(nil)

	if filter.Test(sha1hash) {
		log.Println("Password is in the filter")
	} else {
		log.Println("Password is not in the filter")
	}

}
