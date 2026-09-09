package handlers

import (
	"database/sql"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

// ExternalIngest accepts quiz (+ optional questions) from an external system via POST.
// Auth: X-API-Key or Authorization: Bearer <EXTERNAL_API_KEY>.
func (h *Handler) ExternalIngest(c *gin.Context) {
	if strings.TrimSpace(h.cfg.ExternalAPIKey) == "" {
		writeError(c, http.StatusServiceUnavailable, "external ingest not configured")
		return
	}
	key := strings.TrimSpace(c.GetHeader("X-API-Key"))
	if key == "" {
		authz := strings.TrimSpace(c.GetHeader("Authorization"))
		if strings.HasPrefix(strings.ToLower(authz), "bearer ") {
			key = strings.TrimSpace(authz[7:])
		}
	}
	if key == "" || key != h.cfg.ExternalAPIKey {
		writeError(c, http.StatusUnauthorized, "invalid api key")
		return
	}

	var req struct {
		Name        string   `json:"name"`
		Title       string   `json:"title"`
		Slug        string   `json:"slug"`
		Category    string   `json:"category"`
		Description string   `json:"description"`
		Attachments []string `json:"attachments"`
		Enabled     *bool    `json:"enabled"`
		Questions   []struct {
			Prompt       string `json:"prompt" binding:"required"`
			OptionA      string `json:"option_a" binding:"required"`
			OptionB      string `json:"option_b" binding:"required"`
			OptionC      string `json:"option_c" binding:"required"`
			OptionD      string `json:"option_d" binding:"required"`
			CorrectIndex int    `json:"correct_index"`
			QuestionType string `json:"question_type"`
			Difficulty   string `json:"difficulty"`
			SortOrder    int    `json:"sort_order"`
		} `json:"questions"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		writeError(c, http.StatusBadRequest, "invalid request")
		return
	}

	name := strings.TrimSpace(req.Title)
	if name == "" {
		name = strings.TrimSpace(req.Name)
	}
	if name == "" {
		writeError(c, http.StatusBadRequest, "title required")
		return
	}
	slug := strings.TrimSpace(req.Slug)
	if slug == "" {
		slug = slugify(name)
	} else {
		slug = slugify(slug)
	}
	if slug == "" {
		writeError(c, http.StatusBadRequest, "invalid slug")
		return
	}
	attJSON, err := encodeAttachments(req.Attachments)
	if err != nil {
		writeError(c, http.StatusBadRequest, "invalid attachments")
		return
	}
	enabled := true
	if req.Enabled != nil {
		enabled = *req.Enabled
	}

	ctx := c.Request.Context()
	tx, err := h.db.BeginTx(ctx, nil)
	if err != nil {
		writeError(c, http.StatusInternalServerError, "could not start transaction")
		return
	}
	defer func() { _ = tx.Rollback() }()

	var quizID string
	err = tx.QueryRowContext(ctx, `
		INSERT INTO quizzes (name, slug, category, description, attachments, enabled)
		VALUES ($1, $2, $3, $4, $5::jsonb, $6)
		ON CONFLICT (slug) DO UPDATE SET
			name = EXCLUDED.name,
			category = EXCLUDED.category,
			description = EXCLUDED.description,
			attachments = EXCLUDED.attachments,
			enabled = EXCLUDED.enabled,
			updated_at = NOW()
		RETURNING id
	`, name, slug, strings.TrimSpace(req.Category), strings.TrimSpace(req.Description), string(attJSON), enabled).Scan(&quizID)
	if err != nil {
		writeError(c, http.StatusInternalServerError, "could not upsert quiz")
		return
	}

	inserted := 0
	for i, q := range req.Questions {
		correct := q.CorrectIndex
		if correct < 0 || correct > 3 {
			writeError(c, http.StatusBadRequest, "correct_index must be 0-3")
			return
		}
		diff := q.Difficulty
		if d, ok := normalizeDifficulty(diff); ok {
			diff = d
		} else if strings.TrimSpace(diff) == "" {
			diff = "medium"
		} else {
			writeError(c, http.StatusBadRequest, "invalid difficulty")
			return
		}
		sortOrder := q.SortOrder
		if sortOrder == 0 {
			sortOrder = i + 1
		}
		_, err = tx.ExecContext(ctx, `
			INSERT INTO questions (
				quiz_id, prompt, option_a, option_b, option_c, option_d,
				correct_index, question_type, difficulty, sort_order
			) VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10)
		`, quizID, strings.TrimSpace(q.Prompt), strings.TrimSpace(q.OptionA), strings.TrimSpace(q.OptionB),
			strings.TrimSpace(q.OptionC), strings.TrimSpace(q.OptionD), correct,
			normalizeQuestionType(q.QuestionType), diff, sortOrder)
		if err != nil {
			writeError(c, http.StatusInternalServerError, "could not insert question")
			return
		}
		inserted++
	}

	if err := tx.Commit(); err != nil {
		writeError(c, http.StatusInternalServerError, "could not commit")
		return
	}

	quiz, err := h.getQuizBySlug(c, slug)
	if err == sql.ErrNoRows {
		writeError(c, http.StatusInternalServerError, "quiz missing after ingest")
		return
	}
	if err != nil {
		writeError(c, http.StatusInternalServerError, "could not load quiz")
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"quiz":              quiz,
		"questions_inserted": inserted,
	})
}
