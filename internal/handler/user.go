package handler

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/morozoffnor/home-storage/internal/types"
)

func (h *APIHandler) RegisterUser(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		fmt.Println(err)
		http.Error(w, "Failed parsing body", http.StatusBadRequest)
		return
	}

	var user types.User

	err = json.Unmarshal(body, &user)
	if err != nil {
		fmt.Println(err)
		http.Error(w, "Invalid json", http.StatusUnprocessableEntity)
		return
	}

	exists, err := h.db.User.Exists(user.Email)
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
	token, err := h.auth.Jwt.GenerateToken(user.Email)
	if err != nil {
		fmt.Println(err)
		http.Error(w, "Internal error", http.StatusInternalServerError)
		return
	}

	// adding token to cookies
	ctx, err := h.auth.Jwt.AddTokenToCookies(&w, r, token)
	if err != nil {
		fmt.Println(err)
		return
	}
	r = r.WithContext(ctx)

	// saving user
	err = h.db.User.Create(user.Username, user.Email, user.PassHash)
	if err != nil {
		fmt.Println(err)
		http.Error(w, "Internal error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/plain, utf-8")
	w.WriteHeader(http.StatusOK)
}

func (h *APIHandler) LoginUser(w http.ResponseWriter, r *http.Request) {
	// получаем тело, распаковываем его, если надо
	body, err := io.ReadAll(r.Body)
	if err != nil {
		fmt.Println(err)
		http.Error(w, "Failed parsing body", http.StatusBadRequest)
		return
	}
	type LoginReq struct {
		Email    string `json:"email"`
		PassHash string `json:"password_hash"`
	}

	var loginReq = LoginReq{}

	// вытаскиваем джейсонку в структуру
	err = json.Unmarshal(body, &loginReq)
	if err != nil {
		fmt.Println(err)
		http.Error(w, "Invalid json", http.StatusUnprocessableEntity)
		return
	}

	// находим такого юзера в бд
	dbUser, err := h.db.User.Get(loginReq.Email)
	if err != nil {
		fmt.Println(err)
		http.Error(w, "No such user", http.StatusUnauthorized)
		return
	}

	// хешируем присланный пароль, сравниваем с хешем из бд
	if dbUser.PassHash != loginReq.PassHash {
		http.Error(w, "Wrong login/pass", http.StatusUnauthorized)
		return
	}

	// добавляем токен в куки
	token, err := h.auth.Jwt.GenerateToken(dbUser.Email)
	if err != nil {
		fmt.Println(err)
		http.Error(w, "Internal error", http.StatusInternalServerError)
		return
	}
	ctx, err := h.auth.Jwt.AddTokenToCookies(&w, r, token)
	if err != nil {
		fmt.Println(err)
		return
	}

	r = r.WithContext(ctx)

	w.Header().Set("Content-Type", "text/plain, utf-8")
	w.WriteHeader(http.StatusOK)
}
