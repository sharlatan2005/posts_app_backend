package httphandler

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"

	"github.com/sharlatan2005/chat_app_go_backend_pkg/responseutils"
	"github.com/sharlatan2005/posts_app_backend/services/user/internal/servErrors"
	"github.com/sharlatan2005/posts_app_backend/services/user/internal/service"
)

type UserHandler struct {
	userService *service.UserService
}

func NewUserHandler(userService *service.UserService) *UserHandler {
	return &UserHandler{
		userService: userService,
	}
}

func (h *UserHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		responseutils.StatusMethodNotAllowed(w, "Неправильный метод (должен быть POST)")
		return
	}

	var req struct {
		Username      string `json:"username"`
		Password_hash string `json:"password_hash"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		responseutils.BadRequest(w, "Invalid request body")
		return
	}

	err := h.userService.Create(r.Context(), req.Username, req.Password_hash)
	if err != nil {
		switch {
		case errors.Is(err, servErrors.ErrUsernameEmpty):
			responseutils.BadRequest(w, err.Error())
			return
		case errors.Is(err, servErrors.ErrPasswordEmpty):
			responseutils.BadRequest(w, err.Error())
			return
		case errors.Is(err, servErrors.ErrUserAlreadyExists):
			responseutils.Conflict(w, err.Error())
			return
		default:
			responseutils.InternalServerError(w, "Внутренняя ошибка сервера")
			return
		}
	}

	responseutils.JSON(w, http.StatusCreated, fmt.Sprintf("Пользователь %s успешно сохранен в БД", req.Username))
}

func (h *UserHandler) Exists(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		responseutils.StatusMethodNotAllowed(w, "Неправильный метод (должен быть GET)")
		return
	}

	var req struct {
		Username string `json:"username"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		responseutils.BadRequest(w, "Invalid request body")
		return
	}

	exists, err := h.userService.Exists(r.Context(), req.Username)
	if err != nil {
		switch {
		default:
			responseutils.InternalServerError(w, "Внутренняя ошибка сервера")
			return
		}
	}

	responseutils.JSON(w, http.StatusOK, exists)
}

func (h *UserHandler) GetByUsername(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		responseutils.StatusMethodNotAllowed(w, "Неправильный метод (должен быть GET)")
		return
	}

	username := r.URL.Query().Get("username")

	user, err := h.userService.GetByUsername(r.Context(), username)
	if err != nil {
		switch {
		case errors.Is(err, servErrors.ErrUserNotFound):
			responseutils.NotFound(w, "Нет такого пользователя")
			return
		default:
			log.Println(err.Error())
			responseutils.InternalServerError(w, "Внутренняя ошибка сервера")
			return
		}
	}

	responseutils.JSON(w, http.StatusOK, user)
}
