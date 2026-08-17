package httphandler

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"

	"github.com/google/uuid"
	"github.com/sharlatan2005/chat_app_go_backend_pkg/responseutils"
	"github.com/sharlatan2005/posts_app_backend/services/like/internal/servErrors"
	"github.com/sharlatan2005/posts_app_backend/services/like/internal/service"
)

type LikeHandler struct {
	likeService *service.LikeService
}

func NewLikeHandler(likeService *service.LikeService) *LikeHandler {
	return &LikeHandler{
		likeService: likeService,
	}
}

func (h *LikeHandler) AddLike(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		responseutils.StatusMethodNotAllowed(w, "Неправильный метод (должен быть POST)")
		return
	}

	var req struct {
		PostID string `json:"post_id"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		responseutils.BadRequest(w, "Invalid request body")
		return
	}

	err := h.likeService.AddLike(r.Context(), uuid.MustParse(req.PostID))

	if err != nil {
		switch {
		default:
			log.Println(err.Error())
			responseutils.InternalServerError(w, "Внутренняя ошибка сервера")
			return
		}
	}

	responseutils.JSON(w, http.StatusCreated, "Лайк успешно поставлен.")
}

func (h *LikeHandler) RemoveLike(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		responseutils.StatusMethodNotAllowed(w, "Неправильный метод (должен быть DELETE)")
		return
	}

	postID := r.URL.Query().Get("post_id")

	err := h.likeService.RemoveLike(r.Context(), uuid.MustParse(postID))
	if err != nil {
		switch {
		case errors.Is(err, servErrors.ErrLikeNotFound):
			responseutils.NotFound(w, err.Error())
			return
		case errors.Is(err, servErrors.ErrLikeForbiddenAction):
			responseutils.Forbidden(w, err.Error())
			return
		default:
			log.Println(err.Error())
			responseutils.InternalServerError(w, "Внутренняя ошибка сервера")
			return
		}
	}

	responseutils.JSON(w, http.StatusCreated, "Лайк успешно убран с поста.")
}

func (h *LikeHandler) GetLikedUserIDs(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		responseutils.StatusMethodNotAllowed(w, "Неправильный метод (должен быть GET)")
		return
	}

	postID := r.URL.Query().Get("post_id")

	likedUserIDs, err := h.likeService.GetLikedUserIDs(r.Context(), uuid.MustParse(postID))
	if err != nil {
		switch {
		default:
			log.Println(err.Error())
			responseutils.InternalServerError(w, "Внутренняя ошибка сервера")
			return
		}
	}

	responseutils.JSON(w, http.StatusCreated, likedUserIDs)
}

func (h *LikeHandler) DeleteLikesByPost(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		responseutils.StatusMethodNotAllowed(w, "Неправильный метод (должен быть DELETE)")
		return
	}

	postID := r.URL.Query().Get("post_id")

	err := h.likeService.DeleteAllLikesByPost(r.Context(), uuid.MustParse(postID))
	if err != nil {
		switch {
		case errors.Is(err, servErrors.ErrNoLikesOnPost):
			responseutils.NotFound(w, err.Error())
			return
		default:
			log.Println(err.Error())
			responseutils.InternalServerError(w, "Внутренняя ошибка сервера")
			return
		}
	}

	responseutils.JSON(w, http.StatusCreated, "Лайки успешно удалены с поста.")
}
