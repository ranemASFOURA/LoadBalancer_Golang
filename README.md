# Load Balancer with Least Connections Algorithm in Golang

## Overview
This project implements a simple Load Balancer in Golang using the **Least Connections** algorithm. The Load Balancer distributes incoming HTTP requests among multiple web servers running concurrently and ensures that requests are forwarded to the server with the least active connections. It also performs **health checks** to ensure that only healthy servers receive requests.

## Features
- **Configuration via User Input**: Users can manually specify ports for the web servers and the Load Balancer.
- **Concurrent Request Handling**: The Load Balancer can handle multiple incoming requests simultaneously.
- **Least Connections Algorithm**: Requests are distributed to the server with the fewest active connections.
- **Health Checks**: Periodic checks ensure that only healthy servers receive traffic.
- **Server Logs**: Each server logs received requests and the total requests handled.

## Installation & Setup
### Prerequisites
- Golang installed (`go version` to check)

### Step 1: Clone the Repository
```sh
git clone https://github.com/ranemASFOURA/LoadBalancer_Golang.git
cd load-balancer
```

### Step 2: Start Multiple Web Servers
To simulate multiple servers, open multiple terminal windows and run the following command in each:
```sh
go run servers/server.go
```
Enter the port number for each server when prompted (e.g., `8081`, `8082`, `8083`).

### Step 3: Start the Load Balancer
In a separate terminal window, run:
```sh
go run load_balancer.go
```
Enter the port for the Load Balancer (e.g., `8080`).

### Step 4: Simulate Multiple Requests
To test Load Balancer behavior, run the following command to send multiple concurrent requests:
```sh
for /L %i in (1,1,8) do start /b curl -s http://localhost:8080
```
This command sends 8 requests to the Load Balancer.


## Expected Output
### Load Balancer Logs:
```
Load Balancer is running on port 8080
Load Balancer: Redirecting request to Server1 (Active connections: 1)
Load Balancer: Redirecting request to Server2 (Active connections: 1)
Load Balancer: Redirecting request to Server3 (Active connections: 1)
...
```
### Web Server Logs:
```
2025/02/28 19:54:27 Server Server1 received a request. Total requests handled: 1
2025/02/28 19:54:27 Server Server2 received a request. Total requests handled: 1
2025/02/28 19:54:27 Server Server3 received a request. Total requests handled: 1
...
```


## Conclusion
This project demonstrates how to implement a **Load Balancer in Golang** using **Least Connections** for intelligent request distribution. It ensures efficient resource utilization and enhances fault tolerance through periodic health checks.
---

