package main

import (
	"bufio"
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"log"
	"os"
	"strconv"

	bolt "go.etcd.io/bbolt"
)

func main() {

	file, err := os.Open("./pwned-passwords-sha1-ordered-by-hash-v8.txt")
	if err != nil {
		log.Fatalf("Can't open file %v", err)
	}
	defer file.Close()

	path := "/var/lib/hibp/bbolt"

	db, err := bolt.Open(path, 0666, nil)
	if err != nil {
		log.Fatalf("Can't open database %v", err)
	}
	defer db.Close()

	//create bucket
	err = db.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucket([]byte("hibp"))
		if err != nil {
			return fmt.Errorf("create bucket: %s", err)
		}
		return nil
	})
	if err != nil {
		log.Fatalf("Can't create bucket %v", err)
	}

	counter := 0
	txn, err := db.Begin(true)
	if err != nil {
		log.Fatalf("Can't start transaction %v", err)
	}

	bucket := txn.Bucket([]byte("hibp"))

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
		err = bucket.Put(decodedHex, appearsBytes)
		if err != nil {
			txn.Rollback()
			log.Fatalf("Set value into database failed %v", err)
		}

		if counter > 100_000 {
			if err := txn.Commit(); err != nil {
				log.Fatalf("Commit error %v", err)
			}
			counter = 0

			txn, err = db.Begin(true)
			bucket = txn.Bucket([]byte("hibp"))
			if err != nil {
				log.Fatalf("Can't start transaction %v", err)
			}
		}

		counter++
	}

	if err := txn.Commit(); err != nil {
		log.Fatalf("Commit error %v", err)
	}

	if err := scanner.Err(); err != nil {
		log.Fatalf("Scanner failed %v", err)
	}

}
