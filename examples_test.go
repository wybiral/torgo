package torgo_test

import (
	"crypto/rsa"
	"fmt"
	"io"
	"log"
	"os"

	"github.com/wybiral/torgo"
	"golang.org/x/crypto/ed25519"
)

var onion *torgo.Onion
var controller *torgo.Controller
var privateKeyRSA *rsa.PrivateKey
var publicKeyRSA *rsa.PublicKey
var privateKeyEd25519 ed25519.PrivateKey
var publicKeyEd25519 ed25519.PublicKey

// Return a new Controller interface.
func ExampleNewController() {
	// Address of control port
	addr := "127.0.0.1:9051"
	controller, err := torgo.NewController(addr)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(controller)
}

// Authenticate without password or cookie file.
func ExampleController_AuthenticateNone() {
	err := controller.AuthenticateNone()
	if err != nil {
		log.Fatal(err)
	}
}

// Authenticate with password.
func ExampleController_AuthenticatePassword() {
	err := controller.AuthenticatePassword("pa55w0rd")
	if err != nil {
		log.Fatal(err)
	}
}

// Authenticate with cookie file.
func ExampleController_AuthenticateCookie() {
	err := controller.AuthenticateCookie()
	if err != nil {
		log.Fatal(err)
	}
}

// Print external IP address.
func ExampleController_GetAddress() {
	address, err := controller.GetAddress()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(address)
}

// Print total bytes read (downloaded).
func ExampleController_GetBytesRead() {
	n, err := controller.GetBytesRead()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(n)
}

// Print total bytes written (uploaded).
func ExampleController_GetBytesWritten() {
	n, err := controller.GetBytesWritten()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(n)
}

// Print path to Tor configuration file.
func ExampleController_GetConfigFile() {
	config, err := controller.GetConfigFile()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(config)
}

// Print PID of Tor process.
func ExampleController_GetTorPid() {
	pid, err := controller.GetTorPid()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(pid)
}

// Return version of Tor server.
func ExampleController_GetVersion() {
	version, err := controller.GetVersion()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(version)
}

// Delete an onion by its ServiceID
func ExampleController_DeleteOnion() {
	err := controller.DeleteOnion(onion.ServiceID)
	if err != nil {
		log.Fatal(err)
	}
}

// Add onion and generate private key.
func ExampleController_AddOnion() {
	// Define onion that maps virtual port 80 to local port 8080
	onion := &torgo.Onion{Ports: map[int]string{
		80: "127.0.0.1:8080",
	}}
	// Add onion to controller
	err := controller.AddOnion(onion)
	if err != nil {
		log.Fatal(err)
	}
	// Print onion ID (address without ".onion" ending)
	fmt.Println(onion.ServiceID)
}

// Add onion and generate private key (using ED25519-V3 key if supported).
func ExampleController_AddOnion_ed25519() {
	// Define onion that maps virtual port 80 to local port 8080
	onion := &torgo.Onion{
		Ports: map[int]string{
			80: "127.0.0.1:8080",
		},
		PrivateKeyType: "NEW",
		PrivateKey:     "ED25519-V3",
	}
	// Add onion to controller
	err := controller.AddOnion(onion)
	if err != nil {
		log.Fatal(err)
	}
	// Print onion ID (address without ".onion" ending)
	fmt.Println(onion.ServiceID)
}

// Send signal to reload configuration.
func ExampleController_Signal_reload() {
	err := controller.Signal("RELOAD")
	if err != nil {
		log.Fatal(err)
	}
}

// Send NEWSYM signal to switch to new clean circuits.
func ExampleController_Signal_newnym() {
	err := controller.Signal("NEWNYM")
	if err != nil {
		log.Fatal(err)
	}
}

// Create a new Tor SOCKS HTTP client and request current IP from httpbin.org.
func ExampleNewClient_httpget() {
	// Create client from SOCKS proxy address
	client, err := torgo.NewClient("127.0.0.1:9050")
	if err != nil {
		log.Fatal(err)
	}
	// Perform HTTP GET request
	resp, err := client.Get("http://httpbin.org/ip")
	if err != nil {
		log.Fatal(err)
	}
	// Copy response to Stdout
	io.Copy(os.Stdout, resp.Body)
}

// Create an Onion and start hidden service from an ed25519.PrivateKey.
func ExampleOnionFromEd25519() {
	// Create Onion from private key (does not start hidden service)
	onion, err := torgo.OnionFromEd25519(privateKeyEd25519)
	if err != nil {
		log.Fatal(err)
	}
	// Set port mapping from hidden service 80 to localhost:8080
	onion.Ports[80] = "localhost:8080"
	// Print service ID for Onion
	fmt.Println(onion.ServiceID)
	// Start hidden service
	controller.AddOnion(onion)
}

// Calculate Tor service ID from ed25519.PublicKey.
func ExampleServiceIDFromEd25519() {
	// Calculate service ID
	serviceID, err := torgo.ServiceIDFromEd25519(publicKeyEd25519)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(serviceID)
}

// Create an Onion and start hidden service from an rsa.PrivateKey.
func ExampleOnionFromRSA() {
	// Create Onion from private key (does not start hidden service)
	onion, err := torgo.OnionFromRSA(privateKeyRSA)
	if err != nil {
		log.Fatal(err)
	}
	// Set port mapping from hidden service 80 to localhost:8080
	onion.Ports[80] = "localhost:8080"
	// Print service ID for Onion
	fmt.Println(onion.ServiceID)
	// Start hidden service
	controller.AddOnion(onion)
}

// Calculate Tor service ID from *rsa.PublicKey.
func ExampleServiceIDFromRSA() {
	// Calculate service ID
	serviceID, err := torgo.ServiceIDFromRSA(publicKeyRSA)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(serviceID)
}
