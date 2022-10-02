package order

import (
	"fmt"
	client_loyalty_points "go-gofermart-loyalty-system/internal/pkg/client-loyalty-points"
	"go.uber.org/zap"
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
	orderService        *orderService
	clientLoyaltyPoints *client_loyalty_points.ClientLoyaltyPoints
}

func (w *Worker) Loop() {
	for job := range w.q.ch {
		if err := w.orderService.SetProcessStatusById(job.orderNumber); err != nil {
			w.log.Error("error change processing status for order", zap.Error(err), zap.String("orderNumber", job.orderNumber))
			continue
		}

		if err := w.clientLoyaltyPoints.GetOrder(job.orderNumber); err != nil {
			// TODO: Добавить смену статуса на ошибочную
			w.log.Error("error fetch info order", zap.Error(err), zap.String("orderNumber", job.orderNumber))
			continue
		}

		// TODO: Добавить проставление статуса "PROCESSED" и заполнение инфы

		fmt.Println(job)
	}
}

type WorkerPool struct {
	q *Queue
}

// TODO: Реализовать обработку ошибок в worker и перезапускать
func NewWorkerPool(log *zap.Logger, orderService *orderService, clientLoyaltyPoints *client_loyalty_points.ClientLoyaltyPoints, nWorker int) *WorkerPool {
	q := NewQueue()

	for i := 0; i < nWorker; i++ {
		worker := &Worker{
			ID:           i,
			q:            q,
			log:          log,
			orderService: orderService,
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
