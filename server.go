package main

import (
	"king/bootstrap"
	_ "king/routes/server"
	_ "king/rpc/server"
	"king/utils/db"
	"log"
	"os"
)

func main() {

	port := "3000"
	argLen := len(os.Args)

	if argLen > 1 {
		port = os.Args[2]
	}

	log.Println("Running server on port", port)

	db.Connect()
	bootstrap.Start(port, func() {})
}
