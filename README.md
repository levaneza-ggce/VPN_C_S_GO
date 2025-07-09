# Go VPN

This is a simple VPN client and server written in Go. It uses a TUN interface to capture network traffic and forwards it over a TCP connection. All traffic is encrypted using AES-256.

## How to Run

### Server

**Windows:**

```bash
cd server
go run server.go --ip <server-tun-ip> --subnet <server-tun-subnet> --port <server-port> --psk <pre-shared-key> [--tap-component-id <component-id>]
```

**Example (Windows):**

```bash
go run server.go --ip 10.0.0.1 --subnet 255.255.255.0 --port 8080 --psk "mysecretkey" --tap-component-id "tap0901"
```

**Linux (CentOS 7 Example):**

1.  **Compile for Linux:**
    ```bash
    cd server
    GOOS=linux GOARCH=amd64 go build -o server
    ```
2.  **Run the server (as root or with sudo):**
    ```bash
    sudo ./server --ip <server-tun-ip> --subnet <server-tun-subnet> --port <server-port> --psk <pre-shared-key>
    ```

    **Example (Linux):**

    ```bash
    sudo ./server --ip 10.0.0.1 --subnet 255.255.255.0 --port 8080 --psk "mysecretkey"
    ```

    **Note:** On Linux, you might need to install `tuntap` utilities if not already present (e.g., `yum install tunctl` or `iproute2`). Also, ensure IP forwarding is enabled (`sysctl -w net.ipv4.ip_forward=1`) and set up NAT/masquerading if you want clients to access the internet through the VPN server (e.g., `iptables -t nat -A POSTROUTING -o eth0 -j MASQUERADE`).

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