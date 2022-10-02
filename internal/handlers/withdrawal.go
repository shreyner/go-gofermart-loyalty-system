package handlers

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"time"

	"go.uber.org/zap"

	"go-gofermart-loyalty-system/internal/balance"
	"go-gofermart-loyalty-system/internal/pkg/jwtauth"
	"go-gofermart-loyalty-system/internal/withdrawal"
)

type WithdrawalHandlers struct {
	log     *zap.Logger
	service *withdrawal.WithdrawalService
}

func NewWithdrawalHandlers(log *zap.Logger, service *withdrawal.WithdrawalService) *WithdrawalHandlers {
	return &WithdrawalHandlers{
		log:     log,
		service: service,
	}
}

type RequestCreateWithdrawalDTO struct {
	Order string `json:"order"`
	Sum   int    `json:"sum"`
}

func (w *WithdrawalHandlers) CreateWithdrawal(wr http.ResponseWriter, r *http.Request) {
	jwtData, _ := jwtauth.JwtDataFromContext(r.Context()) // TODO: Авторизацию надо отрефачить

	bytes, err := io.ReadAll(r.Body)
	defer r.Body.Close()

	if err != nil {
		// TODO: Похож на шаблонный код. Можно вынести
		w.log.Error("can't read all bytes from body", zap.Error(err))
		http.Error(wr, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)

		return
	}

	var requestCreateWithdrawalDTO *RequestCreateWithdrawalDTO

	if err := json.Unmarshal(bytes, &requestCreateWithdrawalDTO); err != nil {
		// TODO: Подумать и может обработать ошибку не валидного JSON как BadRequest
		w.log.Error("can't read all bytes from body", zap.Error(err))
		http.Error(wr, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)

		return
	}

	err = w.service.Create(r.Context(), jwtData.ID, requestCreateWithdrawalDTO.Order, requestCreateWithdrawalDTO.Sum)

	if err != nil {
		if errors.Is(err, withdrawal.ErrWithdrawalOrderNumberIsInvalid) {
			w.log.Warn("invalid order number", zap.String("order", requestCreateWithdrawalDTO.Order))
			http.Error(wr, http.StatusText(http.StatusUnprocessableEntity), http.StatusUnprocessableEntity)

			return
		}

		if errors.Is(err, withdrawal.ErrWithdrawalOrderIsExist) {
			w.log.Warn("order is already exists for withdrawal", zap.String("order", requestCreateWithdrawalDTO.Order))
			http.Error(wr, http.StatusText(http.StatusConflict), http.StatusConflict)

			return
		}

		if errors.Is(err, balance.ErrBalanceCannotBeNegative) {
			http.Error(wr, http.StatusText(http.StatusPaymentRequired), http.StatusPaymentRequired)
			return
		}

		w.log.Error("unknown error when create withdrawal", zap.Error(err))
		http.Error(wr, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)

		return
	}

	wr.WriteHeader(http.StatusOK)
}

type ResponseWithdrawalsDTO struct {
	Order       string    `json:"order"`
	Sum         int       `json:"sum"`
	ProcessedAt time.Time `json:"processed_at"`
}

func (w ResponseWithdrawalsDTO) MarshalJSON() ([]byte, error) {
	type ResponseWithdrawalsDTOAlias ResponseWithdrawalsDTO

	aliasValue := struct {
		ResponseWithdrawalsDTOAlias

		ProcessedAt string `json:"processed_at"`
	}{
		ResponseWithdrawalsDTOAlias: ResponseWithdrawalsDTOAlias(w),

		ProcessedAt: w.ProcessedAt.Format(time.RFC3339),
	}

	return json.Marshal(aliasValue)
}

func (w *WithdrawalHandlers) GetAllByUser(wr http.ResponseWriter, r *http.Request) {
	jwtData, _ := jwtauth.JwtDataFromContext(r.Context()) // TODO: Авторизацию надо отрефачить

	withdrawals, err := w.service.GetAllByUser(r.Context(), jwtData.ID)

	if err != nil {
		w.log.Error(
			"unknown error when get all withdrawals by user",
			zap.Error(err),
			zap.String("userId", jwtData.ID),
		)
		http.Error(wr, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)

		return
	}

	if len(withdrawals) == 0 {
		wr.WriteHeader(http.StatusNoContent)

		return
	}

	withdrawalsResponse := make([]ResponseWithdrawalsDTO, 0, len(withdrawals))

	for _, withdrawalEntity := range withdrawals {
		responseWithdrawal := ResponseWithdrawalsDTO{
			Order:       withdrawalEntity.OrderNumber,
			Sum:         withdrawalEntity.Sum,
			ProcessedAt: withdrawalEntity.CreatedAt,
		}

		withdrawalsResponse = append(withdrawalsResponse, responseWithdrawal)
	}

	responseBytes, err := json.Marshal(withdrawalsResponse)

	if err != nil {
		w.log.Error("error when withdrawalsEntity json marshaling", zap.Error(err))
		http.Error(wr, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)

		return
	}

	// TODO: Код повторяется, можно вынести в функцию
	wr.Header().Add("Content-Type", "application/json")
	wr.WriteHeader(http.StatusOK)
	_, _ = wr.Write(responseBytes)
}
