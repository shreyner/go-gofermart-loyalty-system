package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"go-gofermart-loyalty-system/internal/pkg/jwtauth"
	"go.uber.org/zap"
	"io"
	"net/http"

	"go-gofermart-loyalty-system/internal/auth"
	"go-gofermart-loyalty-system/internal/user"
)

type AuthHandlers struct {
	log     *zap.Logger
	service *auth.AuthService
}

func NewAuthHandlers(log *zap.Logger, service *auth.AuthService) *AuthHandlers {
	return &AuthHandlers{
		log:     log,
		service: service,
	}
}

type RegistryUserDTO struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

func (h *AuthHandlers) Registry(wr http.ResponseWriter, r *http.Request) {
	bytes, err := io.ReadAll(r.Body)
	defer r.Body.Close()

	if err != nil {
		h.log.Error("can't read all bytes from body", zap.Error(err))
		http.Error(wr, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)

		return
	}

	var registryUserDTO *RegistryUserDTO

	if err := json.Unmarshal(bytes, &registryUserDTO); err != nil {
		h.log.Error("can't read all bytes from body", zap.Error(err))
		http.Error(wr, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)

		return
	}

	u, err := h.service.RegistryByLogin(r.Context(), registryUserDTO.Login, registryUserDTO.Password)

	if errors.Is(err, user.ErrLoginAlreadyExist) {
		http.Error(wr, http.StatusText(http.StatusConflict), http.StatusConflict)

		return
	}

	if err != nil {
		h.log.Error("auth registry by login error", zap.Error(err))
		http.Error(wr, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)

		return
	}

	tokenString, err := jwtauth.CreateJwtToken(&jwtauth.JwtData{ID: u.ID})

	if err != nil {
		h.log.Error("can't sign jwt token", zap.Error(err))

		http.Error(wr, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)

		return
	}

	// TODO: Need refactoring and move to pkg/jwtauth
	authCookie := http.Cookie{
		Name:  "auth",
		Value: tokenString,
	}

	http.SetCookie(wr, &authCookie)

	wr.WriteHeader(http.StatusOK)
}

type LoginUserDTO struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

func (h *AuthHandlers) Login(wr http.ResponseWriter, r *http.Request) {
	bytes, err := io.ReadAll(r.Body)
	defer r.Body.Close()

	if err != nil {
		h.log.Error("can't read all bytes from body", zap.Error(err))
		http.Error(wr, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)

		return
	}

	var loginUserDTO *LoginUserDTO

	if err := json.Unmarshal(bytes, &loginUserDTO); err != nil {
		h.log.Error("can't read all bytes from body", zap.Error(err))
		http.Error(wr, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)

		return
	}

	fmt.Println(loginUserDTO)

	u, err := h.service.Login(r.Context(), loginUserDTO.Login, loginUserDTO.Password)

	if errors.Is(err, user.ErrUserNotFound) || errors.Is(err, user.ErrUserPasswordIncorrect) {
		http.Error(wr, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)

		return
	}

	if err != nil {
		h.log.Error("unknown error when find and verify password", zap.Error(err))
		http.Error(wr, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)

		return
	}

	tokenString, err := jwtauth.CreateJwtToken(&jwtauth.JwtData{ID: u.ID})

	// TODO: Need refactoring and move to pkg/jwtauth
	authCookie := http.Cookie{
		Name:  "auth",
		Value: tokenString,
	}

	http.SetCookie(wr, &authCookie)

	wr.WriteHeader(http.StatusOK)
}
