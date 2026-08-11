package main

import (
	"log"
	"net/http"

	db "github.com/tushar7789/E_commerce_ms/cmd/users_service/data/db/out"
	"github.com/tushar7789/E_commerce_ms/internal/json"
)

type handler struct {
	service Service
}

type UserRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type RefreshRequest struct {
	RefreshToken string `json:"refresh_token"`
}

type UserResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

func NewHandler(service Service) *handler {
	return &handler{
		service: service,
	}
}

func (h *handler) CreateUser(w http.ResponseWriter, r *http.Request) {
	var reqData UserRequest

	if !json.WriteData(w, r, &reqData) {
		return
	}

	passHash, err := HashPassword(reqData.Password)
	if err != nil {
		log.Print(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	params := db.CreateUserParams{
		Username:     reqData.Username,
		PasswordHash: passHash,
	}

	user, token, err := h.service.CreateUser(r.Context(), params)
	if err != nil {
		log.Print(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.Write(w, http.StatusCreated, struct {
		User   interface{} `json:"user"`
		Tokens TokenPair   `json:"tokens"`
	}{
		User:   user,
		Tokens: token,
	})
}

func (h *handler) Login(w http.ResponseWriter, r *http.Request) {
	var reqData UserRequest

	if !json.WriteData(w, r, &reqData) {
		return
	}

	tokens, err := h.service.Login(r.Context(), reqData.Username, reqData.Password)
	if err != nil {
		http.Error(w, "invalid credentials", http.StatusUnauthorized)
		return
	}

	json.Write(w, http.StatusOK, UserResponse{
		AccessToken:  tokens.AccessToken,
		RefreshToken: tokens.RefreshToken,
	})
}

func (h *handler) Refresh(w http.ResponseWriter, r *http.Request) {
	var reqData RefreshRequest

	if !json.WriteData(w, r, &reqData) {
		return
	}

	tokens, err := h.service.Refresh(reqData.RefreshToken)
	if err != nil {
		http.Error(w, "invalid refresh token", http.StatusUnauthorized)
		return
	}

	json.Write(w, http.StatusOK, UserResponse{
		AccessToken:  tokens.AccessToken,
		RefreshToken: tokens.RefreshToken,
	})
}
