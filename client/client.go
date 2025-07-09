package main

import (
	"flag"
	"fmt"
	"log"
	"net"
	"os"
	"os/exec"

	"github.com/songgao/water"
)

var (
	serverIP = flag.String("server-ip", "localhost", "IP address of the VPN server")
	serverPort = flag.Int("server-port", 8080, "Port of the VPN server")
	clientIP = flag.String("client-ip", "10.0.0.2", "IP address for the TUN interface")
	clientSubnet = flag.String("client-subnet", "255.255.255.0", "Subnet mask for the TUN interface")
	psk = flag.String("psk", "this-is-a-very-secret-key-123456", "Pre-shared key for encryption")
)

func main() {
	flag.Parse()

	// Connect to the server
	conn, err := net.Dial("tcp", fmt.Sprintf("%s:%d", *serverIP, *serverPort))
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	log.Println("Connected to server")

	// Create a new TUN interface
	ifce, err := water.New(water.Config{
		DeviceType: water.TUN,
	})
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("Interface Name: %s\n", ifce.Name())

	// Configure the TUN interface (Windows specific)
	cmd := exec.Command("netsh", "interface", "ip", "set", "address",
		fmt.Sprintf("name=%s", ifce.Name()), "static", *clientIP, *clientSubnet, "none")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		log.Fatalf("Error configuring TUN interface: %v", err)
	}

	log.Printf("Configured TUN interface %s with IP %s/%s\n", ifce.Name(), *clientIP, *clientSubnet)

	// Goroutine to read from the server and write to the TUN interface
	go func() {
		buf := make([]byte, 2048) // Increased buffer size for encrypted data
		for {
			n, err := conn.Read(buf)
			if err != nil {
				log.Println("Error reading from server:", err)
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
	}()

	// Read from the TUN interface and write to the server
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
			log.Println("Error writing to server:", err)
			return
		}
	}
}
