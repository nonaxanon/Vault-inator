# Vault-inator 🔐

A secure, local-first password manager built with Go and React. Vault-inator keeps your passwords encrypted and stored locally, giving you complete control over your data while providing a modern, user-friendly interface for password management.

## Features ✨

- 🔒 Military-grade AES-256-GCM encryption
- 🎨 Modern, responsive dark theme UI
- 🔍 Search and sort functionality
- 📋 One-click copy for usernames, passwords, and URLs
- 👁️ Password visibility toggle
- 🔄 Master password management
- 📱 Mobile-friendly design

## Security 🔐

Vault-inator implements several security measures to protect your passwords:

- AES-256-GCM encryption for all stored passwords
- SHA-256 key derivation from master password
- Bcrypt hashing for master password storage
- Unique encryption nonce for each password
- Secure database storage with PostgreSQL
- SSL support for database connections

## Prerequisites 📋

- Docker and Docker Compose
- Git

## Installation 🚀

1. Clone the repository:
   ```bash
   git clone https://github.com/yourusername/vault-inator.git
   cd vault-inator
   ```

2. Start the application using Docker Compose:
   ```bash
   docker compose up --build
   ```

The application will be available at:
- Frontend: http://localhost:3000
- Backend API: http://localhost:8080
- PostgreSQL: localhost:5432

## Usage 📖

1. Open your browser and navigate to `http://localhost:3000`
2. Set up your master password when first launching the application
3. Start adding your passwords with the "Add New Password" button
4. Use the search and sort features to organize your passwords
5. Click the copy button to copy usernames, passwords, or URLs
6. Use the eye icon to toggle password visibility

## Development 🛠️

### Project Structure

```
vault-inator/
├── cmd/
│   └── api/
│       └── main.go
├── internal/
│   ├── api/
│   ├── config/
│   ├── encryption/
│   ├── services/
│   └── storage/
├── web/
│   ├── src/
│   │   ├── App.js
│   │   └── App.css
│   ├── nginx.conf
│   └── package.json
├── docker-compose.yml
├── go.mod
└── README.md
```

### Running Tests

```bash
# Backend tests
go test ./...

# Frontend tests
cd web
npm test
```

### Development with Docker

For development, you can use the following commands:

```bash
# Start all services
docker compose up

# Start services in detached mode
docker compose up -d

# View logs
docker compose logs -f

# Stop all services
docker compose down

# Rebuild and start services
docker compose up --build
```

## Contributing 🤝

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add some amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## Security Considerations 🔒

- Always use a strong master password
- Keep your master password secure and don't share it
- Regularly update your master password
- Ensure your database is properly secured
- Keep the application and its dependencies updated

## License 📄

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## Acknowledgments 🙏

- [Go](https://golang.org/)
- [React](https://reactjs.org/)
- [PostgreSQL](https://www.postgresql.org/)
- [AES-GCM](https://en.wikipedia.org/wiki/Galois/Counter_Mode)
- [Docker](https://www.docker.com/)

## Support 💬

If you encounter any issues or have questions, please open an issue in the GitHub repository.
