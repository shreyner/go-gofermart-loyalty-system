package handlers

import (
	"go.uber.org/zap"
	"net/http"
)

type UsersHandler struct {
	log *zap.Logger
}

func NewUsersHandlers(logger *zap.Logger) *UsersHandler {
	log := logger.Named("UsersHandlers")

	return &UsersHandler{
		log,
	}
}

func (h *UsersHandler) GetMe(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(http.StatusOK)

	_, err := w.Write([]byte("Hello world"))

	if err != nil {
		h.log.Error("Get error:", zap.Error(err))

		return
	}
}
