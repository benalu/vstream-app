package config

import (
	"fmt"
	"log/slog"
	"os"
	"strings"
)

// Required adalah daftar env yang wajib ada dan tidak boleh kosong
var required = []string{
	"ADMIN_KEY",
	"JWT_SECRET",
	"TMDB_API_KEY",
	"ALLOWED_ORIGIN",
}

// Validate memeriksa semua env yang wajib ada.
// Kalau ada yang kosong, program exit dengan pesan jelas.
func Validate() {
	var missing []string
	for _, key := range required {
		if strings.TrimSpace(os.Getenv(key)) == "" {
			missing = append(missing, key)
		}
	}
	if len(missing) > 0 {
		slog.Error("Startup aborted: missing required environment variables",
			"missing", strings.Join(missing, ", "),
		)
		fmt.Fprintf(os.Stderr, "\n[FATAL] Environment variables berikut wajib diisi:\n")
		for _, k := range missing {
			fmt.Fprintf(os.Stderr, "  - %s\n", k)
		}
		fmt.Fprintln(os.Stderr)
		os.Exit(1)
	}

	slog.Info("Environment validated", "vars_checked", len(required))
}
