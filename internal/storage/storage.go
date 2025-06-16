package storage

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/google/uuid"
	_ "github.com/lib/pq"
	"github.com/nonaxanon/vault-inator/config"
	"github.com/nonaxanon/vault-inator/internal/encryption"
)

// PasswordEntry represents a stored password entry.
type PasswordEntry struct {
	ID       uuid.UUID
	Title    string
	Username string
	Password string
	Notes    string
}

// DB holds the database connection and encryption.
type DB struct {
	*sql.DB
	encryptor *encryption.Encryptor
}

// NewDB creates a new database connection using the provided connection string and master password.
func NewDB(connStr string, masterPassword string) (*DB, error) {
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, err
	}
	if err := db.Ping(); err != nil {
		return nil, err
	}

	encryptor, err := encryption.NewEncryptor(masterPassword)
	if err != nil {
		return nil, fmt.Errorf("failed to create encryptor: %v", err)
	}

	return &DB{db, encryptor}, nil
}

// InitDB initializes the database by creating the vaultinator schema and the passwords table if they don't exist.
func (db *DB) InitDB() error {
	// Create the vaultinator schema if it doesn't exist
	createSchemaQuery := `CREATE SCHEMA IF NOT EXISTS vaultinator;`
	_, err := db.Exec(createSchemaQuery)
	if err != nil {
		return err
	}

	// Set the search path to the vaultinator schema
	setSearchPathQuery := `SET search_path TO vaultinator;`
	_, err = db.Exec(setSearchPathQuery)
	if err != nil {
		return err
	}

	// Enable the uuid-ossp extension
	enableUUIDQuery := `CREATE EXTENSION IF NOT EXISTS "uuid-ossp";`
	_, err = db.Exec(enableUUIDQuery)
	if err != nil {
		return err
	}

	// Create the passwords table in the vaultinator schema
	createTableQuery := `
	CREATE TABLE IF NOT EXISTS vaultinator.passwords (
		id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
		title TEXT NOT NULL,
		username TEXT NOT NULL,
		password TEXT NOT NULL,
		notes TEXT
	);`
	_, err = db.Exec(createTableQuery)
	return err
}

// AddPassword adds a new password entry to the database.
func (db *DB) AddPassword(entry PasswordEntry) error {
	// Encrypt the password before storing
	encryptedPassword, err := db.encryptor.Encrypt(entry.Password)
	if err != nil {
		return fmt.Errorf("failed to encrypt password: %v", err)
	}

	query := `
	INSERT INTO vaultinator.passwords (id, title, username, password, notes)
	VALUES ($1, $2, $3, $4, $5)
	RETURNING id;`
	var id uuid.UUID
	err = db.QueryRow(query, uuid.New(), entry.Title, entry.Username, encryptedPassword, entry.Notes).Scan(&id)
	if err != nil {
		return err
	}
	log.Printf("Added password entry with ID: %s", id)
	return nil
}

// GetPassword retrieves a password entry by its ID.
func (db *DB) GetPassword(id uuid.UUID) (PasswordEntry, error) {
	var entry PasswordEntry
	var encryptedPassword string
	query := `SELECT id, title, username, password, notes FROM vaultinator.passwords WHERE id = $1;`
	err := db.QueryRow(query, id).Scan(&entry.ID, &entry.Title, &entry.Username, &encryptedPassword, &entry.Notes)
	if err != nil {
		return PasswordEntry{}, err
	}

	// Decrypt the password
	decryptedPassword, err := db.encryptor.Decrypt(encryptedPassword)
	if err != nil {
		return PasswordEntry{}, fmt.Errorf("failed to decrypt password: %v", err)
	}
	entry.Password = decryptedPassword

	return entry, nil
}

// GetAllPasswords retrieves all password entries from the database.
func (db *DB) GetAllPasswords() ([]PasswordEntry, error) {
	query := `SELECT id, title, username, password, notes FROM vaultinator.passwords;`
	rows, err := db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var entries []PasswordEntry
	for rows.Next() {
		var entry PasswordEntry
		var encryptedPassword string
		if err := rows.Scan(&entry.ID, &entry.Title, &entry.Username, &encryptedPassword, &entry.Notes); err != nil {
			return nil, err
		}

		// Decrypt the password
		decryptedPassword, err := db.encryptor.Decrypt(encryptedPassword)
		if err != nil {
			return nil, fmt.Errorf("failed to decrypt password: %v", err)
		}
		entry.Password = decryptedPassword

		entries = append(entries, entry)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return entries, nil
}

// DeletePassword deletes a password entry by its ID.
func (db *DB) DeletePassword(id uuid.UUID) error {
	query := `DELETE FROM vaultinator.passwords WHERE id = $1;`
	result, err := db.Exec(query, id)
	if err != nil {
		return err
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return fmt.Errorf("no password entry found with ID: %s", id)
	}
	log.Printf("Deleted password entry with ID: %s", id)
	return nil
}

// UpdateMasterPassword updates the master password and re-encrypts all stored passwords.
func (db *DB) UpdateMasterPassword(currentPassword, newPassword string) error {
	// Create new encryptor with the new password
	newEncryptor, err := encryption.NewEncryptor(newPassword)
	if err != nil {
		return fmt.Errorf("failed to create new encryptor: %v", err)
	}

	// Get all passwords
	entries, err := db.GetAllPasswords()
	if err != nil {
		return fmt.Errorf("failed to get passwords: %v", err)
	}

	// Begin transaction
	tx, err := db.Begin()
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %v", err)
	}
	defer tx.Rollback()

	// Update each password with new encryption
	for _, entry := range entries {
		// Re-encrypt the password with the new master password
		encryptedPassword, err := newEncryptor.Encrypt(entry.Password)
		if err != nil {
			return fmt.Errorf("failed to re-encrypt password: %v", err)
		}

		// Update the password in the database
		query := `UPDATE vaultinator.passwords SET password = $1 WHERE id = $2;`
		if _, err := tx.Exec(query, encryptedPassword, entry.ID); err != nil {
			return fmt.Errorf("failed to update password: %v", err)
		}
	}

	// Update the master password in the configuration
	cfg, err := config.LoadConfig()
	if err != nil {
		return fmt.Errorf("failed to load config: %v", err)
	}

	cfg.MasterPassword = newPassword
	if err := cfg.SaveConfig(); err != nil {
		return fmt.Errorf("failed to save config: %v", err)
	}

	// Commit transaction
	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %v", err)
	}

	// Update the encryptor in the DB struct
	db.encryptor = newEncryptor

	return nil
}
