<h1 align="center">Distributed Network Monitoring System (NMS)<img align="center" src="https://github.com/devicons/devicon/blob/master/icons/go/go-original-wordmark.svg" target="_blank" title="GO" alt="GO" width="50" height="50"/></h1>
<h1 align="center">Final Grade: 17/20🌟</h1>

This project implements a **Distributed Network Monitoring System** designed for monitoring and analyzing the state of network links and devices in a distributed environment. The system leverages custom-built application-layer protocols, `NetTask` (UDP-based) and `AlertFlow` (TCP-based), to achieve efficient communication between clients (agents) and the central server.

## Features
- **NetTask Protocol**:

  - UDP-based, handles task distribution, metric collection, and agent registration.
  - Implements reliability mechanisms like ACKs, retransmissions, and packet ID management.

- **AlertFlow Protocol**:
  - TCP-based, ensures reliable transmission of critical alerts when thresholds are exceeded.

- **Network Monitoring Tools**:
  - Uses tools such as `ping` and `iperf` for metrics including latency, bandwidth, jitter, and packet loss.
  - Supports monitoring device metrics like CPU and RAM usage.

- **JSON-Based Task Management**:
  - Tasks and thresholds are defined in JSON configuration files.

- **Scalable and Distributed**:
  - Supports multiple agents reporting to a single server.

## Requirements
- **Languages/Tools**:
  - Go (Golang)
  - CORE Network Emulator (v7.5 or later)
  - Utilities: `ping`, `iperf`

- **Environment**:

  - Linux-based systems for network emulation.
  - JSON files for configuration.

## Project Structure
```
.
├── LICENSE
├── Makefile              # Automates build and run processes
├── README.md             # Project documentation
├── cmd/                  # Entrypoints for server and agent
│   └── nms/
│       ├── agent/
│       │   └── runner.go
│       └── server/
│           └── runner.go
├── configs/              # Configuration files
│   └── tasks.json        # JSON tasks configuration
├── docs/                 # Documentation
│   ├── report.pdf        # Technical report
│   ├── CC Enunciado TP2 24-25.pdf # Project statement
├── internal/             # Core system modules
│   ├── agent/
│   │   ├── alertflow/
│   │   │   └── tcp.go
│   │   └── nettask/
│   │       ├── handlePingTask.go
│   │       └── other-task-handlers.go
│   ├── server/
│   │   ├── alertflow/
│   │   │   └── handleAlert.go
│   │   └── nettask/
│   │       └── task-handlers.go
│   ├── jsonParse/
│   │   ├── parser.go
│   │   └── task.go
│   └── utils/
│       ├── tcp.go
│       └── utils.go
└── topology/
    └── CC-Topologia.imn  # Core emulator topology file

```

## Usage Instructions
### 1. Installation
Clone the repository and navigate to the project directory:
```bash
git clone git@github.com:2101dudu/Network-Monitoring-System.git
cd Network-Monitoring-System
```

### 2. Build and Run
**Server**
1. Build and run the server:
```bash
make server
```
2. Build and start the server in verbose mode for debugging:
```bash
make server-verbose
```

**Agent**
1. Build and run the agent:
```bash
make agent
```
2. Build and start the agent in verbose mode for debugging:
```bash
make agent-verbose
```

### 3. Configuration
Edit the `tasks.json` file in the `configs/` directory to define tasks, thresholds, and monitoring parameters.

## Testing
- Use the **CORE Network Emulator** to simulate networks and validate the system.
- Deploy multiple agents and simulate scenarios with packet loss and alerts.
- Ensure agents register correctly, execute tasks, and report metrics.

## Sample Test Flow
1. Start the server.
1. Deploy agents across different simulated devices.
1. Configure and load tasks using `tasks.json`.
1. Observe metric collection and alert notifications.

## Group Members
- [Edgar Ferreira](https://www.github.com/Edegare)
- [Eduardo Faria](https://www.github.com/2101dudu)
- [Nuno Siva](https://www.github.com/NunoMRS7)
