package handler

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/morozoffnor/home-storage/internal/database"
	"github.com/morozoffnor/home-storage/internal/types"
)

type Item struct {
	db *database.Database
}

func (i *Item) Create(w http.ResponseWriter, r *http.Request) {
	containerID, ok := r.Context().Value("container_id").(int)
	if !ok {
		http.Error(w, http.StatusText(400), 400)
	}
	var req struct {
		Name        string `json:"name"`
		Description string `json:"description"`
		Category    string `json:"category"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	if req.Name == "" {
		http.Error(w, "Name is required", http.StatusBadRequest)
		return
	}
	_, err := i.db.Item.Create(req.Name, req.Description, req.Category, containerID)
	if err != nil {
		fmt.Printf("error creating container: %v\n", err)
		http.Error(w, "Failed to create item", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"message": "Item created successfully"})
}

func (i *Item) GetAllInContainer(w http.ResponseWriter, r *http.Request) {
	containerID, ok := r.Context().Value("container_id").(int)
	if !ok {
		http.Error(w, http.StatusText(400), 400)
	}
	containers, err := i.db.Item.GetAllInContainer(containerID)
	if err != nil {
		fmt.Printf("error getting all items: %v\n", err)
		http.Error(w, "Failed to get items", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(containers)
}

func (i *Item) Get(w http.ResponseWriter, r *http.Request) {
	itemID, ok := r.Context().Value("item_id").(int)
	if !ok {
		http.Error(w, http.StatusText(400), 400)
	}

	home, err := i.db.Item.Get(itemID)
	if err != nil {
		fmt.Printf("error getting item: %v\n", err)
		http.Error(w, "Couldn't get item", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(home)
}

func (i *Item) Update(w http.ResponseWriter, r *http.Request) {
	containerID, ok := r.Context().Value("container_id").(int)
	if !ok {
		http.Error(w, http.StatusText(400), 400)
	}
	itemID, ok := r.Context().Value("item_id").(int)
	if !ok {
		http.Error(w, http.StatusText(400), 400)
	}
	var req struct {
		Name        string `json:"name"`
		Description string `json:"description"`
		Category    string `json:"category"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	item := &types.Item{
		ID:          itemID,
		Name:        req.Name,
		Description: req.Description,
		Category:    req.Category,
		ContainerID: containerID,
	}

	err := i.db.Item.Update(item)
	if err != nil {
		fmt.Printf("error updating item: %v\n", err)
		http.Error(w, "Failed to update item", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"message": "Item updated successfully"})
}

func (i *Item) Delete(w http.ResponseWriter, r *http.Request) {
	itemID, ok := r.Context().Value("item_id").(int)
	if !ok {
		http.Error(w, http.StatusText(400), 400)
	}

	err := i.db.Item.Delete(itemID)
	if err != nil {
		fmt.Printf("error deleting item: %v\n", err)
		http.Error(w, "Failed to delete item", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"message": "Item deleted successfully"})
}
