package main

import (
	"fmt"
	"net/http"

	"github.com/wawandco/maildoor/sample/web"
)

// This is a sample application that demonstrates how to use
// the maildoor package.
func main() {
	server, err := web.NewApp()
	if err != nil {
		panic(err)
	}

	fmt.Println("Listening on port 8080")
	if err := http.ListenAndServe(":8080", server); err != nil {
		panic(fmt.Errorf("error starting server: %w", err))
	}

}
