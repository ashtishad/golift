## GoLift: A Scalable Load Balancing Solution in Go

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

### Advantages

- **Reduced Server Overload**: Smartly balances load to prevent any server from being overwhelmed.
- **Enhanced Reliability**: More reliable service delivery by evenly distributing requests based on server capacity and current load.

### Limitations

- **Troubleshooting Complexity**: The dynamic nature of the strategy can complicate issue diagnosis.
- **Increased Processing**: The need for constant computation of server loads and decision-making.
- **Capacity Ignorance**: Focuses on connection counts without considering the actual capacity of servers.

### How I Overcame the Limitations

- **Efficient Data Structures**: Implemented optimized data handling to reduce processing overhead.
- **Health Checks**: Integrated server health checks to dynamically adjust the pool based on real-time server status, addressing capacity concerns.
- **Logging and Monitoring**: Enhanced diagnostics with detailed logging and monitoring for better insight and quicker troubleshooting.
