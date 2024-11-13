# Distributed Network Monitoring System

This project implements a Distributed Network Monitoring System (NMS) capable of gathering and analyzing network metrics from various devices. It uses a client-server architecture, where:

- **NMS_Agent**: Deployed on network devices, collects metrics periodically, and communicates with the server.
- **NMS_Server**: Centralized server that receives and processes data from the agents and provides alerts on critical network conditions.

## Features
- **UDP Protocol (NetTask)**: For task distribution and continuous metric collection.
- **TCP Protocol (AlertFlow)**: For critical alerts and notifications.
- **Network Monitoring Tools**: Uses tools like `ping` and `iperf` for gathering metrics such as latency, bandwidth, jitter, packet loss, etc.
- **JSON Task Configuration**: Tasks are assigned to agents via a JSON configuration file.

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
