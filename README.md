# 📦 Notification Batch

A Go project built using the [Gin Web Framework](https://github.com/gin-gonic/gin).  
This service is designed to automate the daily notification process between **SystemI (AS400)** and a mobile app via API.

## 📌 Overview

Every day, a fixed-length text file is received from **SystemI (AS400)** via FTP.  
The system processes the file, calls a mobile app API for notification delivery, collects the results, and then generates a fixed-length result file.  
This result file is then sent back to **SystemI** via FTP for further use.

## 🚀 How to Run

This project is designed to run on a virtual machine (VM). You can run it directly using:

```bash
go run main.go [env]

For background execution (recommended for production), you may use:
nohup go run main.go [env] > applicationlog.txt 2>&1 &

⚙️ Configuration
All environment-specific settings are placed under the /config directory with the filename format:
{env}.yaml (e.g., dev.yaml, prod.yaml).

You must pass the environment name as a CLI parameter when running the program.

🧩 Dependencies
Make sure to install Go modules before running the project:
go mod tidy

🛠️ Features
Fetch input files via FTP
Parse and process fixed-length text files
Call mobile notification API per record
Collect and log the result of each notification
Generate fixed-length result file
Upload result file to FTP for SystemI

🗄️ Database
The project uses PostgreSQL and text file for data persistence and logging.
Ensure the connection settings are properly defined in your config YAML file.

👨‍💻 Author
SYE Section
Mr. Akkharasarans

📁 Project Structure
.
├── config/             # Environment-specific configuration files (YAML)
│   ├── dev.yaml
│   ├── sit.yaml
│   ├── uat.yaml
│   └── prod.yaml
├── cmd/                # FTP-related logic
│   └── main.go         # Application entry point
├── internal/           # Internal module folder
│   ├── ftp/            # FTP-related logic
│   ├── api/            # API calling logic
│   ├── batch/          # Rule of Business (Batch)
│   ├── config/         # Config Loading
│   ├── Logger/         # Write API and application log
│   ├── routes/         # application routes
│   ├── util/           # utility
│   ├── scheduler/      # Scheduler Batch
│   └── model/          # PostgreSQL interactions and struct
├── log/                # Application Log
├── go.sum              # Go sum
└── go.mod              # Go module definition
