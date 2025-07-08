# Go VPN

This is a simple VPN client and server written in Go. It uses a TUN interface to capture network traffic and forwards it over a TCP connection. All traffic is encrypted using AES-256.

## How to Run

### Server

```bash
cd server
go run server.go
```

### Client

```bash
cd client
go run client.go
```
