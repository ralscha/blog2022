package main

import (
	"context"
	"encoding/binary"
	"encoding/hex"
	"errors"
	"github.com/cockroachdb/pebble"
	"log"
	"math/rand"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"syscall"
	"time"
)

const dbPath = "/var/lib/hibp/pebble"

var db *pebble.DB

func main() {
	var err error
	db, err = pebble.Open(dbPath, &pebble.Options{})
	if err != nil {
		log.Fatalf("Can't open database %v", err)
	}
	defer db.Close()

	mux := http.NewServeMux()
	mux.HandleFunc("/", hipbPasswordFunc)

	server := &http.Server{
		Addr:         ":8080",
		Handler:      mux,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  120 * time.Second,
	}

	go func() {
		quit := make(chan os.Signal, 1)
		signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

		s := <-quit
		log.Printf("Caught signal: %s\n", s.String())
		log.Println("Shutting down")

		ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
		defer cancel()
		err := server.Shutdown(ctx)
		if err != nil {
			log.Fatalf("Shutting down server failed %v", err)
		}
	}()

	err = server.ListenAndServe()
	if !errors.Is(err, http.ErrServerClosed) {
		log.Fatalf("Shutting down server failed %v", err)
	}

}

func hipbPasswordFunc(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not supported", http.StatusInternalServerError)
		return
	}
	if !strings.HasPrefix(r.URL.Path, "/range/") {
		http.Error(w, "invalid query", http.StatusNotFound)
		return
	}

	hashPrefix := strings.TrimPrefix(r.URL.Path, "/range/")
	if !isValidHashPrefix(hashPrefix) {
		http.Error(w, "The hash prefix was not in a valid format", http.StatusBadRequest)
		return
	}

	hashPrefix = strings.ToUpper(hashPrefix)

	lower, err := hex.DecodeString(hashPrefix + "00000000000000000000000000000000000")
	if err != nil {
		http.Error(w, "decoding failed", http.StatusInternalServerError)
		return
	}

	upper, err := hex.DecodeString(hashPrefix + "FFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFF")
	if err != nil {
		http.Error(w, "decoding failed", http.StatusInternalServerError)
		return
	}

	iter := db.NewIter(&pebble.IterOptions{
		LowerBound: lower,
		UpperBound: upper,
	})

	counter := 0

	for iter.First(); iter.Valid(); iter.Next() {
		key := iter.Key()
		value, err := iter.ValueAndErr()
		if err != nil {
			http.Error(w, "getting value failed", http.StatusInternalServerError)
		}
		if writeResult(w, key, value, hashPrefix) {
			return
		}
		counter++
	}

	// upper bound is not included in the iterator (exclusive)
	appearedBytes, closer, err := db.Get(upper)
	if err != nil && err != pebble.ErrNotFound {
		http.Error(w, "getting upper bound value failed", http.StatusInternalServerError)
		return
	}

	if err == nil {
		defer closer.Close()
		writeResult(w, upper, appearedBytes, hashPrefix)
	}

	//  add-padding header
	addPaddingHeader := strings.ToLower(r.Header.Get("add-padding"))
	if addPaddingHeader == "true" {

		minimum := 800 - counter
		random := rand.Intn(200)
		for i := 0; i < minimum+random; i++ {

			rh, err := randomHex(18)
			if err != nil {
				http.Error(w, "randomHex failed", http.StatusInternalServerError)
				return
			}

			_, err = w.Write([]byte(rh[:35] + ":0\n"))
			if err != nil {
				http.Error(w, "writing response failed", http.StatusInternalServerError)
				return
			}
		}
	}

}

func writeResult(w http.ResponseWriter, key, value []byte, hashPrefix string) bool {
	hashStr := strings.TrimPrefix(strings.ToUpper(hex.EncodeToString(key)), hashPrefix)
	appeared, _ := binary.Uvarint(value)
	response := hashStr + ":" + strconv.FormatUint(appeared, 10) + "\n"
	_, err := w.Write([]byte(response))
	if err != nil {
		http.Error(w, "writing response failed", http.StatusInternalServerError)
		return true
	}
	return false
}

func isValidHashPrefix(hashPrefix string) bool {
	if len(hashPrefix) != 5 {
		return false
	}
	for _, c := range hashPrefix {
		if !(c >= '0' && c <= '9') && !(c >= 'A' && c <= 'F') && !(c >= 'a' && c <= 'f') {
			return false
		}
	}
	return true
}

func randomHex(n int) (string, error) {
	bytes := make([]byte, n)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return strings.ToUpper(hex.EncodeToString(bytes)), nil
}
