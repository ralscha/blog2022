package main

import (
	"crypto/sha1"
	"encoding/binary"
	"fmt"
	"log"
	"os"

	badger "github.com/dgraph-io/badger/v3"
	// "github.com/dgraph-io/badger/v3/options"
)

func main() {
	if len(os.Args) != 2 {
		log.Fatalf("Password is missing.\nUsage: %s <password>", os.Args[0])
	}

	path := "/var/lib/hibp/badger"
	db, err := badger.Open(badger.DefaultOptions(path)) // .WithCompression(options.ZSTD))
	if err != nil {
		log.Fatalf("Can't open database %v", err)
	}
	defer db.Close()

	password := os.Args[1]

	h := sha1.New()
	h.Write([]byte(password))
	sha1hash := h.Sum(nil)

	err = db.View(func(txn *badger.Txn) error {
		item, err := txn.Get(sha1hash)
		if err != nil {
			return err
		}

		err = item.Value(func(val []byte) error {
			value, _ := binary.Uvarint(val)
			fmt.Println(value)
			return nil
		})

		return err
	})

	if err != nil {
		log.Fatalf("Error fetching value %v", err)
	}

}
