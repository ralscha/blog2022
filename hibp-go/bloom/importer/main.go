package main

import (
	"bufio"
	"encoding/hex"
	"fmt"
	"github.com/bits-and-blooms/bloom/v3"
	"log"
	"os"
	"sort"
	"strconv"
	"strings"
)

func main() {
	inputDir := "E:\\_download\\pwned" // "./pwned"

	files, err := os.ReadDir(inputDir)
	if err != nil {
		log.Fatalf("Can't read input directory %v", err)
	}
	var fileNames []string
	for _, file := range files {
		fileNames = append(fileNames, file.Name())
	}
	sort.Strings(fileNames)

	// count number of password hashes
	var nPasswords uint
	for _, fileName := range fileNames {
		file, err := os.Open(inputDir + "/" + fileName)
		if err != nil {
			log.Fatalf("Can't open file %s, %v", fileName, err)
		}
		scanner := bufio.NewScanner(file)
		for scanner.Scan() {
			nPasswords++
		}
		err = file.Close()
		if err != nil {
			log.Fatalf("Can't close file %s, %v", fileName, err)
		}
	}

	fmt.Println("nPasswords: ", nPasswords)

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

	var count uint

	for _, fileName := range fileNames {
		file, err := os.Open(inputDir + "/" + fileName)
		if err != nil {
			log.Fatalf("Can't open file %s, %v", fileName, err)
		}
		hashPrefix := strings.Split(fileName, ".")[0]
		scanner := bufio.NewScanner(file)
		for scanner.Scan() {
			line := scanner.Text()
			line = hashPrefix + line
			hexString := line[0:40]

			decodedHex, err := hex.DecodeString(hexString)
			if err != nil {
				log.Fatalf("Failed to decode hex string %s, %v", hexString, err)
			}
			count++
			for _, filter := range filters {
				filter.Add(decodedHex)
			}

			if count%10_000_000 == 0 {
				fmt.Println(count*100/nPasswords, "%")
			}
		}
	}

	for ix, filter := range filters {
		fmt.Println("count: ", filter.BitSet().Count())
		encode, err := filter.GobEncode()
		if err != nil {
			log.Fatal("gob encode failed", err)
		}

		file, err := os.Create("./pwned-passwords-sha1-ordered-by-hash-v8_" + strconv.Itoa(ix+1) + ".gob")
		if err != nil {
			log.Fatal("Can't create bloom file", err)
		}

		_, err = file.Write(encode)
		if err != nil {
			log.Fatal("Can't write bloom file", err)
		}

		fmt.Println("bytes: ", len(encode))
		err = file.Close()
		if err != nil {
			log.Fatal("Can't close bloom file", err)
		}
		fmt.Println()
	}

}
