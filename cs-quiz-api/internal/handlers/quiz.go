package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strings"

	"github.com/ArminDashti/cs-quiz-api/internal/auth"
	"github.com/ArminDashti/cs-quiz-api/internal/models"
	"github.com/gin-gonic/gin"
)

func scanQuiz(row interface {
	Scan(dest ...any) error
}) (models.Quiz, error) {
	var q models.Quiz
	var attRaw []byte
	err := row.Scan(&q.ID, &q.Name, &q.Slug, &q.Category, &q.Description, &attRaw, &q.Enabled, &q.CreatedAt, &q.UpdatedAt)
	if err != nil {
		return q, err
	}
	q.Attachments = decodeAttachments(attRaw)
	return q, nil
}

func decodeAttachments(raw []byte) []string {
	if len(raw) == 0 {
		return []string{}
	}
	var out []string
	if err := json.Unmarshal(raw, &out); err != nil || out == nil {
		return []string{}
	}
	return out
}

func encodeAttachments(list []string) ([]byte, error) {
	if list == nil {
		list = []string{}
	}
	clean := make([]string, 0, len(list))
	for _, a := range list {
		a = strings.TrimSpace(a)
		if a != "" {
			clean = append(clean, a)
		}
	}
	return json.Marshal(clean)
}

const quizSelect = `SELECT id, name, slug, category, description, attachments, enabled, created_at, updated_at FROM quizzes`

func (h *Handler) getQuizBySlug(c *gin.Context, slug string) (models.Quiz, error) {
	return scanQuiz(h.db.QueryRowContext(c.Request.Context(), quizSelect+` WHERE slug = $1`, slug))
}

// ListQuizzes returns enabled quizzes (public).
func (h *Handler) ListQuizzes(c *gin.Context) {
	rows, err := h.db.QueryContext(c.Request.Context(), quizSelect+` WHERE enabled = TRUE ORDER BY name`)
	if err != nil {
		writeError(c, http.StatusInternalServerError, "could not list quizzes")
		return
	}
	defer rows.Close()

	out := make([]models.Quiz, 0)
	for rows.Next() {
		q, err := scanQuiz(rows)
		if err != nil {
			writeError(c, http.StatusInternalServerError, "could not list quizzes")
			return
		}
		out = append(out, q)
	}
	c.JSON(http.StatusOK, out)
}

// GetQuiz returns quiz metadata by slug (enabled only).
func (h *Handler) GetQuiz(c *gin.Context) {
	slug := strings.TrimSpace(c.Param("slug"))
	q, err := h.getQuizBySlug(c, slug)
	if err == sql.ErrNoRows || (err == nil && !q.Enabled) {
		writeError(c, http.StatusNotFound, "quiz not found")
		return
	}
	if err != nil {
		writeError(c, http.StatusInternalServerError, "could not load quiz")
		return
	}
	c.JSON(http.StatusOK, q)
}

// PlayQuiz returns questions without correct answers.
func (h *Handler) PlayQuiz(c *gin.Context) {
	slug := strings.TrimSpace(c.Param("slug"))
	q, err := h.getQuizBySlug(c, slug)
	if err == sql.ErrNoRows || (err == nil && !q.Enabled) {
		writeError(c, http.StatusNotFound, "quiz not found")
		return
	}
	if err != nil {
		writeError(c, http.StatusInternalServerError, "could not load quiz")
		return
	}

	rows, err := h.db.QueryContext(c.Request.Context(), `
		SELECT id, quiz_id, prompt, option_a, option_b, option_c, option_d, sort_order
		FROM questions
		WHERE quiz_id = $1
		ORDER BY sort_order, created_at
	`, q.ID)
	if err != nil {
		writeError(c, http.StatusInternalServerError, "could not load questions")
		return
	}
	defer rows.Close()

	questions := make([]gin.H, 0)
	for rows.Next() {
		var id, quizID, prompt, a, b, cOpt, d string
		var sortOrder int
		if err := rows.Scan(&id, &quizID, &prompt, &a, &b, &cOpt, &d, &sortOrder); err != nil {
			writeError(c, http.StatusInternalServerError, "could not load questions")
			return
		}
		questions = append(questions, gin.H{
			"id":         id,
			"quiz_id":    quizID,
			"prompt":     prompt,
			"option_a":   a,
			"option_b":   b,
			"option_c":   cOpt,
			"option_d":   d,
			"sort_order": sortOrder,
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"quiz":      q,
		"questions": questions,
	})
}

// SubmitQuiz grades answers and stores a score.
func (h *Handler) SubmitQuiz(c *gin.Context) {
	slug := strings.TrimSpace(c.Param("slug"))
	var req struct {
		Answers []struct {
			QuestionID    string `json:"question_id" binding:"required"`
			SelectedIndex int    `json:"selected_index"`
		} `json:"answers" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		writeError(c, http.StatusBadRequest, "invalid request")
		return
	}

	q, err := h.getQuizBySlug(c, slug)
	if err == sql.ErrNoRows || (err == nil && !q.Enabled) {
		writeError(c, http.StatusNotFound, "quiz not found")
		return
	}
	if err != nil {
		writeError(c, http.StatusInternalServerError, "could not load quiz")
		return
	}

	rows, err := h.db.QueryContext(c.Request.Context(), `
		SELECT id, correct_index FROM questions WHERE quiz_id = $1
	`, q.ID)
	if err != nil {
		writeError(c, http.StatusInternalServerError, "could not load questions")
		return
	}
	defer rows.Close()

	correctMap := map[string]int{}
	for rows.Next() {
		var id string
		var correct int
		if err := rows.Scan(&id, &correct); err != nil {
			writeError(c, http.StatusInternalServerError, "could not load questions")
			return
		}
		correctMap[id] = correct
	}

	if len(correctMap) == 0 {
		writeError(c, http.StatusBadRequest, "quiz has no questions")
		return
	}

	correctCount := 0
	seen := map[string]bool{}
	for _, ans := range req.Answers {
		if seen[ans.QuestionID] {
			continue
		}
		expected, ok := correctMap[ans.QuestionID]
		if !ok {
			writeError(c, http.StatusBadRequest, "unknown question")
			return
		}
		seen[ans.QuestionID] = true
		if ans.SelectedIndex >= 0 && ans.SelectedIndex <= 3 && ans.SelectedIndex == expected {
			correctCount++
		}
	}

	total := len(correctMap)
	userID := auth.UserIDFromContext(c)
	var scoreID string
	err = h.db.QueryRowContext(c.Request.Context(), `
		INSERT INTO scores (user_id, quiz_id, correct, total)
		VALUES ($1, $2, $3, $4)
		RETURNING id
	`, userID, q.ID, correctCount, total).Scan(&scoreID)
	if err != nil {
		writeError(c, http.StatusInternalServerError, "could not save score")
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"id":        scoreID,
		"quiz_id":   q.ID,
		"quiz_slug": q.Slug,
		"correct":   correctCount,
		"total":     total,
	})
}
