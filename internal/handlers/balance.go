package handlers

import (
	"encoding/json"
	"errors"
	"net/http"

	"go.uber.org/zap"

	"go-gofermart-loyalty-system/internal/balance"
	"go-gofermart-loyalty-system/internal/pkg/jwtauth"
)

type BalanceHandlers struct {
	log     *zap.Logger
	service *balance.BalanceService
}

func NewBalanceHandlers(log *zap.Logger, service *balance.BalanceService) *BalanceHandlers {
	return &BalanceHandlers{
		log:     log,
		service: service,
	}
}

type ResponseBalanceDTO struct {
	Current   int `json:"current"`
	Withdrawn int `json:"withdrawn"`
}

func (h *BalanceHandlers) GetUserBalance(wr http.ResponseWriter, r *http.Request) {
	jwtData, _ := jwtauth.JwtDataFromContext(r.Context())

	b, err := h.service.GetByUserID(r.Context(), jwtData.ID)

	if errors.Is(err, balance.ErrBalanceNotFound) {
		h.log.Error("can't find balance for user", zap.Error(err), zap.String("userId", jwtData.ID))
		http.Error(wr, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	if err != nil {
		h.log.Error("unknown error when find balance by userId", zap.Error(err))
		http.Error(wr, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	responseBalanceDTO := ResponseBalanceDTO{
		Current:   b.Current,
		Withdrawn: b.Withdrawn,
	}

	balanceResponse, err := json.Marshal(responseBalanceDTO)

	if err != nil {
		h.log.Error("error when balance json marshaling", zap.Error(err))
		http.Error(wr, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)

		return
	}

	wr.Header().Add("Content-Type", "application/json")
	wr.WriteHeader(http.StatusOK)
	_, _ = wr.Write(balanceResponse)
}
