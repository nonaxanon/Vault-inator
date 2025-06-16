# Vault-inator 🔐

A secure, local-first password manager built with Go and React. Vault-inator keeps your passwords encrypted and stored locally, giving you complete control over your data while providing a modern, user-friendly interface for password management.

## Features ✨

- �� Military-grade AES-256-GCM encryption
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

- Go 1.20 or higher
- Node.js 14 or higher
- PostgreSQL 12 or higher
- npm or yarn

## Installation 🚀

1. Clone the repository:
   ```bash
   git clone https://github.com/yourusername/vault-inator.git
   cd vault-inator
   ```

2. Set up the database:
   ```bash
   # Create a PostgreSQL database
   createdb vaultinator
   ```

3. Configure environment variables:
   ```bash
   # Create a .env file in the root directory
   DATABASE_URL=postgres://username:password@localhost:5432/vaultinator?sslmode=disable
   CONFIG_PATH=/path/to/config.json
   ```

4. Build and run the backend:
   ```bash
   cd backend
   go mod download
   go run cmd/vault-inator/main.go
   ```

5. Build and run the frontend:
   ```bash
   cd web
   npm install
   npm start
   ```

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
├── backend/
│   ├── cmd/
│   │   └── vault-inator/
│   │       └── main.go
│   │
│   ├── internal/
│   │   ├── api/
│   │   ├── config/
│   │   ├── encryption/
│   │   ├── services/
│   │   └── storage/
│   └── go.mod
├── web/
│   ├── src/
│   │   ├── App.js
│   │   └── App.css
│   └── package.json
└── README.md
```

### Running Tests

```bash
# Backend tests
cd backend
go test ./...

# Frontend tests
cd web
npm test
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

## Support 💬

If you encounter any issues or have questions, please open an issue in the GitHub repository.
