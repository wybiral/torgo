package torgo_test

import (
	"log"
	"fmt"
	"github.com/wybiral/torgo"
)

var onion *torgo.Onion
var controller *torgo.Controller

func ExampleNewController() {
	// Address of control port
	addr := "127.0.0.1:9051"
	controller, err := torgo.NewController(addr)
	if err != nil {
		log.Fatal(err)
	}
	// Authenticate with cookie file
	err = controller.AuthenticateCookie()
	if err != nil {
		log.Fatal(err)
	}
}

func ExampleController_DeleteOnion() {
	// Delete an onion by its ServiceID
	err := controller.DeleteOnion(onion.ServiceID)
	if err != nil {
		log.Fatal(err)
	}
}

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