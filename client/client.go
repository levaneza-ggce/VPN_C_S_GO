package main

import (
	"log"
	"net"

	"github.com/songgao/water"
)

// A pre-shared key for encryption. In a real-world application, you would want to
// use a more secure key exchange mechanism.
var psk = []byte("this-is-a-very-secret-key-123456")

func main() {
	// Connect to the server
	conn, err := net.Dial("tcp", "localhost:8080")
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

	log.Printf("Interface Name: %s
", ifce.Name())

	// Goroutine to read from the server and write to the TUN interface
	go func() {
		buf := make([]byte, 2048) // Increased buffer size for encrypted data
		for {
			n, err := conn.Read(buf)
			if err != nil {
				log.Println("Error reading from server:", err)
				return
			}

			decrypted, err := Decrypt(buf[:n], psk)
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

		encrypted, err := Encrypt(buf[:n], psk)
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