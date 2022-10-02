package order

import (
	"context"
	client_loyalty_points "go-gofermart-loyalty-system/internal/pkg/client-loyalty-points"
	"go.uber.org/zap"
	"time"
)

type OrderJob struct {
	orderNumber string
}

type Queue struct {
	ch chan *OrderJob
}

func NewQueue() *Queue {
	q := Queue{
		ch: make(chan *OrderJob, 100),
	}

	return &q
}

func (q *Queue) Add(orderJob *OrderJob) {
	q.ch <- orderJob
}

func (q *Queue) Stop() {
	close(q.ch)
}

type Worker struct {
	ID                  int
	q                   *Queue
	log                 *zap.Logger
	orderService        *OrderService
	clientLoyaltyPoints *client_loyalty_points.ClientLoyaltyPoints
}

func (w *Worker) Loop() {
	for job := range w.q.ch {
		ctx := context.Background()

		w.log.Info("Start processing order", zap.String("orderNumber", job.orderNumber))

		w.log.Info("update status on processing", zap.String("orderNumber", job.orderNumber))
		if err := w.orderService.SetProcessingStatusByNumber(ctx, job.orderNumber); err != nil {
			w.log.Error(
				"error change processing status for order",
				zap.Error(err),
				zap.String("orderNumber", job.orderNumber),
			)

			continue
		}

		w.log.Info("request to internal system", zap.String("orderNumber", job.orderNumber))
		response, err := w.clientLoyaltyPoints.GetOrder(ctx, job.orderNumber)

		// TODO: Сюда нужны разные виды ошибок
		if err != nil {
			// TODO: Добавить смену статуса на ошибочную
			w.log.Error("error fetch info order", zap.Error(err), zap.String("orderNumber", job.orderNumber))
			continue
		}

		if response.Status == client_loyalty_points.ClientResponseOrderStatusRegistered ||
			response.Status == client_loyalty_points.ClientResponseOrderStatusProcessing {
			w.log.Info("order processing in external system", zap.String("orderNumber", job.orderNumber))

			time.Sleep(500 * time.Millisecond)

			w.q.Add(job)

			return
		}

		if response.Status == client_loyalty_points.ClientResponseOrderStatusInvalid {
			w.log.Info("order invalid processing in external system", zap.String("orderNumber", job.orderNumber))
			err := w.orderService.SetInvalidStatusByNumber(ctx, job.orderNumber)

			if err != nil {
				w.log.Error("can't update status on failed", zap.Error(err))
			}

			return
		}

		if response.Status == client_loyalty_points.ClientResponseOrderStatusProcessed {
			w.log.Info("order processed in external system", zap.String("orderNumber", job.orderNumber))
			accuralInt64, err := response.Accrual.Int64()

			if err != nil {
				w.log.Error(
					"can't convert Accrual to int",
					zap.String("orderNumber", job.orderNumber),
					zap.Error(err),
				)

				return
			}

			err = w.orderService.SetProcessedStatusByNumber(
				ctx,
				job.orderNumber,
				int(accuralInt64),
			)

			if err != nil {
				w.log.Error("can't update status on processed", zap.Error(err))

				return
			}

			return
		}
	}
}

type WorkerPool struct {
	q *Queue
}

// TODO: Реализовать обработку ошибок в worker и перезапускать
func NewWorkerPool(log *zap.Logger, orderService *OrderService, clientLoyaltyPoints *client_loyalty_points.ClientLoyaltyPoints, nWorker int) *WorkerPool {
	q := NewQueue()

	for i := 0; i < nWorker; i++ {
		worker := &Worker{
			ID:                  i,
			q:                   q,
			log:                 log.With(zap.Int("workerID", i)),
			orderService:        orderService,
			clientLoyaltyPoints: clientLoyaltyPoints,
		}

		go worker.Loop()
	}

	return &WorkerPool{
		q: q,
	}
}

func (w *WorkerPool) AddJob(job *OrderJob) {
	w.q.Add(job)
}

func (w *WorkerPool) Stop() {
	w.q.Stop()
}
