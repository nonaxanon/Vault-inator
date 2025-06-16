package api

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/nonaxanon/vault-inator/internal/services"
	"github.com/nonaxanon/vault-inator/internal/storage"
	"github.com/rs/cors"
	"github.com/sirupsen/logrus"
)

// Server holds the API server dependencies.
type Server struct {
	router          *mux.Router
	db              *storage.DB
	logger          *logrus.Logger
	authService     *services.AuthService
	passwordService *services.PasswordService
}

// NewServer creates a new API server with the provided database connection and services.
func NewServer(db *storage.DB, authService *services.AuthService, passwordService *services.PasswordService) *Server {
	logger := logrus.New()
	logger.SetFormatter(&logrus.JSONFormatter{})

	s := &Server{
		router:          mux.NewRouter(),
		db:              db,
		logger:          logger,
		authService:     authService,
		passwordService: passwordService,
	}
	s.routes()
	return s
}

// routes sets up the API routes.
func (s *Server) routes() {
	// Auth endpoints
	s.router.HandleFunc("/api/auth/initialize", s.handleInitializeMasterPassword).Methods("POST")
	s.router.HandleFunc("/api/auth/verify", s.handleVerifyMasterPassword).Methods("POST")
	s.router.HandleFunc("/api/auth/change", s.handleChangeMasterPassword).Methods("POST")
	s.router.HandleFunc("/api/auth/status", s.handleAuthStatus).Methods("GET")

	// Password endpoints
	s.router.HandleFunc("/api/passwords", s.handleAddPassword).Methods("POST")
	s.router.HandleFunc("/api/passwords", s.handleGetAllPasswords).Methods("GET")
	s.router.HandleFunc("/api/passwords/{id}", s.handleGetPassword).Methods("GET")
	s.router.HandleFunc("/api/passwords/{id}", s.handleDeletePassword).Methods("DELETE")

	// Serve static files (React frontend)
	s.router.PathPrefix("/").Handler(http.FileServer(http.Dir("./web/build")))
}

// ServeHTTP implements the http.Handler interface.
func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// Add CORS middleware
	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"http://localhost:5432", "http://localhost:3000"},
		AllowedMethods:   []string{"GET", "POST", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Content-Type"},
		AllowCredentials: true,
	})
	c.Handler(s.router).ServeHTTP(w, r)
}

// handleInitializeMasterPassword handles the POST request to initialize the master password.
func (s *Server) handleInitializeMasterPassword(w http.ResponseWriter, r *http.Request) {
	s.logger.Info("Received POST request to /api/auth/initialize")
	var req struct {
		Password string `json:"password" binding:"required"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		s.logger.WithError(err).Error("Error decoding request body")
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := s.authService.InitializeMasterPassword(req.Password); err != nil {
		s.logger.WithError(err).Error("Error initializing master password")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	s.logger.Info("Successfully initialized master password")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Master password initialized"})
}

// handleVerifyMasterPassword handles the POST request to verify the master password.
func (s *Server) handleVerifyMasterPassword(w http.ResponseWriter, r *http.Request) {
	s.logger.Info("Received POST request to /api/auth/verify")
	var req struct {
		Password string `json:"password" binding:"required"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		s.logger.WithError(err).Error("Error decoding request body")
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if !s.authService.VerifyMasterPassword(req.Password) {
		s.logger.Error("Invalid master password")
		http.Error(w, "Invalid password", http.StatusUnauthorized)
		return
	}

	s.logger.Info("Successfully verified master password")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Password verified"})
}

// handleChangeMasterPassword handles the POST request to change the master password.
func (s *Server) handleChangeMasterPassword(w http.ResponseWriter, r *http.Request) {
	s.logger.Info("Received POST request to /api/auth/change")
	var req struct {
		CurrentPassword string `json:"currentPassword" binding:"required"`
		NewPassword     string `json:"newPassword" binding:"required"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		s.logger.WithError(err).Error("Error decoding request body")
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := s.authService.ChangeMasterPassword(req.CurrentPassword, req.NewPassword); err != nil {
		s.logger.WithError(err).Error("Error changing master password")
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	s.logger.Info("Successfully changed master password")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Password changed successfully"})
}

// handleAuthStatus handles the GET request to check authentication status.
func (s *Server) handleAuthStatus(w http.ResponseWriter, r *http.Request) {
	s.logger.Info("Received GET request to /api/auth/status")
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]bool{
		"initialized": s.authService.IsInitialized(),
	})
}

// handleAddPassword handles the POST request to add a new password entry.
func (s *Server) handleAddPassword(w http.ResponseWriter, r *http.Request) {
	s.logger.Info("Received POST request to /api/passwords")
	var password services.Password
	if err := json.NewDecoder(r.Body).Decode(&password); err != nil {
		s.logger.WithError(err).Error("Error decoding request body")
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := s.passwordService.CreatePassword(&password); err != nil {
		s.logger.WithError(err).Error("Error adding password")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	s.logger.WithField("title", password.Title).Info("Successfully added password entry")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(password)
}

// handleGetAllPasswords handles the GET request to retrieve all password entries.
func (s *Server) handleGetAllPasswords(w http.ResponseWriter, r *http.Request) {
	s.logger.Info("Received GET request to /api/passwords")
	passwords, err := s.passwordService.GetAllPasswords()
	if err != nil {
		s.logger.WithError(err).Error("Error fetching passwords")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	s.logger.WithField("count", len(passwords)).Info("Successfully fetched password entries")
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(passwords)
}

// handleGetPassword handles the GET request to retrieve a specific password entry by ID.
func (s *Server) handleGetPassword(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	s.logger.WithField("id", id).Info("Received GET request to /api/passwords/{id}")
	passwords, err := s.passwordService.GetAllPasswords()
	if err != nil {
		s.logger.WithError(err).Error("Error fetching passwords")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Find the password with the matching ID
	for _, p := range passwords {
		if p.ID == id {
			s.logger.WithField("id", id).Info("Successfully fetched password entry")
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(p)
			return
		}
	}

	s.logger.WithField("id", id).Error("Password not found")
	http.Error(w, "Password not found", http.StatusNotFound)
}

// handleDeletePassword handles the DELETE request to remove a password entry by ID.
func (s *Server) handleDeletePassword(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	s.logger.WithField("id", id).Info("Received DELETE request to /api/passwords/{id}")
	if err := s.passwordService.DeletePassword(id); err != nil {
		s.logger.WithError(err).WithField("id", id).Error("Error deleting password")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	s.logger.WithField("id", id).Info("Successfully deleted password entry")
	w.WriteHeader(http.StatusNoContent)
}
