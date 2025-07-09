# Go VPN

This is a simple VPN client and server written in Go. It uses a TUN interface to capture network traffic and forwards it over a TCP connection. All traffic is encrypted using AES-256.

## How to Run

### Server

```bash
cd server
go run server.go --ip <server-tun-ip> --subnet <server-tun-subnet> --port <server-port> --psk <pre-shared-key> [--tap-component-id <component-id>]
```

**Example:**

```bash
go run server.go --ip 10.0.0.1 --subnet 255.255.255.0 --port 8080 --psk "mysecretkey" --tap-component-id "tap0901"
```

### Client

```bash
cd client
go run client.go --server-ip <server-public-ip> --server-port <server-port> --client-ip <client-tun-ip> --client-subnet <client-tun-subnet> --psk <pre-shared-key> [--tap-component-id <component-id>]
```

**Example:**

```bash
go run client.go --server-ip 192.168.1.100 --server-port 8080 --client-ip 10.0.0.2 --client-subnet 255.255.255.0 --psk "mysecretkey" --tap-component-id "tap0901"
```

**Note:** You will need administrator privileges to run these commands as they configure network interfaces. The `tap-component-id` is optional; if not provided, the `water` library will try to find a default TAP device.
