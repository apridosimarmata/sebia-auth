package main

import (
	"fmt"
	"mini-wallet/presentation"
	"net/http"
)

func main() {
	router := presentation.InitServer()

	if err := http.ListenAndServe(":3001", router); err == nil {
		fmt.Println("server listening on port 3001")
	}
}
