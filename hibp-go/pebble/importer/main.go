package main

import (
	"bufio"
	"encoding/binary"
	"encoding/hex"
	"log"
	"os"
	"sort"
	"strconv"
	"strings"

	"github.com/cockroachdb/pebble"
)

func main() {

	inputDir := "./pwned"

	files, err := os.ReadDir(inputDir)
	if err != nil {
		log.Fatalf("Can't read input directory %v", err)
	}
	var fileNames []string
	for _, file := range files {
		fileNames = append(fileNames, file.Name())
	}
	sort.Strings(fileNames)

	path := "/var/lib/hibp/pebble"
	err = os.MkdirAll(path, os.ModePerm)
	if err != nil {
		log.Fatalf("Can't create data directory %v", err)
	}

	db, err := pebble.Open(path, &pebble.Options{DisableWAL: true})
	if err != nil {
		log.Fatalf("Can't open database %v", err)
	}
	defer func(db *pebble.DB) {
		err := db.Close()
		if err != nil {
			log.Fatalf("Can't close database %v", err)
		}
	}(db)

	for _, fileName := range fileNames {
		file, err := os.Open(inputDir + "/" + fileName)

		batch := db.NewBatch()

		hashPrefix := strings.Split(fileName, ".")[0]

		scanner := bufio.NewScanner(file)
		for scanner.Scan() {
			line := scanner.Text()
			line = hashPrefix + line
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
		}

		err = batch.Commit(pebble.NoSync)
		if err != nil {
			log.Fatalf("Commit failed %v", err)
		}

		if err := scanner.Err(); err != nil {
			log.Fatalf("Scanner failed %v", err)
		}
	}

	err = db.Flush()
	if err != nil {
		log.Fatalf("Flush failed %v", err)
	}

}
