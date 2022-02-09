package main

import (
	"crypto/sha1"
	"encoding/binary"
	"fmt"
	"log"
	"os"

	"github.com/cockroachdb/pebble"
)

func main() {
	if len(os.Args) != 2 {
		log.Fatalf("Password is missing.\nUsage: %s <password>", os.Args[0])
	}

	dbPath := "/var/lib/hibp/pebble"
	db, err := pebble.Open(dbPath, &pebble.Options{})
	if err != nil {
		log.Fatalf("Can't open database %v", err)
	}
	defer db.Close()
	password := os.Args[1]

	h := sha1.New()
	h.Write([]byte(password))
	sha1hash := h.Sum(nil)

	data, closer, err := db.Get(sha1hash)
	if err != nil && err != pebble.ErrNotFound {
		log.Fatalf("Can't get value for key %v", sha1hash)
	}

	if err == nil {
		defer closer.Close()
		value, _ := binary.Uvarint(data)
		fmt.Printf("Password found. It appears %d times in the database\n", value)
	} else {
		fmt.Printf("Password not found in the HIBP database")
	}

}
