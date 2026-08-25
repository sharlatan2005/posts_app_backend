package httphandler

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"

	"github.com/google/uuid"
	"github.com/sharlatan2005/chat_app_go_backend_pkg/responseutils"
	"github.com/sharlatan2005/posts_app_backend/services/comment/internal/servErrors"
	"github.com/sharlatan2005/posts_app_backend/services/comment/internal/service"
)

type CommentHandler struct {
	commentService *service.CommentService
}

func NewCommentHandler(commentService *service.CommentService) *CommentHandler {
	return &CommentHandler{
		commentService: commentService,
	}
}

func (h *CommentHandler) CreateComment(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		responseutils.StatusMethodNotAllowed(w, "Неправильный метод (должен быть POST)")
		return
	}

	var req struct {
		PostID string `json:"post_id"`
		Text   string `json:"text"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		responseutils.BadRequest(w, "Invalid request body")
		return
	}

	err := h.commentService.Create(r.Context(), uuid.MustParse(req.PostID), req.Text)

	if err != nil {
		switch {
		case errors.Is(err, servErrors.ErrCommentTextEmpty):
			responseutils.BadRequest(w, err.Error())
			return
		default:
			log.Println(err.Error())
			responseutils.InternalServerError(w, "Внутренняя ошибка сервера")
			return
		}
	}

	responseutils.JSON(w, http.StatusCreated, "Комментарий успешно создан.")
}

func (h *CommentHandler) DeleteComment(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		responseutils.StatusMethodNotAllowed(w, "Неправильный метод (должен быть DELETE)")
		return
	}

	commentID := r.URL.Query().Get("comment_id")

	err := h.commentService.Delete(r.Context(), uuid.MustParse(commentID))
	if err != nil {
		switch {
		case errors.Is(err, servErrors.ErrCommentNotFound):
			responseutils.NotFound(w, err.Error())
			return
		case errors.Is(err, servErrors.ErrCommentForbiddenAction):
			responseutils.Forbidden(w, err.Error())
			return
		default:
			log.Println(err.Error())
			responseutils.InternalServerError(w, "Внутренняя ошибка сервера")
			return
		}
	}

	responseutils.JSON(w, http.StatusOK, "Комментарий успешно удалён.")
}

func (h *CommentHandler) UpdateComment(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
		responseutils.StatusMethodNotAllowed(w, "Неправильный метод (должен быть PUT)")
		return
	}

	var req struct {
		CommentID string `json:"comment_id"`
		NewText   string `json:"new_text"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		responseutils.BadRequest(w, "Invalid request body")
		return
	}

	err := h.commentService.Update(r.Context(), uuid.MustParse(req.CommentID), req.NewText)
	if err != nil {
		switch {
		case errors.Is(err, servErrors.ErrCommentNotFound):
			responseutils.NotFound(w, err.Error())
			return
		case errors.Is(err, servErrors.ErrCommentTextEmpty):
			responseutils.BadRequest(w, err.Error())
			return
		case errors.Is(err, servErrors.ErrCommentForbiddenAction):
			responseutils.Forbidden(w, err.Error())
			return
		default:
			log.Println(err.Error())
			responseutils.InternalServerError(w, "Внутренняя ошибка сервера")
			return
		}
	}

	responseutils.JSON(w, http.StatusOK, "Текст поста успешно изменен.")
}

func (h *CommentHandler) GetAllPostComments(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		responseutils.StatusMethodNotAllowed(w, "Неправильный метод (должен быть GET)")
		return
	}

	postID := r.URL.Query().Get("post_id")

	posts, err := h.commentService.GetAllPostComments(r.Context(), uuid.MustParse(postID))
	if err != nil {
		switch {
		default:
			log.Println(err.Error())
			responseutils.InternalServerError(w, "Внутренняя ошибка сервера")
			return
		}
	}

	responseutils.JSON(w, http.StatusCreated, posts)
}

func (h *CommentHandler) DeleteCommentsByPost(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		responseutils.StatusMethodNotAllowed(w, "Неправильный метод (должен быть DELETE)")
		return
	}

	postID := r.URL.Query().Get("post_id")

	err := h.commentService.DeleteAllCommentsByPost(r.Context(), uuid.MustParse(postID))
	if err != nil {
		switch {
		case errors.Is(err, servErrors.ErrNoCommentsOnPost):
			responseutils.NotFound(w, err.Error())
			return
		default:
			log.Println(err.Error())
			responseutils.InternalServerError(w, "Внутренняя ошибка сервера")
			return
		}
	}

	responseutils.JSON(w, http.StatusOK, "Комментарии успешно удалены с поста.")
}
