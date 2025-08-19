package handler

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/morozoffnor/home-storage/internal/auth"
	"github.com/morozoffnor/home-storage/internal/config"
	"github.com/morozoffnor/home-storage/internal/database"
	"github.com/morozoffnor/home-storage/internal/types"
)

type APIHandler struct {
	cfg  *config.Config
	db   *database.Database
	auth *auth.Auth
}

func New(cfg *config.Config, db *database.Database, a *auth.Auth) *APIHandler {
	return &APIHandler{
		cfg:  cfg,
		db:   db,
		auth: a,
	}
}

func (h *APIHandler) CreateHome(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Name        string `json:"name"`
		Description string `json:"description"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	if req.Name == "" {
		http.Error(w, "Name is required", http.StatusBadRequest)
		return
	}

	err := h.db.Home.Create(req.Name, req.Description)
	if err != nil {
		fmt.Printf("error creating home: %v\n", err)
		http.Error(w, "Failed to create home", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"message": "Home created successfully"})
}

func (h *APIHandler) GetAllHomes(w http.ResponseWriter, r *http.Request) {
	homes, err := h.db.Home.GetAll()
	if err != nil {
		fmt.Printf("error getting all homes: %v\n", err)
		http.Error(w, "Failed to get homes", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(homes)
}

func (h *APIHandler) GetHome(w http.ResponseWriter, r *http.Request) {
	homeID, ok := r.Context().Value("home_id").(int)
	if !ok {
		http.Error(w, http.StatusText(400), 400)
	}
	fmt.Println(homeID)

	home, err := h.db.Home.Get(homeID)
	if err != nil {
		fmt.Printf("error getting home: %v\n", err)
		http.Error(w, "Home not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(home)
}

func (h *APIHandler) UpdateHome(w http.ResponseWriter, r *http.Request) {
	homeID, ok := r.Context().Value("home_id").(int)
	if !ok {
		http.Error(w, http.StatusText(400), 400)
	}
	var req struct {
		Name        string `json:"name"`
		Description string `json:"description"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	home := &types.Home{
		ID:          homeID,
		Name:        req.Name,
		Description: req.Description,
	}

	err := h.db.Home.Update(home)
	if err != nil {
		fmt.Printf("error updating home: %v\n", err)
		http.Error(w, "Failed to update home", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"message": "Home updated successfully"})
}

func (h *APIHandler) DeleteHome(w http.ResponseWriter, r *http.Request) {
	homeID, ok := r.Context().Value("home_id").(int)
	if !ok {
		http.Error(w, http.StatusText(400), 400)
	}

	err := h.db.Home.Delete(homeID)
	if err != nil {
		fmt.Printf("error deleting home: %v\n", err)
		http.Error(w, "Failed to delete home", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"message": "Home deleted successfully"})
}
