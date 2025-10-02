package server

import (
	"context"

	"github.com/meehighlov/grats/internal/clients/clients/telegram"
)

func (s *Server) worker(workerID int, updates *telegram.UpdatesChannel) {
	defer s.wgWorkerPool.Done()

	s.logger.Debug("Worker started", "worker_id", workerID)

	for {
		select {
		case update, ok := <-*updates:
			if !ok {
				s.logger.Debug("Worker stopped - channel closed", "worker_id", workerID)
				return
			}

			ctx, cancel := context.WithTimeout(s.workerCtx, s.handleTimeout)

			if err := s.HandleUpdate(ctx, &update); err != nil {
				s.logger.Error("Worker failed to handle update",
					"worker_id", workerID,
					"update_id", update.UpdateId,
					"error", err)
			}

			cancel()

		case <-s.workerCtx.Done():
			s.logger.Debug("Worker stopped - context cancelled", "worker_id", workerID)
			return
		}
	}
}
