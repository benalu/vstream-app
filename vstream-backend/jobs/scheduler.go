package jobs

import (
	"log/slog"
	"time"
)

func Start() {
	go runEvery(24*time.Hour, "cleanup-logs", CleanupOldLogs)
	go runEvery(7*24*time.Hour, "backup-db", BackupDatabase)
}

func runEvery(interval time.Duration, name string, fn func()) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	slog.Info("Job scheduled", "name", name, "interval", interval.String())

	for range ticker.C {
		slog.Info("Job running", "name", name)
		fn()
	}
}
