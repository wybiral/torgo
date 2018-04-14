package main

import (
	"flag"
	"log"
	"fmt"
	"github.com/wybiral/torgo"
	"strconv"
	"strings"
	"io/ioutil"
	"os"
)

func main() {
	// Setup flags
	var controlPort int
	flag.IntVar(&controlPort, "c", 9051, "Tor control port")
	var keyFile string
	flag.StringVar(&keyFile, "k", "", "Private key file")
	flag.Usage = func() {
		fmt.Println("NAME:")
		fmt.Println("  torgo\n")
		fmt.Println("USAGE:")
		fmt.Println("  torgo [options] port_mapping\n")
		fmt.Println("OPTIONS:")
		flag.PrintDefaults()
	}
	flag.Parse()

	// Port mapping parsed from positional args
	args := flag.Args()
	if len(args) == 0 {
		flag.Usage()
		return
	}
	ports, err := parsePorts(args)
	if err != nil {
		log.Fatal(err)
	}

	// Create instance of controller
	addr := fmt.Sprintf("127.0.0.1:%d", controlPort)
	c, err := torgo.NewController(addr)
	if err != nil {
		log.Fatal(err)
	}

	// Perform basic cookie auth
	// TODO: Add other auth methods
	err = c.AuthenticateCookie()
	if err != nil {
		log.Fatal(err)
	}

	onion := &torgo.Onion{Ports: ports}

	// If a key file is specified but doesn't exist it will be created
	createKeyFile := false
	if len(keyFile) > 0 {
		key, err := ioutil.ReadFile(keyFile)
		if err != nil {
			if os.IsNotExist(err) {
				createKeyFile = true
			} else {
				log.Fatal(err)
			}
		} else {
			split := strings.SplitN(string(key), ":", 2)
			onion.PrivateKeyType = split[0]
			onion.PrivateKey = split[1]
		}
	}

	// Add onion to controller
	err = c.Add(onion)
	if err != nil {
		log.Fatal(err)
	}

	// Create new key file with resulting private key
	if createKeyFile {
		key := []byte(onion.PrivateKeyType + ":" + onion.PrivateKey)
		err = ioutil.WriteFile(keyFile, key, 0644)
		if err != nil {
			log.Fatal(err)
		}
	}

	fmt.Println(onion.ServiceID + ".onion")
	// Block forever
	select {}
}

// Parse array of port mappings in the form of virtualPort,localAddress
// If a virtual port isn't specified it defaults to 80
func parsePorts(args []string) (map[int]string, error) {
	ports := make(map[int]string)
	for _, raw := range args {
		if strings.Contains(raw, ",") {
			split := strings.SplitN(raw, ",", 2)
			port, err := strconv.Atoi(split[0])
			if err != nil {
				return nil, err
			}
			ports[port] = split[1]
		} else {
			ports[80] = raw
		}
	}
	return ports, nil
}