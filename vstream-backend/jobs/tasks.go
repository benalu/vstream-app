package jobs

import (
	"fmt"
	"io"
	"log/slog"
	"os"
	"path/filepath"
	"time"
	"vstream-backend/database"
	"vstream-backend/models"
)

func CleanupOldLogs() {
	cutoff := time.Now().AddDate(0, 0, -30)
	result := database.DB.
		Where("reported_at < ?", cutoff).
		Delete(&models.PlaybackLog{})

	if result.Error != nil {
		slog.Error("Cleanup logs failed", "err", result.Error)
		return
	}

	slog.Info("Cleanup logs done", "deleted_rows", result.RowsAffected)
}

func BackupDatabase() {
	src := filepath.Join("data", "db", "vstream.db")
	backupDir := filepath.Join("data", "db", "backups")

	if err := os.MkdirAll(backupDir, 0755); err != nil {
		slog.Error("Backup: failed to create backup dir", "err", err)
		return
	}

	date := time.Now().Format("2006-01-02")
	dst := filepath.Join(backupDir, fmt.Sprintf("vstream_%s.db", date))

	if _, err := os.Stat(dst); err == nil {
		slog.Info("Backup already exists for today, skipping", "file", dst)
		return
	}

	// Checkpoint WAL dulu sebelum copy
	database.DB.Exec("PRAGMA wal_checkpoint(FULL)")

	if err := copyFile(src, dst); err != nil {
		slog.Error("Backup: copy failed", "src", src, "dst", dst, "err", err)
		return
	}

	slog.Info("Backup created", "file", dst)
	pruneOldBackups(backupDir, 4)
}

func copyFile(src, dst string) error {
	in, err := os.Open(src)
	if err != nil {
		return err
	}
	defer in.Close()

	out, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, in)
	return err
}

func pruneOldBackups(dir string, keep int) {
	entries, err := os.ReadDir(dir)
	if err != nil {
		slog.Error("Backup prune: failed to read dir", "err", err)
		return
	}

	type fileInfo struct {
		path    string
		modTime time.Time
	}

	var files []fileInfo
	for _, e := range entries {
		if e.IsDir() {
			continue
		}
		info, err := e.Info()
		if err != nil {
			continue
		}
		files = append(files, fileInfo{
			path:    filepath.Join(dir, e.Name()),
			modTime: info.ModTime(),
		})
	}

	if len(files) <= keep {
		return
	}

	for i := 0; i < len(files)-1; i++ {
		for j := i + 1; j < len(files); j++ {
			if files[j].modTime.Before(files[i].modTime) {
				files[i], files[j] = files[j], files[i]
			}
		}
	}

	toDelete := files[:len(files)-keep]
	for _, f := range toDelete {
		if err := os.Remove(f.path); err != nil {
			slog.Error("Backup prune: failed to delete", "file", f.path, "err", err)
		} else {
			slog.Info("Backup pruned", "file", f.path)
		}
	}
}
