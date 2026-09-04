package handler

import (
	"encoding/json"
	"net/http"

	"github.com/SephirothGit/warehouse/internal/service"
)

type AuthHandler struct {
	authService service.AuthService
}

type refreshRequest struct {
	RefreshToken string `json:"refresh_token"`
}

type registerRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type loginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func NewAuthHandler(authService service.AuthService) *AuthHandler {
	return &AuthHandler{
		authService: authService,
	}
}

func (h *AuthHandler) RegisterHandler(w http.ResponseWriter, r *http.Request) {
	var req registerRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	userID, err := h.authService.Register(req.Email, req.Password)
	if err != nil {
		http.Error(w, "unable to register user", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]int{"user_id": userID})
}

func (h *AuthHandler) LoginHandler(w http.ResponseWriter, r *http.Request) {
	var req loginRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "unable to login user", http.StatusBadRequest)
		return
	}

	tokens, err := h.authService.Login(req.Email, req.Password)
	if err != nil {
		http.Error(w, "invalid email or password", http.StatusUnauthorized)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(tokens)
}

func (h *AuthHandler) RefreshHandler(w http.ResponseWriter, r *http.Request) {
	var req refreshRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	tokens, err := h.authService.RefreshAccessToken(req.RefreshToken)
	if err != nil {
		http.Error(w, "invalid or expired refresh token", http.StatusUnauthorized)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(tokens)
}

func (h *AuthHandler) LogoutHandler(w http.ResponseWriter, r *http.Request) {
	var req refreshRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	err = h.authService.Logout(req.RefreshToken)
	if err != nil {
		http.Error(w, "unable to logout", http.StatusUnauthorized)
		return
	}

	w.WriteHeader(http.StatusOK)
}
