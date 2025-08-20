package handler

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/morozoffnor/home-storage/internal/database"
	"github.com/morozoffnor/home-storage/internal/types"
)

type Container struct {
	db *database.Database
}

func (c *Container) Create(w http.ResponseWriter, r *http.Request) {
	homeID, ok := r.Context().Value("home_id").(int)
	if !ok {
		http.Error(w, http.StatusText(400), 400)
	}
	var req struct {
		Name        string `json:"name"`
		Description string `json:"description"`
		Category    string `json:"category"`
		Location    string `json:"location"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	if req.Name == "" {
		http.Error(w, "Name is required", http.StatusBadRequest)
		return
	}
	_, err := c.db.Container.Create(req.Name, req.Description, req.Category, req.Location, homeID)
	if err != nil {
		fmt.Printf("error creating container: %v\n", err)
		http.Error(w, "Failed to create container", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"message": "Container created successfully"})
}

func (c *Container) GetAllInHome(w http.ResponseWriter, r *http.Request) {
	homeID, ok := r.Context().Value("home_id").(int)
	if !ok {
		http.Error(w, http.StatusText(400), 400)
	}
	containers, err := c.db.Container.GetAllInHome(homeID)
	if err != nil {
		fmt.Printf("error getting all containers: %v\n", err)
		http.Error(w, "Failed to get containers", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(containers)
}

func (c *Container) Get(w http.ResponseWriter, r *http.Request) {
	containerID, ok := r.Context().Value("container_id").(int)
	if !ok {
		http.Error(w, http.StatusText(400), 400)
	}

	home, err := c.db.Container.Get(containerID)
	if err != nil {
		fmt.Printf("error getting containers: %v\n", err)
		http.Error(w, "Home not containers", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(home)
}

func (c *Container) Update(w http.ResponseWriter, r *http.Request) {
	homeID, ok := r.Context().Value("home_id").(int)
	if !ok {
		http.Error(w, http.StatusText(400), 400)
	}
	containerID, ok := r.Context().Value("container_id").(int)
	if !ok {
		http.Error(w, http.StatusText(400), 400)
	}
	var req struct {
		Name        string `json:"name"`
		Description string `json:"description"`
		Category    string `json:"category"`
		Location    string `json:"location"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	home := &types.Container{
		ID:          containerID,
		Name:        req.Name,
		Description: req.Description,
		Category:    req.Category,
		Location:    req.Location,
		HomeID:      homeID,
	}

	err := c.db.Container.Update(home)
	if err != nil {
		fmt.Printf("error updating container: %v\n", err)
		http.Error(w, "Failed to update container", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"message": "Container updated successfully"})
}
