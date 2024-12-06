# Distributed Network Monitoring System (NMS)

This project implements a **Distributed Network Monitoring System** designed for monitoring and analyzing the state of network links and devices in a distributed environment. The system leverages custom-built application-layer protocols, `NetTask` (UDP-based) and `AlertFlow` (TCP-based), to achieve efficient communication between clients (agents) and the central server.

## Features
- **NetTask Protocol**:

-- UDP-based, handles task distribution, metric collection, and agent registration.
-- Implements reliability mechanisms like ACKs, retransmissions, and packet ID management.

- **AlertFlow Protocol**:
-- TCP-based, ensures reliable transmission of critical alerts when thresholds are exceeded.

- **Network Monitoring Tools**:
-- Uses tools such as `ping` and `iperf` for metrics including latency, bandwidth, jitter, and packet loss.
-- Supports monitoring device metrics like CPU and RAM usage.

- **JSON-Based Task Management**:
-- Tasks and thresholds are defined in JSON configuration files.

- **Scalable and Distributed**:
-- Supports multiple agents reporting to a single server.

## Requirements
- **Languages/Tools**:
-- Go (Golang)
-- CORE Network Emulator (v7.5 or later)
-- Utilities: `ping`, `iperf`

- **Environment**:

-- Linux-based systems for network emulation.
-- JSON files for configuration.

## Project Structure
```
.
├── LICENSE
├── Makefile
├── README.md
├── go.mod
├── cmd/
│   └── nms/
│       ├── agent/
│       │   └── runner.go
│       └── server/
│           └── runner.go
├── internal/
│   ├── agent/
│   │   ├── tcp.go
│   │   └── udp.go
│   └── server/
│       ├── tcp.go
│       └── udp.go
├── pkg/
│   ├── packet/
│   │   ├── ack.go
│   │   ├── ackMap.go
│   │   └── registration.go
│   └── utils/
│       ├── parser.go
│       ├── udp.go
│       └── utils.go
├── configs/
│   └── tasks.json
├── docs/
│   └── CC Enunciado TP2 24-25.pdf
├── Users/
│   └── eduardofaria/
└── reminders.txt
```


## Requirements
- GO (or your chosen programming language)
- CORE Network Emulator 7.5
- Network utilities: `ping`, `iperf`

## Usage
1. Configure your agents by editing the JSON task file.
2. Run the NMS_Server to start receiving data.
3. Deploy NMS_Agents on the network devices to monitor network health.
4. Alerts will be triggered if critical thresholds are exceeded.

## Installation
1. Clone the repository:
   ```bash
   git clone git@github.com:2101dudu/Network-Monitoring-System.git
   cd Network-Monitoring-System
   ```

2. On one terminal, run:
    ```bash
    make server
    ```

3. On the other terminal, run:
    ```bash
    make agent
    ```

## Group Members
- [Edgar Ferreira](https://www.github.com/Edegare)
- [Eduardo Faria](https://www.github.com/2101dudu)
- [Nuno Siva](https://www.github.com/NunoMRS7)
