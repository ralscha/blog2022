package main

import (
	"crypto/sha1"
	"encoding/binary"
	"fmt"
	"log"
	"os"

	"github.com/syndtr/goleveldb/leveldb"
)

func main() {
	if len(os.Args) != 2 {
		log.Fatalf("Password is missing.\nUsage: %s <password>", os.Args[0])
	}

	dbPath := "/var/lib/hibp/goleveldb"
	db, err := leveldb.OpenFile(dbPath, nil)
	if err != nil {
		log.Fatalf("Can't open database %v", err)
	}
	defer db.Close()
	password := os.Args[1]

	h := sha1.New()
	h.Write([]byte(password))
	sha1 := h.Sum(nil)

	data, err := db.Get(sha1, nil)
	if err != nil && err != leveldb.ErrNotFound {
		log.Fatalf("Can't get value for key %v", sha1)
	}

	if err == nil {
		value, _ := binary.Uvarint(data)
		fmt.Printf("Password found. It appears %d times in the database\n", value)
	} else {
		fmt.Printf("Password not found in the HIBP database")
	}

}
