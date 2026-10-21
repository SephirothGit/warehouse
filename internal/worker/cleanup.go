package worker

import (
	"context"
	"log/slog"
	"time"

	"github.com/SephirothGit/warehouse/internal/repository"
)

func RunRefreshTokenCleanup(ctx context.Context, repo repository.RefreshTokenRepo, interval time.Duration) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			opCtx, cancel := context.WithTimeout(ctx, 30 * time.Second)
			n, err := repo.DeleteExpired(opCtx, time.Now())
			cancel()
			if err != nil {
				slog.Error("refresh token cleanup failed", "error", err)
				continue
			}
			slog.Info("expired refresh tokens deleted", "count", n)
		}
	}
}