package main

import (
	"bufio"
	"encoding/binary"
	"encoding/hex"
	"log"
	"os"
	"strconv"

	"github.com/cockroachdb/pebble"
)

func main() {

	file, err := os.Open("./pwned-passwords-sha1-ordered-by-hash-v8.txt")
	if err != nil {
		log.Fatalf("Can't open hibp file %v", err)
	}
	defer file.Close()

	path := "/var/lib/hibp/pebble"
	err = os.MkdirAll(path, os.ModePerm)
	if err != nil {
		log.Fatalf("Can't create data directory %v", err)
	}

	db, err := pebble.Open(path, &pebble.Options{DisableWAL: true})
	if err != nil {
		log.Fatalf("Can't open database %v", err)
	}
	defer db.Close()

	counter := 0
	batch := db.NewBatch()

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
		err = batch.Set(decodedHex, appearsBytes, pebble.NoSync)
		if err != nil {
			log.Fatalf("Set value into database failed %v", err)
		}
		counter++
		if counter > 5_000_000 {
			err := batch.Commit(pebble.NoSync)
			if err != nil {
				log.Fatalf("Commit failed %v", err)
			}
			batch = db.NewBatch()
			counter = 0
		}

	}

	err = batch.Commit(pebble.NoSync)
	if err != nil {
		log.Fatalf("Commit failed %v", err)
	}

	if err := scanner.Err(); err != nil {
		log.Fatalf("Scanner failed %v", err)
	}

}
