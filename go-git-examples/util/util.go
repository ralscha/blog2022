package util

import "log"

func CheckError(err error) {
	if err != nil {
		log.Fatalf("error occurred, %v", err)
	}
}
