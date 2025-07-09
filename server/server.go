package main

import (
	"flag"
	"fmt"
	"log"
	"net"
	"os"
	"os/exec"
	"runtime"

	"github.com/songgao/water"
)

var (
	ip = flag.String("ip", "10.0.0.1", "IP address for the TUN interface")
	subnet = flag.String("subnet", "255.255.255.0", "Subnet mask for the TUN interface")
	port = flag.Int("port", 8080, "Listening port for the VPN server")
	psk = flag.String("psk", "this-is-a-very-secret-key-123456", "Pre-shared key for encryption")
	tapComponentID = flag.String("tap-component-id", "", "Optional: Component ID of the TAP device (e.g., tap0901)")
)

func main() {
	flag.Parse()

	config := water.Config{
		DeviceType: water.TUN,
	}

	if *tapComponentID != "" {
		config.ComponentID = *tapComponentID
	}

	// Create a new TUN interface
	ifce, err := water.New(config)
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("Interface Name: %s\n", ifce.Name())

	// Configure the TUN interface based on OS
	switch runtime.GOOS {
	case "windows":
		cmd := exec.Command("netsh", "interface", "ip", "set", "address",
			fmt.Sprintf("name=%s", ifce.Name()), "static", *ip, *subnet, "none")
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		if err := cmd.Run(); err != nil {
			log.Fatalf("Error configuring TUN interface on Windows: %v", err)
		}
	case "linux":
		// Set IP address
		cmd := exec.Command("ip", "addr", "add", fmt.Sprintf("%s/%s", *ip, *subnet), "dev", ifce.Name())
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		if err := cmd.Run(); err != nil {
			log.Fatalf("Error setting IP address on Linux: %v", err)
		}

		// Bring up the interface
		cmd = exec.Command("ip", "link", "set", "dev", ifce.Name(), "up")
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		if err := cmd.Run(); err != nil {
			log.Fatalf("Error bringing up interface on Linux: %v", err)
		}
	default:
		log.Fatalf("Unsupported operating system: %s", runtime.GOOS)
	}

	log.Printf("Configured TUN interface %s with IP %s/%s\n", ifce.Name(), *ip, *subnet)

	// Set up the server to listen for incoming connections
	ln, err := net.Listen("tcp", fmt.Sprintf(":%d", *port))
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("Listening for client connection on port %d...\n", *port)

	for {
		conn, err := ln.Accept()
		if err != nil {
			log.Println("Error accepting connection:", err)
			continue
		}
		go handleClient(conn, ifce)
	}
}

func handleClient(conn net.Conn, ifce *water.Interface) {
	defer conn.Close()
	log.Printf("Client connected: %s\n", conn.RemoteAddr())

	// Goroutine to read from the TUN interface and write to the client
	go func() {
		buf := make([]byte, 1500)
		for {
			n, err := ifce.Read(buf)
			if err != nil {
				log.Println("Error reading from TUN interface:", err)
				return
			}

			encrypted, err := Encrypt(buf[:n], []byte(*psk))
			if err != nil {
				log.Println("Error encrypting data:", err)
				return
			}

			_, err = conn.Write(encrypted)
			if err != nil {
				log.Println("Error writing to client:", err)
				return
			}
		}
	}()

	// Read from the client and write to the TUN interface
	buf := make([]byte, 2048) // Increased buffer size for encrypted data
	for {
		n, err := conn.Read(buf)
		if err != nil {
			log.Println("Error reading from client:", err)
				return
			}

			decrypted, err := Decrypt(buf[:n], []byte(*psk))
			if err != nil {
				log.Println("Error decrypting data:", err)
				return
			}

			_, err = ifce.Write(decrypted)
			if err != nil {
				log.Println("Error writing to TUN interface:", err)
				return
			}
		}
	}