// Create a "Hello world!" ephemeral hidden service.
//
// Ephemeral hidden services are useful for when you want to serve something
// from a temporary address without having to change your torrc configuration or
// store any keys on your filesystem. Once the process has been stopped with
// Ctrl-C the hidden service will stop being available and the private key will
// also no longer be available.
//
// This means that you won't be able to recreate the same onion address twice,
// but you can also run ephemeral services by supplying your own keys (see the
// OnionFromRSA and OnionFromEd25519 methods) stored, for instance in an
// encrypted database.

package main

import (
	"fmt"
	"log"
	"net"
	"net/http"

	"github.com/wybiral/torgo"
)

// PrivateKeyType is the type of private key to use for onion.
// Options:
//     "BEST" (recommended)
//     "RSA1024" (short, old type)
//     "ED25519-V3" (long, new type)
const PrivateKeyType = "BEST"

// ControllerAddr is the Tor controller interface address
// Note: Tor Browser Bundle uses 9151 instead of 9051 like the daemon
const ControllerAddr = "127.0.0.1:9051"

func main() {
	// Setup handler
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hello world!")
	})
	// Create local listener on next available port
	listener, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		log.Fatal(err)
	}
	// Get listener port
	port := listener.Addr().(*net.TCPAddr).Port
	addr := fmt.Sprintf("127.0.0.1:%d", port)
	// Create controller
	controller, err := torgo.NewController(ControllerAddr)
	if err != nil {
		log.Fatal(err)
	}
	// Authenticate to controller using filesystem cookie
	// You may need to change this depending on your torrc configuration
	err = controller.AuthenticateCookie()
	if err != nil {
		log.Fatal(err)
	}
	// Configure onion to route hidden service port 80 to server address
	onion := &torgo.Onion{
		Ports:          map[int]string{80: addr},
		PrivateKeyType: PrivateKeyType,
	}
	// Start the hidden service
	err = controller.AddOnion(onion)
	if err != nil {
		log.Fatal(err)
	}
	// Print newly created onion address
	fmt.Println("Local port is", port)
	fmt.Println("Serving at http://" + onion.ServiceID + ".onion")
	// Start serving
	http.Serve(listener, nil)
}
