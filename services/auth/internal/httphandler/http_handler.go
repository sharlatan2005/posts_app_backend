package httphandler

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"

	"github.com/sharlatan2005/chat_app_go_backend_pkg/clients"
	"github.com/sharlatan2005/chat_app_go_backend_pkg/responseutils"
	"github.com/sharlatan2005/posts_app_backend/services/auth/internal/service"
)

type AuthHandler struct {
	authService *service.AuthService
}

func NewAuthHandler(authService *service.AuthService) *AuthHandler {
	return &AuthHandler{
		authService: authService,
	}
}

func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		responseutils.StatusMethodNotAllowed(w, "Неправильный метод. Нужен POST")
	}

	var req struct {
		Username     string `json:"username"`
		PasswordHash string `json:"password_hash"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		responseutils.BadRequest(w, "Invalid request body")
		return
	}

	result, err := h.authService.RegisterUser(r.Context(), req.Username, req.PasswordHash)
	if err != nil {
		var extErr *clients.ExternalServiceError
		switch {
		case errors.As(err, &extErr):
			responseutils.Error(w, extErr.StatusCode, extErr.ErrorText)
			return
		default:
			log.Printf("Unhandled error: %+v", err)
			responseutils.InternalServerError(w, "Внутренняя ошибка сервера")
			return
		}
	}

	responseutils.JSON(w, http.StatusCreated, result)
}
