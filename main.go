package main

import (
	"fmt"
	"mini-wallet/presentation"
	"net/http"
)

func main() {
	router, port := presentation.InitServer()

	if err := http.ListenAndServe(fmt.Sprintf(":%s", port), router); err == nil {
		fmt.Println("server listening on port " + port)
	}
}
