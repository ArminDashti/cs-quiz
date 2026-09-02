package handlers

import (
	"database/sql"
	"net/http"
	"strings"

	"github.com/ArminDashti/cs-quiz-api/internal/auth"
	appdb "github.com/ArminDashti/cs-quiz-api/internal/db"
	"github.com/gin-gonic/gin"
)

type registerRequest struct {
	FirstName  string `json:"first_name" binding:"required"`
	LastName   string `json:"last_name" binding:"required"`
	Email      string `json:"email" binding:"required,email"`
	Password   string `json:"password" binding:"required,min=8"`
	InviteCode string `json:"invite_code" binding:"required"`
}

type loginRequest struct {
	Username string `json:"username"`
	Email    string `json:"email"` // accepted for backward compatibility
	Password string `json:"password" binding:"required"`
}

// Register creates a new user account.
func (h *Handler) Register(c *gin.Context) {
	var req registerRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		writeError(c, http.StatusBadRequest, "invalid request")
		return
	}

	ctx := c.Request.Context()
	invite, err := appdb.GetSetting(ctx, h.db, "invite_code")
	if err != nil {
		writeError(c, http.StatusInternalServerError, "invite code unavailable")
		return
	}
	if strings.TrimSpace(req.InviteCode) != invite {
		writeError(c, http.StatusForbidden, "invalid invite code")
		return
	}

	email := normalizeEmail(req.Email)
	hash, err := auth.HashPassword(req.Password)
	if err != nil {
		writeError(c, http.StatusInternalServerError, "could not hash password")
		return
	}

	isAdmin := isAdminEmail(h.cfg.AdminEmail, email)
	var id string
	err = h.db.QueryRowContext(ctx, `
		INSERT INTO users (email, password_hash, first_name, last_name, is_admin)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id
	`, email, hash, strings.TrimSpace(req.FirstName), strings.TrimSpace(req.LastName), isAdmin).Scan(&id)
	if err != nil {
		if fmtUniqueViolation(err) {
			writeError(c, http.StatusConflict, "email already registered")
			return
		}
		writeError(c, http.StatusInternalServerError, "could not create user")
		return
	}

	user, err := h.getUserByID(c, id)
	if err != nil {
		writeError(c, http.StatusInternalServerError, "could not load user")
		return
	}

	token, err := auth.IssueToken(h.cfg.JWTSecret, user.ID, user.Email, user.IsAdmin)
	if err != nil {
		writeError(c, http.StatusInternalServerError, "could not issue token")
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"token": token,
		"user":  h.userPublic(user),
	})
}

// Login authenticates a user by username (preferred) or email.
func (h *Handler) Login(c *gin.Context) {
	var req loginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		writeError(c, http.StatusBadRequest, "invalid request")
		return
	}

	loginID := strings.TrimSpace(req.Username)
	if loginID == "" {
		loginID = strings.TrimSpace(req.Email)
	}
	if loginID == "" {
		writeError(c, http.StatusBadRequest, "username is required")
		return
	}

	user, err := h.getUserByLogin(c, loginID)
	if err == sql.ErrNoRows {
		writeError(c, http.StatusUnauthorized, "invalid username or password")
		return
	}
	if err != nil {
		writeError(c, http.StatusInternalServerError, "could not load user")
		return
	}
	if !auth.CheckPassword(user.PasswordHash, req.Password) {
		writeError(c, http.StatusUnauthorized, "invalid username or password")
		return
	}

	if isAdminEmail(h.cfg.AdminEmail, user.Email) && !user.IsAdmin {
		_, _ = h.db.ExecContext(c.Request.Context(), `UPDATE users SET is_admin = TRUE, updated_at = NOW() WHERE id = $1`, user.ID)
		user.IsAdmin = true
	}

	token, err := auth.IssueToken(h.cfg.JWTSecret, user.ID, user.Email, user.IsAdmin)
	if err != nil {
		writeError(c, http.StatusInternalServerError, "could not issue token")
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"token": token,
		"user":  h.userPublic(user),
	})
}

// Me returns the current user.
func (h *Handler) Me(c *gin.Context) {
	user, err := h.getUserByID(c, auth.UserIDFromContext(c))
	if err == sql.ErrNoRows {
		writeError(c, http.StatusUnauthorized, "user not found")
		return
	}
	if err != nil {
		writeError(c, http.StatusInternalServerError, "could not load user")
		return
	}
	c.JSON(http.StatusOK, h.userPublic(user))
}
