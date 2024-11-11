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
├── README.md
├── go.mod
├── cmd/
│   └── nms/
│       └── main.go
├── configs/
│   └── settings.json
├── docs/
│   └── CC Enunciado TP2 24-25.pdf
├── internal/
│   ├── agent/
│   │   ├── agent_config/
│   │   │   ├── tcp_config.go
│   │   │   └── udp_config.go
│   │   └── agent_runner/
│   │       └── runner.go
│   └── server/
│       ├── server_config/
│       │   ├── tcp_config.go
│       │   └── udp_config.go
│       └── server_runner/
│           └── runner.go
├── pkg/
│   ├── packet/
│   │   ├── ack.go
│   │   └── registration.go
│   └── utils/
│       ├── ack.go
│       ├── parser.go
│       ├── udp.go
│       └── utils.go
├── Users/
│   └── eduardofaria
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
   git clone https://github.com/your-username/distributed-network-monitoring.git
   cd distributed-network-monitoring
   ```

2. On one terminal, run:
    ```bash
    go run internal/server/server_runner/runner.go
    ```

3. On the other terminal, run:
    ```bash
    go run internal/agent/agent_runner/runner.go
    ```

## Group Members
- [Edgar Ferreira](https://www.github.com/Edegare)
- [Eduardo Faria](https://www.github.com/2101dudu)
- [Nuno Siva](https://www.github.com/NunoMRS7)
