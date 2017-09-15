package main

import (
	"log"
)

func checkError(ctx string, e error) {
	if e != nil {
		log.Fatalf("Error (%s): %s\n", ctx, e)
	}
}

