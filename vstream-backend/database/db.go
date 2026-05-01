package database

import (
	"log/slog"
	"os"
	"path/filepath"
	"time"

	"vstream-backend/models"

	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
)

var DB *gorm.DB

func InitDB() {
	dir := filepath.Join("data", "db")
	dbPath := filepath.Join(dir, "vstream.db")

	if _, err := os.Stat(dir); os.IsNotExist(err) {
		_ = os.MkdirAll(dir, 0755)
	}

	db, err := gorm.Open(sqlite.Open(dbPath), &gorm.Config{PrepareStmt: true})
	if err != nil {
		slog.Error("DB connection failed", "path", dbPath, "err", err)
		os.Exit(1)
	}

	sqlDB, _ := db.DB()
	sqlDB.SetMaxOpenConns(1)
	sqlDB.SetMaxIdleConns(1)
	sqlDB.SetConnMaxLifetime(time.Hour)

	db.Exec("PRAGMA journal_mode=WAL;")
	db.Exec("PRAGMA synchronous=NORMAL;")

	if err := db.AutoMigrate(&models.Movie{}, &models.WatchSession{}, &models.Season{}, &models.Episode{}, &models.Subtitle{}, &models.HeroSlide{}, &models.PlaybackLog{}); err != nil {
		slog.Error("Migration failed", "err", err)
		os.Exit(1)
	}

	DB = db
	slog.Info("Database initialized", "mode", "WAL", "path", dbPath)
}
