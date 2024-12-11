# Project: User-Service, Task-Service, and Notification-Service

This project consists of three microservices: **User-Service**, **Task-Service**, and **Notification-Service**. Each service is responsible for specific functionality, and they communicate with each other to achieve the overall system goals.

### Key Components

- **User-Service**: Manages user signup, login, and authentication using JWT.
- **Task-Service**: Manages task creation, approval, and updates. It communicates with Notification-Service to send notifications.
- **Notification-Service**: Sends email notifications to users for task creation, approval, and completion. It interacts with an SMTP server to send emails.

---

## Prerequisites

Ensure you have the following tools installed:

- **GoLang** (version 1.19 or higher)
- **PostgreSQL** (for database management)
- **Docker** (optional, for containerized setup)

Set up your environment variables for the following:

- JWT secret key
- SMTP credentials (for email sending)
- Database configuration (PostgreSQL connection details)

---

## Docker Compose

To build and start the services, use the following command:

```bash
docker-compose up --build
```


## Set Up Database

You need to create PostgreSQL databases for each service:

- **service-db**: Stores user data, task and approval data.

### Database Configuration

Before proceeding, ensure that PostgreSQL is properly installed and configured. Set the following environment variables for database access:

- **DB_HOST**: The host where the PostgreSQL server is running (e.g., `localhost` if using local setup).
- **DB_PORT**: The port where PostgreSQL is accessible (default is `5432`).
- **DB_USER**: PostgreSQL username (e.g., `postgres`).
- **DB_PASSWORD**: PostgreSQL password for the specified user.
- **DB_NAME**: The name of the database to be used (e.g., `user-service-db`, `task-service-db`).

---

## Database Schema

You can use the following SQL commands to create the necessary tables in PostgreSQL:

### 1. **User-Service Database Schema**

```sql
-- Table for storing user data
CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    email VARCHAR(255) UNIQUE NOT NULL,
    password VARCHAR(255) NOT NULL
);
```

```sql
-- Table for storing task data
CREATE TABLE tasks (
    id SERIAL PRIMARY KEY,
    title VARCHAR(255) NOT NULL,
    creator_id INT NOT NULL,
    status VARCHAR(50) DEFAULT 'pending',
    created_at TIMESTAMP DEFAULT NOW(),
    FOREIGN KEY (creator_id) REFERENCES users(id)
);
```

```sql
-- Table for storing task approval data
CREATE TABLE approvals (
    id SERIAL PRIMARY KEY,
    task_id INT NOT NULL,
    approver_id INT NOT NULL,
    comment TEXT,
    approved BOOLEAN DEFAULT FALSE,
    FOREIGN KEY (task_id) REFERENCES tasks(id),
    FOREIGN KEY (approver_id) REFERENCES users(id)
);
```

## Database Migration

You can migrate the schema using either **Docker Compose** or **psql** commands. Follow the appropriate instructions based on your setup.

### 1. Using Docker Compose

If you are using Docker Compose to set up your environment, the migration process will be handled as part of the service startup.

#### Steps:

1. Ensure that `docker-compose.yml` file has the correct database configuration.
2. Run the following command to build and start the services:

```bash
docker-compose up --build
```


## cURL Commands

### 1. User-Service

#### Signup: Register a user to get a user ID.

```bash
curl -X POST http://localhost:8000/signup \
-H "Content-Type: application/json" \
-d '{"name": "Chandan Ahire", "email": "chandanahire25@gmail.com", "password": "password123"}'
```
#### Login: Authenticate the user to receive a JWT token.

```bash
curl -X POST http://localhost:8000/login \
-H "Content-Type: application/json" \
-d '{"email": "chandanahire25@gmail.com", "password": "password123"}'
```
#### Response Example:
```json
{
  "token": "your-jwt-token"
}
```

### 2. Task-Service
#### Create Task: Use the JWT token from login to create a task.
```bash
curl -X POST http://localhost:8001/create-task \
-H "Content-Type: application/json" \
-H "Authorization: Bearer your-jwt-token" \
-d '{
  "title": "Approve Task #123",
  "approvers": [1, 2],
  "description": "Task description here"
}'
```

#### Approve Task: Approvers use their JWT token to approve a task.
```bash
curl -X POST http://localhost:8001/tasks/approve \
-H "Content-Type: application/json" \
-H "Authorization: Bearer approver-jwt-token" \
-d '{"task_id": 123, "status": "approved", "comment": "Task is good to go!"}'
```

#### View Task Status: Fetch task details to check its status.
```bash
curl -X GET http://localhost:8001/tasks/123 \
-H "Authorization: Bearer your-jwt-token"
```

#### 3. Notification-Service
The Notification-Service is automatically triggered by the Task-Service when a task is created, approved, or completed. It sends email notifications based on task events. The Notification-Service does not have direct API endpoints, but you can monitor logs or check the SMTP server for email notifications sent out.

