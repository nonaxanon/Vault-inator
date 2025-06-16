# Vault-inator

A local, LastPass-style password manager built with Go, PostgreSQL, and React.

## Features

- Secure password storage using AES-256-GCM encryption
- Local network access only (no internet exposure)
- Simple web UI for managing passwords
- RESTful API for programmatic access

## Prerequisites

- Go 1.21 or later
- PostgreSQL
- Node.js and npm (for React frontend)

## Setup

### 1. Clone the Repository

```bash
git clone https://github.com/nonaxanon/vault-inator.git
cd vault-inator
```

### 2. Set Up PostgreSQL

- Create a PostgreSQL database named `vault_inator`.
- Update the connection string in `cmd/vault-inator/main.go` if needed.

### 3. Build and Run the Go Backend

```bash
go mod download
go build -o vault-inator ./cmd/vault-inator
./vault-inator
```

The server will start on `http://localhost:8080`.

### 4. Set Up the React Frontend

```bash
cd web
npm install
npm run build
```

The frontend will be served by the Go backend at `http://localhost:8080`.

## Usage

- Open your browser and navigate to `http://localhost:8080`.
- Use the web UI to add, view, and delete password entries.

## API Endpoints

- `POST /api/passwords` - Add a new password entry
- `GET /api/passwords` - Retrieve all password entries
- `GET /api/passwords/{id}` - Retrieve a specific password entry
- `DELETE /api/passwords/{id}` - Delete a password entry

## Security Considerations

- The service is designed to run on your local network only.
- Use HTTPS in production (even locally) for added security.
- Regularly backup your database.

## License

MIT
