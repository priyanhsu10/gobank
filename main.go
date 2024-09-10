package main

import (
	"log"
)

func main() {
	store, err := NewPostgressStore()
	if err != nil {
		log.Fatal(err)
	}
	if err1 := store.init(); err1 != nil {
		log.Fatal(err1)
	}
	server := NewAPIServer("localhost:3000", store)
	server.run()
}
