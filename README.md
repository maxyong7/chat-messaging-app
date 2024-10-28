# Chat Messaging App

A clean and scalable chat messaging app built with **Golang**, designed to maintain modularity, testability, and independence in business logic. This project leverages **PostgreSQL** and the **Gin framework** to create a lightweight and efficient backend solution.

## Features
- **Modular architecture**: Clear separation of business logic and infrastructure.
- **RESTful APIs**: Powered by the Gin framework.
- **PostgreSQL integration**: Ensures reliable data storage.
- **Dependency injection**: Promotes flexibility and scalability.
- **Swagger documentation**: Simplifies API exploration.
- **Testing support**: Unit and integration tests included.

## Installation

1. Clone the repository:
   ```bash
   git clone https://github.com/maxyong7/chat-messaging-app.git
   cd chat-messaging-app
   ```
2. Set up PostgreSQL and configure the .env file with database credentials.
3. Install project dependencies:
    ```bash
    go mod download
    ```
4. Run the application:
    ```bash
    go run main.go
    ```

## Folder Structure
- **/api:** API endpoints and routing.
- **/models:** Database models.
- **/services:** Core business logic.
- **/config:** Application configurations.
