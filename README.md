# Load Balancer with Least Connections Algorithm in Golang

## Overview
This project implements a simple Load Balancer in Golang using the **Least Connections** algorithm. The Load Balancer distributes incoming HTTP requests among multiple web servers running concurrently and ensures that requests are forwarded to the server with the least active connections. It also performs **health checks** to ensure that only healthy servers receive requests.

## Features
- **Configuration via JSON file**: The Load Balancer reads server configuration (names, URLs, and health check intervals) from a JSON file.
- **Concurrent Request Handling**: The Load Balancer can handle multiple incoming requests at the same time.
- **Least Connections Algorithm**: Requests are directed to the server with the fewest active connections.
- **Health Checks**: Regular health checks are conducted to ensure servers are healthy and available to handle requests.
- **Server Logs**: Each server logs the number of requests it has processed and provides a health check endpoint.

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
go run servers/server.go <ServerName>

```
Replace <ServerName> with the name of the server (e.g., Server1, Server2, or Server3). Each server will run on a different port (e.g., 8081, 8082, 8083).

### Step 3: Start the Load Balancer
In a separate terminal window, run:
```sh
go run load_balancer.go
```
This will start the Load Balancer, which will listen on the port specified in the config.json file (e.g., 8080).

### Step 4: Simulate Multiple Requests
To test the Load Balancer's behavior with multiple concurrent requests, we have added a function to simulate sending multiple requests. The function will simulate sending `numRequests` requests to the Load Balancer, each of which is processed concurrently.

To simulate the requests, run the following command:
```sh
go run simulate_request.go
```
## Conclusion
This project demonstrates how to implement a **Load Balancer in Golang** using **Least Connections** for intelligent request distribution. It ensures efficient resource utilization and enhances fault tolerance through periodic health checks.
---

