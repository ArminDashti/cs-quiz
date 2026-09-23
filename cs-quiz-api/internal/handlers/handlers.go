package handlers

import (
	"database/sql"
	"fmt"
	"net/http"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/ArminDashti/cs-quiz-api/internal/config"
	"github.com/ArminDashti/cs-quiz-api/internal/models"
	"github.com/gin-gonic/gin"
)

// Handler exposes HTTP endpoints.
type Handler struct {
	db  *sql.DB
	cfg config.Config
}

// New creates a Handler.
func New(db *sql.DB, cfg config.Config) *Handler {
	return &Handler{db: db, cfg: cfg}
}

// Health returns a simple liveness payload.
func (h *Handler) Health(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}

// RecordVisit stores basic request metadata for usage analytics. It intentionally
// records no credentials, request bodies, or authentication tokens.
func (h *Handler) RecordVisit() gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.Request.Method != http.MethodOptions {
			ip := c.ClientIP()
			os, browser := visitClient(c.GetHeader("User-Agent"))
			_, _ = h.db.ExecContext(c.Request.Context(), `
				INSERT INTO visit_logs (ip_address, operating_system, browser, path)
				VALUES ($1, $2, $3, $4)
			`, ip, os, browser, c.Request.URL.Path)
		}
		c.Next()
	}
}

var browserVersion = regexp.MustCompile(`(?:edg|opr|chrome|firefox|version)/([0-9.]+)`)

func visitClient(userAgent string) (string, string) {
	ua := strings.ToLower(userAgent)
	os := "Unknown"
	switch {
	case strings.Contains(ua, "windows"):
		os = "Windows"
	case strings.Contains(ua, "android"):
		os = "Android"
	case strings.Contains(ua, "iphone") || strings.Contains(ua, "ipad"):
		os = "iOS"
	case strings.Contains(ua, "mac os") || strings.Contains(ua, "macintosh"):
		os = "macOS"
	case strings.Contains(ua, "linux"):
		os = "Linux"
	}

	browser := "Unknown"
	switch {
	case strings.Contains(ua, "edg/"):
		browser = "Edge"
	case strings.Contains(ua, "opr/"):
		browser = "Opera"
	case strings.Contains(ua, "firefox/"):
		browser = "Firefox"
	case strings.Contains(ua, "chrome/"):
		browser = "Chrome"
	case strings.Contains(ua, "safari/"):
		browser = "Safari"
	}
	if match := browserVersion.FindStringSubmatch(ua); len(match) == 2 && browser != "Unknown" {
		browser += " " + match[1]
	}
	return os, browser
}

func (h *Handler) avatarURL(avatarPath *string) *string {
	if avatarPath == nil || *avatarPath == "" {
		return nil
	}
	url := "/uploads/" + filepath.ToSlash(*avatarPath)
	return &url
}

func (h *Handler) userPublic(u models.User) models.User {
	u.AvatarURL = h.avatarURL(u.AvatarPath)
	u.AvatarPath = nil
	return u
}

func scanUser(row interface {
	Scan(dest ...any) error
}) (models.User, error) {
	var u models.User
	var username, avatar sql.NullString
	err := row.Scan(
		&u.ID, &u.Email, &u.PasswordHash, &u.FirstName, &u.LastName,
		&username, &avatar, &u.IsAdmin, &u.CreatedAt, &u.UpdatedAt,
	)
	if err != nil {
		return u, err
	}
	if username.Valid {
		u.Username = &username.String
	}
	if avatar.Valid {
		u.AvatarPath = &avatar.String
	}
	return u, nil
}

const userSelect = `
	SELECT id, email, password_hash, first_name, last_name, username, avatar_path, is_admin, created_at, updated_at
	FROM users
`

func (h *Handler) getUserByID(c *gin.Context, id string) (models.User, error) {
	return scanUser(h.db.QueryRowContext(c.Request.Context(), userSelect+` WHERE id = $1`, id))
}

func (h *Handler) getUserByEmail(c *gin.Context, email string) (models.User, error) {
	return scanUser(h.db.QueryRowContext(c.Request.Context(), userSelect+` WHERE email = $1`, strings.ToLower(email)))
}

func (h *Handler) getUserByLogin(c *gin.Context, login string) (models.User, error) {
	login = strings.TrimSpace(login)
	if strings.Contains(login, "@") {
		return h.getUserByEmail(c, login)
	}
	return scanUser(h.db.QueryRowContext(
		c.Request.Context(),
		userSelect+` WHERE LOWER(username) = LOWER($1)`,
		login,
	))
}

func normalizeEmail(email string) string {
	return strings.ToLower(strings.TrimSpace(email))
}

func isAdminEmail(cfgAdmin, email string) bool {
	if cfgAdmin == "" {
		return false
	}
	return normalizeEmail(cfgAdmin) == normalizeEmail(email)
}

func writeError(c *gin.Context, status int, msg string) {
	c.JSON(status, gin.H{"error": msg})
}

func nullString(p *string) sql.NullString {
	if p == nil || *p == "" {
		return sql.NullString{}
	}
	return sql.NullString{String: *p, Valid: true}
}

func fmtUniqueViolation(err error) bool {
	return err != nil && strings.Contains(err.Error(), "duplicate key")
}

func absUpload(cfg config.Config, rel string) string {
	return filepath.Join(cfg.UploadDir, rel)
}

func safeExt(filename string) string {
	ext := strings.ToLower(filepath.Ext(filename))
	switch ext {
	case ".jpg", ".jpeg", ".png", ".webp", ".gif":
		return ext
	default:
		return ".jpg"
	}
}

func avatarRelPath(userID, ext string) string {
	return filepath.ToSlash(filepath.Join("avatars", fmt.Sprintf("%s%s", userID, ext)))
}

func slugify(name string) string {
	s := strings.ToLower(strings.TrimSpace(name))
	s = strings.ReplaceAll(s, " ", "-")
	var b strings.Builder
	for _, r := range s {
		if (r >= 'a' && r <= 'z') || (r >= '0' && r <= '9') || r == '-' {
			b.WriteRune(r)
		}
	}
	out := b.String()
	for strings.Contains(out, "--") {
		out = strings.ReplaceAll(out, "--", "-")
	}
	return strings.Trim(out, "-")
}
