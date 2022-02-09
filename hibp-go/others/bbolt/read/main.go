package main

import (
	"crypto/sha1"
	"encoding/binary"
	"fmt"
	"log"
	"os"

	bolt "go.etcd.io/bbolt"
)

func main() {
	if len(os.Args) != 2 {
		log.Fatalf("Password is missing.\nUsage: %s <password>", os.Args[0])
	}

	dbPath := "/var/lib/hibp/bbolt"
	db, err := bolt.Open(dbPath, 0666, &bolt.Options{ReadOnly: true})
	if err != nil {
		log.Fatalf("Can't open database %v", err)
	}
	defer db.Close()
	password := os.Args[1]

	h := sha1.New()
	h.Write([]byte(password))
	sha1hash := h.Sum(nil)

	err = db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("hibp"))
		v := b.Get(sha1hash)
		if v != nil {
			value, _ := binary.Uvarint(v)
			fmt.Printf("Password found. It appears %d times in the database\n", value)
		} else {
			fmt.Printf("Password not found in the HIBP database")
		}
		return nil
	})

	if err != nil {
		log.Fatalf("Error fetching value %v", err)
	}

}
