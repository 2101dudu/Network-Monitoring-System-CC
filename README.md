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
- **NMS_Agent**: Collects metrics and sends them to the NMS_Server.
- **NMS_Server**: Receives data from agents and generates alerts based on network conditions.

## Requirements
- Python (or your chosen programming language)
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

## Group Members
- [Edgar Ferreira](https://www.github.com/Edegare)
- [Eduardo Faria](https://www.github.com/2101dudu)
- [Nuno Siva](https://www.github.com/NunoMRS7)
