# Virtual Server 

## Table of Contents

- [Description](#description)
- [Features](#features)
- [Prerequisites](#prerequisites)
- [Setup & Installation](#setup--installation)
- [API Endpoints](#api-endpoints)
- [Project Structure](#project-structure)
- [Testing](#testing)

## Description

This project provides a backend API to simulate the management of virtual servers. It allows for the creation, retrieval, and state-transition of server instances. A core feature is the implementation of a Finite State Machine (FSM) to ensure valid transitions between server states (e.g., running, stopped, terminated). All significant lifecycle events are logged and can be retrieved.

## Features

- Server Creation: Provision new virtual server instances.
- Server Retrieval: Fetch details of a single server or a list of all servers.
-  Server Actions (FSM): Perform actions like start, stop, reboot, and terminate on servers.Enforces valid state transitions (e.g., cannot start a terminated server).Returns HTTP 409 Conflict for invalid state transitions.
- Lifecycle Logging: Records all significant server events (creation, status changes, actions, denials) as structured logs.

## Technologies Used

* Go (Golang): The primary programming language.
* Gin Web Framework: For building the RESTful API endpoints.
* GORM: An ORM (Object-Relational Mapper) for database interactions.
* PostgreSQL: The relational database used for persistence.
* Docker & Docker Compose: For containerizing the PostgreSQL database, ensuring easy setup and environment consistency.

## Prerequisites
Before you begin, ensure you have the following installed on your system:

- Go: Version 1.18 or higher 
- Docker & Docker Compose

## Setup & Installation
Follow these steps to get the project up and running locally:

### 1. Clone the Repository
```bash
git clone https://github.com/gitshubham45/virtualServer.git 
cd virtualServer
```

### 2. Configure Environment Variables
Create a .env file in the root of your project directory and add your PostgreSQL database credentials and application port.

```bash
# .env
DB_USER=your_postgres_user       
DB_PASSWORD=your_postgres_password 
DB_HOST=localhost
DB_NAME=your_database_name       
DB_PORT=5432
DB_TIME_ZONE=Asia/Kolkata        
APP_PORT=8080                    
```
Note: Ensure DB_HOST is localhost if running Docker locally.

### 3. Start PostgreSQL with Docker Compose
Navigate to the project root directory where your docker-compose.yml is located and start the database container:

```bash
docker-compose up -d
```
This will create and start a PostgreSQL container named my_postgres.


### 4. Run the Go Application

```bash
go mod tidy 
go run cmd/main.go
```
Your API server should now be running, typically on http://localhost:8080.

## API Endpoints
Here's a list of the API endpoints available, along with example curl commands. Replace http://localhost:8080 with your actual server address if different.

### 1. Create a New Server
Creates a new virtual server instance. Initial status is running (or whatever your creation logic sets). A SERVER_CREATED log event is recorded.

- Method: POST
- Path: /servers
- Body:
    ```bash
    {
        "type": "basic",    // "basic", "plus", or "prime"
        "region": "India"   // e.g., "India", "US East", etc.
    }
    ```
- Example curl:
    ```bash
    curl -X POST http://localhost:8080/servers \
        -H "Content-Type: application/json" \
        -d '{"type": "basic", "region": "US East"}'
    ```
- Success Response (201 Created):
    ```bash
    {
        "message": "Server created successfully",
        "server": {
            "id": "a1b2c3d4-e5f6-7890-1234-567890abcdef",
            "serverNumber": 1,
            "upTime" : 5,
            "billingRate": 5,
            "status": "running",
            "region": "US East",
            "type": "basic",
            "createdAt": "2025-07-28T10:00:00Z",
            "updatedAt": "2025-07-28T10:00:00Z",
            "deletedAt": null
        },
    }
    ```
### 2. Get Server Details
Retrieves the details of a specific server by its UUID.

- Method: GET
- Path: /servers/:id (e.g., /servers/a1b2c3d4-e5f6-7890-1234-567890abcdef)
- Example curl:
    ```bash
    curl -X GET http://localhost:8080/servers/a1b2c3d4-e5f6-7890-1234-567890abcdef
    ```
- Success Response (200 OK): Same as the server object in the create response.
- Error Response (404 Not Found): If server ID does not exist.

### 3. Perform Server Action
Initiates a state-changing action on a specific server, enforcing FSM transitions. Logs are recorded for actions and denials.

- Method: POST
- Path: /servers/:id/action
- Body:
    ```bash
    {
        "action": "start" // Can be "start", "stop", "reboot", or "terminate"
    }
    ```
- Example curl (Stop a running server):
    ```bash
    curl -X POST http://localhost:8080/servers/a1b2c3d4-e5f6-7890-1234-567890abcdef/action \
        -H "Content-Type: application/json" \
        -d '{"action": "stop"}'
    ```
- Success Response (200 OK):

    ```bsh
    {
        "message": "Server action completed successfully",
        "server": {
            // Updated server object with new status
        }
    }
    ```
- Error Response (409 Conflict): If the action is invalid for the current server state.

    ```bash
    {
        "message": "Server is already running."
    }
    ```
- Error Response (404 Not Found): If server ID does not exist.
- Error Response (400 Bad Request): If action is missing or unknown.

### 4. Get Server Lifecycle Logs
Retrieves the last 100 lifecycle events for a specific server, ordered by most recent first.

- Method: GET
- Path: /servers/:id/logs (e.g., /servers/a1b2c3d4-e5f6-7890-1234-567890abcdef/logs)
- Example curl:

    ```bash
    curl -X GET "http://localhost:8080/servers/a1b2c3d4-e5f6-7890-1234-567890abcdef/logs"
    ``` 
- Success Response (200 OK):
    ```bash
    {
        "message": "Server logs fetched successfully",
        "logs": [
            {
                "ID": 10,
                "createdAt": "2025-07-28T10:05:00Z",
                "updatedAt": "2025-07-28T10:05:00Z",
                "deletedAt": null,
                "serverId": "a1b2c3d4-e5f6-7890-1234-567890abcdef",
                "eventType": "STATUS_CHANGE",
                "message": "Server status changed from 'running' to 'stopped'",
                "oldStatus": "running",
                "newStatus": "stopped",
                "timestamp": "2025-07-28T10:05:00Z"
            },
            {
                "ID": 9,
                "createdAt": "2025-07-28T10:01:00Z",
                "updatedAt": "2025-07-28T10:01:00Z",
                "deletedAt": null,
                "serverId": "a1b2c3d4-e5f6-7890-1234-567890abcdef",
                "eventType": "SERVER_CREATED",
                "message": "New server provisioned.",
                "newStatus": "running",
                "timestamp": "2025-07-28T10:01:00Z"
            }
            // ... up to 100 log entries
        ]
    }
    ```
- Error Response (500 Internal Server Error): If there's a database error.
