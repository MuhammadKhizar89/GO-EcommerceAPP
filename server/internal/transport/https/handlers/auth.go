package handlers

import (
	"encoding/json"
	"net/http"
	"server/internal/domain/auth"
	"server/internal/util/response"
)

type authHandler struct {
	service auth.Service
}

func NewAuthHandler(s auth.Service) *authHandler {
	return &authHandler{service: s}
}

type authRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (h *authHandler) Signup(w http.ResponseWriter, r *http.Request) {
	var req authRequest
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

func (h *authHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req authRequest
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
