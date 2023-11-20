package handlers

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"strconv"
	"strings"
	"time"

	"go.uber.org/zap"

	"go-gofermart-loyalty-system/internal/order"
	"go-gofermart-loyalty-system/internal/pkg/jwtauth"
)

type OrdersHandlers struct {
	log                         *zap.Logger
	service                     *order.OrderService
	asyncProcessingOrderService *order.AsyncProcessingOrderService
}

func NewOrdersHandlers(
	log *zap.Logger,
	service *order.OrderService,
	asyncProcessingOrderService *order.AsyncProcessingOrderService,
) *OrdersHandlers {
	return &OrdersHandlers{
		log:                         log,
		service:                     service,
		asyncProcessingOrderService: asyncProcessingOrderService,
	}
}

type ResponseOrderDTO struct {
	Number     string      `json:"number"`
	Status     string      `json:"status"`
	Accrual    json.Number `json:"accrual,omitempty"`
	UploadedAt time.Time   `json:"uploaded_at"`
}

func (r ResponseOrderDTO) MarshalJSON() ([]byte, error) {
	type ResponseOrderDTOAlias ResponseOrderDTO

	aliasValue := struct {
		ResponseOrderDTOAlias

		UploadedAt string `json:"uploaded_at"`
	}{
		ResponseOrderDTOAlias: ResponseOrderDTOAlias(r),

		UploadedAt: r.UploadedAt.Format(time.RFC3339),
	}

	return json.Marshal(aliasValue)
}

func (h *OrdersHandlers) GetAllByUser(wr http.ResponseWriter, r *http.Request) {
	jwtData, _ := jwtauth.JwtDataFromContext(r.Context()) // TODO: Авторизацию надо отрефачить

	ordersEntity, err := h.service.GetOrdersByUser(r.Context(), jwtData.ID)

	if err != nil {
		h.log.Error(
			"unknown error when get all order by user",
			zap.Error(err),
			zap.String("userId", jwtData.ID),
		)
		http.Error(wr, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)

		return
	}

	if len(ordersEntity) == 0 {
		wr.WriteHeader(http.StatusNoContent)

		return
	}

	orders := make([]ResponseOrderDTO, 0, len(ordersEntity))

	for _, orderEntity := range ordersEntity {
		var orderResponse = ResponseOrderDTO{
			Number:     orderEntity.Number,
			Status:     orderEntity.Status,
			UploadedAt: orderEntity.CreatedAt,
		}

		if orderEntity.Accrual.Valid {
			// TODO: Сделать из этого памятку
			orderResponse.Accrual = json.Number(strconv.FormatFloat(float64(orderEntity.Accrual.Int64)/100, 'f', 2, 64))
		}

		orders = append(orders, orderResponse)
	}

	responseByte, err := json.Marshal(orders)

	if err != nil {
		h.log.Error("error when ordersEntity json marshaling", zap.Error(err))
		http.Error(wr, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)

		return
	}

	wr.Header().Add("Content-Type", "application/json")
	wr.WriteHeader(http.StatusOK)
	_, _ = wr.Write(responseByte)
}

func (h *OrdersHandlers) AddOrder(wr http.ResponseWriter, r *http.Request) {
	jwtData, _ := jwtauth.JwtDataFromContext(r.Context()) // TODO: Авторизацию надо отрефачить

	bodyRaw, err := io.ReadAll(r.Body)
	defer r.Body.Close()

	if err != nil {
		h.log.Error("error when parsing body", zap.Error(err))
		http.Error(wr, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)

		return
	}

	orderNumber := strings.TrimSpace(string(bodyRaw))

	err = h.asyncProcessingOrderService.CreateOrderAndAddQueue(r.Context(), jwtData.ID, orderNumber)

	if errors.Is(err, order.ErrOrderNumberIsInvalid) {
		wr.WriteHeader(http.StatusUnprocessableEntity)

		return
	}

	if errors.Is(err, order.ErrOrderAlreadyExistAnotherUser) {
		wr.WriteHeader(http.StatusConflict)

		return
	}

	if errors.Is(err, order.ErrOrderIsExist) {
		wr.WriteHeader(http.StatusOK)

		return
	}

	if err != nil {
		h.log.Error("error when create and add order", zap.Error(err))
		http.Error(wr, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)

		return
	}

	wr.WriteHeader(http.StatusAccepted)
}
