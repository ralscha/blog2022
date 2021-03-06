package main

import (
	"bufio"
	"encoding/hex"
	"fmt"
	"github.com/bits-and-blooms/bloom/v3"
	"log"
	"os"
	"strconv"
)

func main() {
	file, err := os.Open("./pwned-passwords-sha1-ordered-by-hash-v8.txt")
	if err != nil {
		log.Fatal("Can't open pwned file", err)
	}
	defer file.Close()

	const nPasswords = 847_223_402

	filters := []*bloom.BloomFilter{
		bloom.NewWithEstimates(nPasswords, 0.1),
		bloom.NewWithEstimates(nPasswords, 0.01),
		bloom.NewWithEstimates(nPasswords, 0.001),
		bloom.NewWithEstimates(nPasswords, 0.0001),
		bloom.NewWithEstimates(nPasswords, 0.00001),
		bloom.NewWithEstimates(nPasswords, 0.000001),
		bloom.NewWithEstimates(nPasswords, 0.0000001),
		bloom.NewWithEstimates(nPasswords, 0.00000001),
	}

	count := 0
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		hexString := line[0:40]

		decodedHex, err := hex.DecodeString(hexString)
		if err != nil {
			log.Fatalf("Failed to decode hex string %s, %v", hexString, err)
		}
		count++
		for _, filter := range filters {
			filter.Add(decodedHex)
		}

		if count%8_500_000 == 0 {
			fmt.Println(count*100/nPasswords, "%")
		}
	}

	for ix, filter := range filters {
		fmt.Println("count: ", filter.BitSet().Count())
		encode, err := filter.GobEncode()
		if err != nil {
			log.Fatal("gob encode failed", err)
		}

		file, err = os.Create("./pwned-passwords-sha1-ordered-by-hash-v8_" + strconv.Itoa(ix+1) + ".gob")
		if err != nil {
			log.Fatal("Can't create bloom file", err)
		}

		_, err = file.Write(encode)
		if err != nil {
			log.Fatal("Can't write bloom file", err)
		}

		fmt.Println("bytes: ", len(encode))
		file.Close()
	}

}
