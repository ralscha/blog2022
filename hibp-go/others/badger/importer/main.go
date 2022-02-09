package main

import (
	"bufio"
	"encoding/binary"
	"encoding/hex"
	"log"
	"os"
	"strconv"

	badger "github.com/dgraph-io/badger/v3"
	// "github.com/dgraph-io/badger/v3/options"
)

func main() {

	file, err := os.Open("./pwned-passwords-sha1-ordered-by-hash-v8.txt")
	if err != nil {
		log.Fatalf("Can't open file %v", err)
	}
	defer file.Close()

	path := "/var/lib/hibp/badger"
	err = os.MkdirAll(path, os.ModePerm)
	if err != nil {
		log.Fatalf("Can't make output directory %v", err)
	}

	db, err := badger.Open(badger.DefaultOptions(path)) // .WithCompression(options.ZSTD))
	if err != nil {
		log.Fatalf("Can't open database %v", err)
	}
	defer db.Close()

	wb := db.NewWriteBatch()
	defer wb.Cancel()

	//read file line by line
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		//extract sha1 and appears value
		line := scanner.Text()
		hexString := line[0:40]
		appears, err := strconv.ParseUint(line[41:], 10, 32)
		if err != nil {
			log.Fatalf("String to int conversion failed %v", err)
		}

		decodedHex, err := hex.DecodeString(hexString)
		if err != nil {
			log.Fatalf("Failed to decode hex string %v", err)
		}

		buf := make([]byte, binary.MaxVarintLen64)
		n := binary.PutUvarint(buf, appears)
		appearsBytes := buf[:n]

		// insert into database
		err = wb.Set(decodedHex, appearsBytes)
		if err != nil {
			log.Fatalf("Set value into database failed %v", err)
		}
	}

	if err := wb.Flush(); err != nil {
		log.Fatalf("Flush error %v", err)
	}

	if err := scanner.Err(); err != nil {
		log.Fatalf("Scanner failed %v", err)
	}

}
