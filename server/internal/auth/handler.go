package auth

import (
	"encoding/json"
	"net/http"
	"server/internal/response"
)

type AuthHandler struct {
	service Service
}

func NewAuthHandler(s Service) *AuthHandler {
	return &AuthHandler{service: s}
}

type AuthRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (h *AuthHandler) Signup(w http.ResponseWriter, r *http.Request) {
	var req AuthRequest
	json.NewDecoder(r.Body).Decode(&req)

	token, err := h.service.Signup(r.Context(), req.Email, req.Password)
	if err != nil {
		response.WriteJson(w, http.StatusBadRequest, response.GernalResponse{
			Success: false,
			Message: err.Error(),
			Data:    nil,
		})
		return
	}

	response.WriteJson(w, http.StatusOK, response.GernalResponse{
		Success: true,
		Message: "User created successfully",
		Data: map[string]string{
			"token": token,
		},
	})
}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req AuthRequest
	json.NewDecoder(r.Body).Decode(&req)

	token, err := h.service.Login(r.Context(), req.Email, req.Password)

	if err != nil {
		response.WriteJson(w, http.StatusBadRequest, response.GernalResponse{
			Success: false,
			Message: err.Error(),
			Data:    nil,
		})
		return
	}

	response.WriteJson(w, http.StatusOK, response.GernalResponse{
		Success: true,
		Message: "success",
		Data: map[string]string{
			"token": token,
		},
	})
}
