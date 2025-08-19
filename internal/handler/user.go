package handler

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/morozoffnor/home-storage/internal/auth"
	"github.com/morozoffnor/home-storage/internal/database"
)

type User struct {
	auth *auth.Auth
	db   *database.Database
}

func (u *User) Register(w http.ResponseWriter, r *http.Request) {
	var raw bytes.Buffer
	if _, err := raw.ReadFrom(r.Body); err != nil {
		http.Error(w, "Invalid body", http.StatusUnprocessableEntity)
		return
	}

	type regReq struct {
		Username string `json:"username"`
		Password string `json:"password"`
		Email    string `json:"email"`
	}

	user := regReq{}

	err := json.Unmarshal(raw.Bytes(), &user)
	if err != nil {
		fmt.Println(err)
		http.Error(w, "Invalid json", http.StatusUnprocessableEntity)
		return
	}

	exists, err := u.db.User.Exists(user.Email)
	if err != nil {
		fmt.Println(err)
		http.Error(w, "Internal error", http.StatusInternalServerError)
		return
	}
	if exists {
		http.Error(w, "User already exists", http.StatusConflict)
		return
	}

	// generating new token
	token, err := u.auth.Jwt.GenerateToken(user.Email)
	if err != nil {
		fmt.Println(err)
		http.Error(w, "Internal error", http.StatusInternalServerError)
		return
	}

	// adding token to cookies
	ctx, err := u.auth.Jwt.AddTokenToCookies(&w, r, token)
	if err != nil {
		fmt.Println(err)
		return
	}
	r = r.WithContext(ctx)

	// hashing password and saving user
	hash := u.auth.HashPassword(user.Password)
	err = u.db.User.Create(user.Username, user.Email, hash)
	if err != nil {
		fmt.Println(err)
		http.Error(w, "Internal error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/plain, utf-8")
	w.WriteHeader(http.StatusOK)
}

func (u *User) Login(w http.ResponseWriter, r *http.Request) {
	var raw bytes.Buffer
	if _, err := raw.ReadFrom(r.Body); err != nil {
		http.Error(w, "Invalid body", http.StatusUnprocessableEntity)
		return
	}
	type LoginReq struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	loginReq := LoginReq{}

	err := json.Unmarshal(raw.Bytes(), &loginReq)
	if err != nil {
		fmt.Println(err)
		http.Error(w, "Invalid json", http.StatusUnprocessableEntity)
		return
	}

	dbUser, err := u.db.User.Get(loginReq.Email)
	if err != nil {
		fmt.Println(err)
		http.Error(w, "No such user", http.StatusUnauthorized)
		return
	}

	hash := u.auth.HashPassword(loginReq.Password)
	if dbUser.PassHash != hash {
		http.Error(w, "Wrong login/pass", http.StatusUnauthorized)
		return
	}

	token, err := u.auth.Jwt.GenerateToken(dbUser.Email)
	if err != nil {
		fmt.Println(err)
		http.Error(w, "Internal error", http.StatusInternalServerError)
		return
	}
	ctx, err := u.auth.Jwt.AddTokenToCookies(&w, r, token)
	if err != nil {
		fmt.Println(err)
		return
	}

	r = r.WithContext(ctx)

	w.Header().Set("Content-Type", "text/plain, utf-8")
	w.WriteHeader(http.StatusOK)
}

func (u *User) Get(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value("user_id").(int)
	if !ok {
		http.Error(w, http.StatusText(400), 400)
		return
	}

	user, err := u.db.User.GetByID(userID)
	if err != nil {
		fmt.Printf("error getting user: %v\n", err)
		http.Error(w, "Failed to get user", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}

func (u *User) GetAll(w http.ResponseWriter, r *http.Request) {
	users, err := u.db.User.GetAll()
	if err != nil {
		fmt.Println("error getting all users")
		http.Error(w, http.StatusText(500), 500)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(users)
}

func (u *User) AddHome(w http.ResponseWriter, r *http.Request) {
	type addHomeReq struct {
		ID int `json:"id"`
	}
	userID, ok := r.Context().Value("user_id").(int)
	if !ok {
		http.Error(w, http.StatusText(400), 400)
		return
	}
	var raw bytes.Buffer
	if _, err := raw.ReadFrom(r.Body); err != nil {
		http.Error(w, "Invalid body", http.StatusUnprocessableEntity)
		return
	}
	homeReq := addHomeReq{}
	err := json.Unmarshal(raw.Bytes(), &homeReq)
	if err != nil {
		fmt.Println(err)
		http.Error(w, "Invalid json", http.StatusUnprocessableEntity)
		return
	}

	err = u.db.User.AddHome(userID, homeReq.ID)
	if err != nil {
		fmt.Println(err)
		http.Error(w, "Error adding a home", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func (u *User) GetHomes(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value("user_id").(int)
	if !ok {
		http.Error(w, http.StatusText(400), 400)
		return
	}
	homes, err := u.db.User.GetHomes(userID)
	if err != nil {
		fmt.Println(err)
		http.Error(w, "Error finding homes", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(homes)
}
