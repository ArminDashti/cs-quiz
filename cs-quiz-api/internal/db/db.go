package db

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"
)

// Open connects to PostgreSQL.
func Open(databaseURL string) (*sql.DB, error) {
	if databaseURL == "" {
		return nil, fmt.Errorf("DATABASE_URL is required")
	}
	db, err := sql.Open("pgx", databaseURL)
	if err != nil {
		return nil, err
	}
	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(5)
	db.SetConnMaxLifetime(30 * time.Minute)

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()
	if err := db.PingContext(ctx); err != nil {
		_ = db.Close()
		return nil, fmt.Errorf("ping database: %w", err)
	}
	return db, nil
}

// Migrate runs SQL files in migrationsDir in lexical order.
func Migrate(db *sql.DB, migrationsDir string) error {
	entries, err := os.ReadDir(migrationsDir)
	if err != nil {
		return fmt.Errorf("read migrations: %w", err)
	}
	var files []string
	for _, e := range entries {
		if e.IsDir() || !strings.HasSuffix(e.Name(), ".sql") {
			continue
		}
		files = append(files, e.Name())
	}
	sort.Strings(files)

	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	for _, name := range files {
		body, err := os.ReadFile(filepath.Join(migrationsDir, name))
		if err != nil {
			return err
		}
		if _, err := db.ExecContext(ctx, string(body)); err != nil {
			return fmt.Errorf("migration %s: %w", name, err)
		}
	}
	return nil
}

// SeedInviteCode inserts the invite code when missing.
func SeedInviteCode(ctx context.Context, db *sql.DB, code string) error {
	_, err := db.ExecContext(ctx, `
		INSERT INTO settings (key, value) VALUES ('invite_code', $1)
		ON CONFLICT (key) DO NOTHING
	`, code)
	return err
}

// DefaultUser holds the local/demo login account seeded on startup.
const (
	DefaultUsername = "armin"
	DefaultPassword = "dopadopa123"
	DefaultEmail    = "armin@local"
	AdminUsername   = "admin"
	AdminEmail      = "admin@local"
)

// SeedDefaultUser upserts the default local login accounts (armin + admin).
// passwordHash must already be a bcrypt hash of DefaultPassword.
func SeedDefaultUser(ctx context.Context, db *sql.DB, passwordHash string) error {
	if err := upsertDemoUser(ctx, db, DefaultEmail, passwordHash, "Armin", "Dashti", DefaultUsername); err != nil {
		return err
	}
	return upsertDemoUser(ctx, db, AdminEmail, passwordHash, "Admin", "User", AdminUsername)
}

func upsertDemoUser(ctx context.Context, db *sql.DB, email, passwordHash, first, last, username string) error {
	_, err := db.ExecContext(ctx, `
		INSERT INTO users (email, password_hash, first_name, last_name, username, is_admin)
		VALUES ($1, $2, $3, $4, $5, TRUE)
		ON CONFLICT (email) DO UPDATE SET
			password_hash = EXCLUDED.password_hash,
			username = EXCLUDED.username,
			first_name = EXCLUDED.first_name,
			last_name = EXCLUDED.last_name,
			is_admin = TRUE,
			updated_at = NOW()
	`, email, passwordHash, first, last, username)
	return err
}

// GetSetting returns a settings value.
func GetSetting(ctx context.Context, db *sql.DB, key string) (string, error) {
	var value string
	err := db.QueryRowContext(ctx, `SELECT value FROM settings WHERE key = $1`, key).Scan(&value)
	return value, err
}

// SetSetting upserts a settings value.
func SetSetting(ctx context.Context, db *sql.DB, key, value string) error {
	_, err := db.ExecContext(ctx, `
		INSERT INTO settings (key, value) VALUES ($1, $2)
		ON CONFLICT (key) DO UPDATE SET value = EXCLUDED.value
	`, key, value)
	return err
}
