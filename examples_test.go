package torgo_test

import (
	"fmt"
	"github.com/wybiral/torgo"
	"log"
)

var onion *torgo.Onion
var controller *torgo.Controller

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
