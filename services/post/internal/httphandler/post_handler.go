package httphandler

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"

	"github.com/google/uuid"
	"github.com/sharlatan2005/chat_app_go_backend_pkg/clients/errorsutils"
	"github.com/sharlatan2005/chat_app_go_backend_pkg/responseutils"
	"github.com/sharlatan2005/posts_app_backend/services/post/internal/servErrors"
	"github.com/sharlatan2005/posts_app_backend/services/post/internal/service"
)

type PostHandler struct {
	postService *service.PostService
}

func NewPostHandler(postService *service.PostService) *PostHandler {
	return &PostHandler{
		postService: postService,
	}
}

func (h *PostHandler) CreatePost(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		responseutils.StatusMethodNotAllowed(w, "Неправильный метод (должен быть POST)")
		return
	}

	var req struct {
		Text string `json:"text"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		responseutils.BadRequest(w, "Invalid request body")
		return
	}

	err := h.postService.Create(r.Context(), req.Text)
	if err != nil {
		switch {
		case errors.Is(err, servErrors.ErrPostTextEmpty):
			responseutils.BadRequest(w, err.Error())
			return
		default:
			log.Println(err.Error())
			responseutils.InternalServerError(w, "Внутренняя ошибка сервера")
			return
		}
	}

	responseutils.JSON(w, http.StatusCreated, "Пост успешно создан.")
}

func (h *PostHandler) DeletePost(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		responseutils.StatusMethodNotAllowed(w, "Неправильный метод (должен быть DELETE)")
		return
	}

	postID := r.URL.Query().Get("post_id")

	err := h.postService.Delete(r.Context(), uuid.MustParse(postID))
	if err != nil {
		var extErr *errorsutils.ExternalServiceError
		switch {
		case errors.Is(err, servErrors.ErrPostNotFound):
			responseutils.NotFound(w, err.Error())
			return
		case errors.Is(err, servErrors.ErrPostForbiddenAction):
			responseutils.Forbidden(w, err.Error())
			return
		case errors.As(err, &extErr):
			responseutils.JSON(w, extErr.StatusCode, extErr.ErrorText)
		default:
			log.Println(err.Error())
			responseutils.InternalServerError(w, "Внутренняя ошибка сервера")
			return
		}
	}

	responseutils.JSON(w, http.StatusOK, "Пост успешно удалён.")
}

func (h *PostHandler) UpdatePost(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
		responseutils.StatusMethodNotAllowed(w, "Неправильный метод (должен быть PUT)")
		return
	}

	var req struct {
		PostID  string `json:"post_id"`
		NewText string `json:"new_text"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		responseutils.BadRequest(w, "Invalid request body")
		return
	}

	err := h.postService.Update(r.Context(), uuid.MustParse(req.PostID), req.NewText)
	if err != nil {
		switch {
		case errors.Is(err, servErrors.ErrPostNotFound):
			responseutils.NotFound(w, err.Error())
			return
		case errors.Is(err, servErrors.ErrPostForbiddenAction):
			responseutils.Forbidden(w, err.Error())
			return
		case errors.Is(err, servErrors.ErrPostTextEmpty):
			responseutils.BadRequest(w, err.Error())
			return
		default:
			log.Println(err.Error())
			responseutils.InternalServerError(w, "Внутренняя ошибка сервера")
			return
		}
	}

	responseutils.JSON(w, http.StatusOK, "Текст поста успешно изменен.")
}

func (h *PostHandler) GetAllUserPosts(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		responseutils.StatusMethodNotAllowed(w, "Неправильный метод (должен быть GET)")
		return
	}

	username := r.URL.Query().Get("username")

	posts, err := h.postService.GetAllUserPosts(r.Context(), username)
	if err != nil {
		var extErr *errorsutils.ExternalServiceError
		switch {
		case errors.As(err, &extErr):
			responseutils.Error(w, extErr.StatusCode, extErr.ErrorText)
			return
		default:
			log.Println(err.Error())
			responseutils.InternalServerError(w, "Внутренняя ошибка сервера")
			return
		}
	}

	responseutils.JSON(w, http.StatusOK, posts)
}
