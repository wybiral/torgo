# torgo [![GoDoc](https://godoc.org/github.com/wybiral/torgo?status.svg)](https://godoc.org/github.com/wybiral/torgo)
This is a Go library for interacting with Tor over the standard controller interface. It simplifies tasks like creating ephemeral hidden services, working with private keys, and making SOCKS proxied client requests on the Tor network.

# Example Usage

## "Hello world" ephemeral hidden service

This program will generate a random onion address and route a local "Hello world" response server to that hidden service. The service is ephemeral so that no directory or configuration changes are required and it stops serving once the process has been ended.

```go
package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/wybiral/torgo"
)

func main() {
	// Setup handler
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hello world!")
	})
	// Create controller
	controller, err := torgo.NewController("127.0.0.1:9051")
	if err != nil {
		log.Fatal(err)
	}
	// Authenticate to controller using filesystem cookie
	err = controller.AuthenticateCookie()
	if err != nil {
		log.Fatal(err)
	}
	// Configure onion to route hidden service port 80 to localhost:8080
	onion := &torgo.Onion{
		Ports: map[int]string{80: "localhost:8080"},
	}
	// Start the hidden service
	err = controller.AddOnion(onion)
	if err != nil {
		log.Fatal(err)
	}
	// Print newly created onion address
	fmt.Println("Serving at http://" + onion.ServiceID + ".onion")
	// Start serving
	http.ListenAndServe(":8080", nil)
}
```
