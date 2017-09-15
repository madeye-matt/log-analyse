package main

import (
	"os"
	"log"
)

func initLogging(filename string) *os.File {
	f, err := os.OpenFile(filename, os.O_CREATE|os.O_RDWR, 0666)
	checkError("initLogging", err)

	// assign it to the standard logger
	log.SetOutput(f)

	return f
}
