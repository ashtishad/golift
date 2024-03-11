## go-lift
GoLift: A Scalable Load Balancing Solution in Go

### Different Load Balancing Algorithms

- **Round Robin**: Distributes requests evenly across all servers, regardless of their current load.
- **Weighted Round Robin**: Allocates more requests to servers with higher capacity, refining the Round Robin approach.
- **Least Connections**: Prefers servers with the fewest active connections, promoting fair load distribution.

### Why Am I Choosing the Least Connection with Round Robin Tiebreaker Algorithm?

The Least Connection strategy, enhanced with a Round Robin tiebreaker, combines efficiency and fairness, especially suitable for high-traffic conditions. It dynamically adapts to server load changes, ensuring optimal resource utilization and user experience without overwhelming any single server.

### Ideal Scenario

This approach excels in distributed environments where demand fluctuates, such as e-commerce sites during sales events. It ensures that incoming requests are evenly distributed, preventing server overload and maintaining smooth operation.

### Algorithmic Explanation

GoLift implements this refined strategy as follows:

1. **Identify Servers**: Determines servers with the lowest active connections.
2. **Single Server Assignment**: Directly assigns requests if one server has the fewest connections.
3. **Round-Robin Tiebreaker**: When multiple servers have the same number of connections, it selects in a Round-Robin manner.

### Expected Behavior

- **Single Server Selection**: For servers `{server3: 0, others: 1+}`, `server3` is chosen.
- **Round-Robin Cycling**: With `{all servers: 1 connection}`, it cycles from `server1` to `server6`, ensuring equitable distribution.
- **No Servers Available**: Returns `nil` if no servers are alive, indicating a need for intervention.
- **Multiple Servers with Least Connections**: Given `{server5: 0, server6: 0, others: 2+}`, selects `server5` or `server6` based on Round-Robin position.
- **Round-Robin Across All Servers**: Continues cycling through all servers `{all servers: 1 connection}`, maintaining fairness.

<p align="right"><a href="#go-lift">↑ Top</a></p>

### Advantages

- **Reduced Server Overload**: Smartly balances load to prevent any server from being overwhelmed.
- **Enhanced Reliability**: More reliable service delivery by evenly distributing requests based on server capacity and current load.

<p align="right"><a href="#go-lift">↑ Top</a></p>

### Limitations

- **Troubleshooting Complexity**: The dynamic nature of the strategy can complicate issue diagnosis.
- **Increased Processing**: The need for constant computation of server loads and decision-making.
- **Capacity Ignorance**: Focuses on connection counts without considering the actual capacity of servers.

<p align="right"><a href="#go-lift">↑ Top</a></p>

### How I Overcame the Limitations

- **Efficient Data Structures**: Implemented optimized data handling to reduce processing overhead.
- **Health Checks**: Integrated server health checks to dynamically adjust the pool based on real-time server status, addressing capacity concerns.
- **Logging and Monitoring**: Enhanced diagnostics with detailed logging and monitoring for better insight and quicker troubleshooting.

<p align="right"><a href="#go-lift">↑ Top</a></p>

### How To Run The App

###### Option 1: Using Makefile

To run the application using the Makefile:

1. Open your terminal and navigate to the project's root directory.
2. (Optional) Adjust the environment variables in the Makefile as necessary to fit your setup.
3. Execute the command: `make run`


###### Option 2: Using Docker

To run the application using Docker:

1. Ensure the Docker Desktop application is running.
2. Open your terminal and navigate to the project's root directory.
3. Build the Docker image for the application by executing: `docker build -t golift:latest .`
4. Start the application using Docker Compose with the command: `docker compose up`

<p align="right"><a href="#go-lift">↑ Top</a></p>

###### Expected Output After Running The App

Regardless of the method chosen to run the application, you should observe output similar to the following in your terminal, indicating that the servers and the load balancer are up and running:


```
server-1         | Server1 listening on port :8000
server-2         | Server2 listening on port :8001
server-3         | Server3 listening on port :8002
server-4         | Server4 listening on port :8003
server-5         | Server5 listening on port :8004
load_balancer    | Load Balancer listening on port 8080

```
<p align="right"><a href="#go-lift">↑ Top</a></p>

#### Project Structure

```plaintext
├── .github
│   └── workflows
│       └── go-ci.yaml             ← GitHub Actions CI workflows (Build, Test, Lint).
├── cmd
│   └── app
│       ├── app.go                 ← Main application logic for consumer servers and load balancer setup.
│       └── handler.go             ← Forwarded http Request with reverse proxy.
├── internal
│   └── domain
│       ├── load_balancer.go       ← Load balancer logic implementation(Least Connection Strategy).
│       ├── load_balancer_test.go  ← Unit Tests for load balancer functionality.
│       ├── server.go              ← Server instance definition and bheaviour.
│       └── server_test.go         ← Unit Tests for server functionality.
│       ├── server_pool.go         ← Server pool for maintaining a list of servers.
│   └── common
│       ├── env_vars.go                ← Utility functions for environment variable management.
│       ├── env_vars_test.go           ← Tests for environment variable utility functions.
│       ├── srvvidgen.go           ← Server ID generation logic(Hash value Server URL and Port).
│       └── srvvidgen_test.go      ← Unit Tests for server ID generation.
│   └── transport
│       └── handler.go             ← Forwarded http Request with reverse proxy.
├── .gitignore                     ← Specifies intentionally untracked files to ignore.
├── .golangci.yaml                 ← Configuration for golangci-lint.
├── compose.yaml                   ← Docker service setup for development environments.
├── Dockerfile                     ← Dockerfile for building the GoLift:latest application image.
├── go.mod                         ← Go module dependencies.
├── main.go                        ← Entry point to start the application services.
├── Makefile                       ← Make commands for building and running the application.
└── readme.md                      ← Project documentation and setup instructions.


```
<p align="right"><a href="#go-lift">↑ Top</a></p>

#### Example Request

```
curl --location '127.0.0.1:8080'

GET -> 127.0.0.1:8080

```
<p align="right"><a href="#go-lift">↑ Top</a></p>
