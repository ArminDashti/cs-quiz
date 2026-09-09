package handlers

import (
	"database/sql"
	"net/http"
	"strings"

	appdb "github.com/ArminDashti/cs-quiz-api/internal/db"
	"github.com/ArminDashti/cs-quiz-api/internal/models"
	"github.com/gin-gonic/gin"
)

// GetInviteCode returns the current invite code (admin only).
func (h *Handler) GetInviteCode(c *gin.Context) {
	code, err := appdb.GetSetting(c.Request.Context(), h.db, "invite_code")
	if err != nil {
		writeError(c, http.StatusInternalServerError, "could not load invite code")
		return
	}
	c.JSON(http.StatusOK, gin.H{"invite_code": code})
}

// UpdateInviteCode sets the invite code (admin only).
func (h *Handler) UpdateInviteCode(c *gin.Context) {
	var req struct {
		InviteCode string `json:"invite_code" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		writeError(c, http.StatusBadRequest, "invalid request")
		return
	}
	code := strings.TrimSpace(req.InviteCode)
	if code == "" {
		writeError(c, http.StatusBadRequest, "invite_code required")
		return
	}
	if err := appdb.SetSetting(c.Request.Context(), h.db, "invite_code", code); err != nil {
		writeError(c, http.StatusInternalServerError, "could not update invite code")
		return
	}
	c.JSON(http.StatusOK, gin.H{"invite_code": code})
}

// AdminListQuizzes returns all quizzes including disabled (admin).
func (h *Handler) AdminListQuizzes(c *gin.Context) {
	rows, err := h.db.QueryContext(c.Request.Context(), quizSelect+` ORDER BY name`)
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

// AdminCreateQuiz creates a quiz.
func (h *Handler) AdminCreateQuiz(c *gin.Context) {
	var req struct {
		Name        string   `json:"name"`
		Title       string   `json:"title"`
		Slug        string   `json:"slug"`
		Category    string   `json:"category"`
		Description string   `json:"description"`
		Attachments []string `json:"attachments"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		writeError(c, http.StatusBadRequest, "invalid request")
		return
	}
	name := strings.TrimSpace(req.Name)
	if name == "" {
		name = strings.TrimSpace(req.Title)
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

	var id string
	err = h.db.QueryRowContext(c.Request.Context(), `
		INSERT INTO quizzes (name, slug, category, description, attachments)
		VALUES ($1, $2, $3, $4, $5::jsonb)
		RETURNING id
	`, name, slug, strings.TrimSpace(req.Category), strings.TrimSpace(req.Description), string(attJSON)).Scan(&id)
	if err != nil {
		if fmtUniqueViolation(err) {
			writeError(c, http.StatusConflict, "quiz slug already exists")
			return
		}
		writeError(c, http.StatusInternalServerError, "could not create quiz")
		return
	}
	q, err := h.getQuizBySlug(c, slug)
	if err != nil {
		writeError(c, http.StatusInternalServerError, "could not load quiz")
		return
	}
	c.JSON(http.StatusCreated, q)
}

// AdminUpdateQuiz patches a quiz by slug.
func (h *Handler) AdminUpdateQuiz(c *gin.Context) {
	slug := strings.TrimSpace(c.Param("slug"))
	q, err := h.getQuizBySlug(c, slug)
	if err == sql.ErrNoRows {
		writeError(c, http.StatusNotFound, "quiz not found")
		return
	}
	if err != nil {
		writeError(c, http.StatusInternalServerError, "could not load quiz")
		return
	}

	var req struct {
		Name        *string  `json:"name"`
		Title       *string  `json:"title"`
		Slug        *string  `json:"slug"`
		Category    *string  `json:"category"`
		Description *string  `json:"description"`
		Attachments *[]string `json:"attachments"`
		Enabled     *bool    `json:"enabled"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		writeError(c, http.StatusBadRequest, "invalid request")
		return
	}

	name := q.Name
	newSlug := q.Slug
	category := q.Category
	desc := q.Description
	enabled := q.Enabled
	attJSON, err := encodeAttachments(q.Attachments)
	if err != nil {
		writeError(c, http.StatusInternalServerError, "invalid attachments")
		return
	}
	if req.Name != nil {
		name = strings.TrimSpace(*req.Name)
	}
	if req.Title != nil && strings.TrimSpace(*req.Title) != "" {
		name = strings.TrimSpace(*req.Title)
	}
	if req.Slug != nil {
		newSlug = slugify(*req.Slug)
		if newSlug == "" {
			writeError(c, http.StatusBadRequest, "invalid slug")
			return
		}
	}
	if req.Category != nil {
		category = strings.TrimSpace(*req.Category)
	}
	if req.Description != nil {
		desc = strings.TrimSpace(*req.Description)
	}
	if req.Attachments != nil {
		attJSON, err = encodeAttachments(*req.Attachments)
		if err != nil {
			writeError(c, http.StatusBadRequest, "invalid attachments")
			return
		}
	}
	if req.Enabled != nil {
		enabled = *req.Enabled
	}

	_, err = h.db.ExecContext(c.Request.Context(), `
		UPDATE quizzes SET name = $1, slug = $2, category = $3, description = $4, attachments = $5::jsonb, enabled = $6, updated_at = NOW()
		WHERE id = $7
	`, name, newSlug, category, desc, string(attJSON), enabled, q.ID)
	if err != nil {
		if fmtUniqueViolation(err) {
			writeError(c, http.StatusConflict, "quiz slug already exists")
			return
		}
		writeError(c, http.StatusInternalServerError, "could not update quiz")
		return
	}
	q, err = h.getQuizBySlug(c, newSlug)
	if err != nil {
		writeError(c, http.StatusInternalServerError, "could not load quiz")
		return
	}
	c.JSON(http.StatusOK, q)
}

// AdminDeleteQuiz deletes a quiz by slug.
func (h *Handler) AdminDeleteQuiz(c *gin.Context) {
	slug := strings.TrimSpace(c.Param("slug"))
	res, err := h.db.ExecContext(c.Request.Context(), `DELETE FROM quizzes WHERE slug = $1`, slug)
	if err != nil {
		writeError(c, http.StatusInternalServerError, "could not delete quiz")
		return
	}
	n, _ := res.RowsAffected()
	if n == 0 {
		writeError(c, http.StatusNotFound, "quiz not found")
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}

const questionAdminSelect = `
	SELECT id, quiz_id, prompt, option_a, option_b, option_c, option_d, correct_index,
	       question_type, difficulty, sort_order, created_at, updated_at
	FROM questions`

func scanQuestionAdmin(row interface {
	Scan(dest ...any) error
}) (models.Question, error) {
	var q models.Question
	var correct int
	err := row.Scan(
		&q.ID, &q.QuizID, &q.Prompt, &q.OptionA, &q.OptionB, &q.OptionC, &q.OptionD,
		&correct, &q.QuestionType, &q.Difficulty, &q.SortOrder, &q.CreatedAt, &q.UpdatedAt,
	)
	if err != nil {
		return q, err
	}
	q.CorrectIndex = &correct
	return q, err
}

func normalizeDifficulty(d string) (string, bool) {
	switch strings.ToLower(strings.TrimSpace(d)) {
	case "easy", "medium", "hard":
		return strings.ToLower(strings.TrimSpace(d)), true
	default:
		return "", false
	}
}

func normalizeQuestionType(t string) string {
	t = strings.TrimSpace(t)
	if t == "" {
		return "multiple_choice"
	}
	return t
}

// AdminListQuestions lists questions for a quiz including answers.
func (h *Handler) AdminListQuestions(c *gin.Context) {
	slug := strings.TrimSpace(c.Param("slug"))
	quiz, err := h.getQuizBySlug(c, slug)
	if err == sql.ErrNoRows {
		writeError(c, http.StatusNotFound, "quiz not found")
		return
	}
	if err != nil {
		writeError(c, http.StatusInternalServerError, "could not load quiz")
		return
	}

	rows, err := h.db.QueryContext(c.Request.Context(), questionAdminSelect+`
		WHERE quiz_id = $1
		ORDER BY sort_order, created_at
	`, quiz.ID)
	if err != nil {
		writeError(c, http.StatusInternalServerError, "could not list questions")
		return
	}
	defer rows.Close()

	out := make([]models.Question, 0)
	for rows.Next() {
		q, err := scanQuestionAdmin(rows)
		if err != nil {
			writeError(c, http.StatusInternalServerError, "could not list questions")
			return
		}
		out = append(out, q)
	}
	c.JSON(http.StatusOK, out)
}

// AdminCreateQuestion adds a question to a quiz.
func (h *Handler) AdminCreateQuestion(c *gin.Context) {
	slug := strings.TrimSpace(c.Param("slug"))
	quiz, err := h.getQuizBySlug(c, slug)
	if err == sql.ErrNoRows {
		writeError(c, http.StatusNotFound, "quiz not found")
		return
	}
	if err != nil {
		writeError(c, http.StatusInternalServerError, "could not load quiz")
		return
	}

	var req struct {
		Prompt       string `json:"prompt" binding:"required"`
		OptionA      string `json:"option_a" binding:"required"`
		OptionB      string `json:"option_b" binding:"required"`
		OptionC      string `json:"option_c" binding:"required"`
		OptionD      string `json:"option_d" binding:"required"`
		CorrectIndex int    `json:"correct_index"`
		QuestionType string `json:"question_type"`
		Difficulty   string `json:"difficulty"`
		SortOrder    *int   `json:"sort_order"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		writeError(c, http.StatusBadRequest, "invalid request")
		return
	}
	if req.CorrectIndex < 0 || req.CorrectIndex > 3 {
		writeError(c, http.StatusBadRequest, "correct_index must be 0-3")
		return
	}
	qType := normalizeQuestionType(req.QuestionType)
	diff := "medium"
	if req.Difficulty != "" {
		var ok bool
		diff, ok = normalizeDifficulty(req.Difficulty)
		if !ok {
			writeError(c, http.StatusBadRequest, "difficulty must be easy, medium, or hard")
			return
		}
	}
	sortOrder := 0
	if req.SortOrder != nil {
		sortOrder = *req.SortOrder
	} else {
		_ = h.db.QueryRowContext(c.Request.Context(), `
			SELECT COALESCE(MAX(sort_order), 0) + 1 FROM questions WHERE quiz_id = $1
		`, quiz.ID).Scan(&sortOrder)
	}

	var q models.Question
	var correct int
	err = h.db.QueryRowContext(c.Request.Context(), `
		INSERT INTO questions (quiz_id, prompt, option_a, option_b, option_c, option_d, correct_index, question_type, difficulty, sort_order)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
		RETURNING id, quiz_id, prompt, option_a, option_b, option_c, option_d, correct_index, question_type, difficulty, sort_order, created_at, updated_at
	`, quiz.ID, strings.TrimSpace(req.Prompt), req.OptionA, req.OptionB, req.OptionC, req.OptionD, req.CorrectIndex, qType, diff, sortOrder).Scan(
		&q.ID, &q.QuizID, &q.Prompt, &q.OptionA, &q.OptionB, &q.OptionC, &q.OptionD,
		&correct, &q.QuestionType, &q.Difficulty, &q.SortOrder, &q.CreatedAt, &q.UpdatedAt,
	)
	if err != nil {
		writeError(c, http.StatusInternalServerError, "could not create question")
		return
	}
	q.CorrectIndex = &correct
	c.JSON(http.StatusCreated, q)
}

// AdminUpdateQuestion updates a question.
func (h *Handler) AdminUpdateQuestion(c *gin.Context) {
	slug := strings.TrimSpace(c.Param("slug"))
	qid := strings.TrimSpace(c.Param("id"))
	quiz, err := h.getQuizBySlug(c, slug)
	if err == sql.ErrNoRows {
		writeError(c, http.StatusNotFound, "quiz not found")
		return
	}
	if err != nil {
		writeError(c, http.StatusInternalServerError, "could not load quiz")
		return
	}

	existing, err := scanQuestionAdmin(h.db.QueryRowContext(c.Request.Context(), questionAdminSelect+`
		WHERE id = $1 AND quiz_id = $2
	`, qid, quiz.ID))
	if err == sql.ErrNoRows {
		writeError(c, http.StatusNotFound, "question not found")
		return
	}
	if err != nil {
		writeError(c, http.StatusInternalServerError, "could not load question")
		return
	}

	var req struct {
		Prompt       *string `json:"prompt"`
		OptionA      *string `json:"option_a"`
		OptionB      *string `json:"option_b"`
		OptionC      *string `json:"option_c"`
		OptionD      *string `json:"option_d"`
		CorrectIndex *int    `json:"correct_index"`
		QuestionType *string `json:"question_type"`
		Difficulty   *string `json:"difficulty"`
		SortOrder    *int    `json:"sort_order"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		writeError(c, http.StatusBadRequest, "invalid request")
		return
	}

	prompt := existing.Prompt
	a, b, cOpt, d := existing.OptionA, existing.OptionB, existing.OptionC, existing.OptionD
	correct := *existing.CorrectIndex
	qType := existing.QuestionType
	diff := existing.Difficulty
	sortOrder := existing.SortOrder
	if req.Prompt != nil {
		prompt = strings.TrimSpace(*req.Prompt)
	}
	if req.OptionA != nil {
		a = *req.OptionA
	}
	if req.OptionB != nil {
		b = *req.OptionB
	}
	if req.OptionC != nil {
		cOpt = *req.OptionC
	}
	if req.OptionD != nil {
		d = *req.OptionD
	}
	if req.CorrectIndex != nil {
		if *req.CorrectIndex < 0 || *req.CorrectIndex > 3 {
			writeError(c, http.StatusBadRequest, "correct_index must be 0-3")
			return
		}
		correct = *req.CorrectIndex
	}
	if req.QuestionType != nil {
		qType = normalizeQuestionType(*req.QuestionType)
	}
	if req.Difficulty != nil {
		var ok bool
		diff, ok = normalizeDifficulty(*req.Difficulty)
		if !ok {
			writeError(c, http.StatusBadRequest, "difficulty must be easy, medium, or hard")
			return
		}
	}
	if req.SortOrder != nil {
		sortOrder = *req.SortOrder
	}

	var q models.Question
	var correctOut int
	err = h.db.QueryRowContext(c.Request.Context(), `
		UPDATE questions
		SET prompt = $1, option_a = $2, option_b = $3, option_c = $4, option_d = $5,
		    correct_index = $6, question_type = $7, difficulty = $8, sort_order = $9, updated_at = NOW()
		WHERE id = $10
		RETURNING id, quiz_id, prompt, option_a, option_b, option_c, option_d, correct_index, question_type, difficulty, sort_order, created_at, updated_at
	`, prompt, a, b, cOpt, d, correct, qType, diff, sortOrder, existing.ID).Scan(
		&q.ID, &q.QuizID, &q.Prompt, &q.OptionA, &q.OptionB, &q.OptionC, &q.OptionD,
		&correctOut, &q.QuestionType, &q.Difficulty, &q.SortOrder, &q.CreatedAt, &q.UpdatedAt,
	)
	if err != nil {
		writeError(c, http.StatusInternalServerError, "could not update question")
		return
	}
	q.CorrectIndex = &correctOut
	c.JSON(http.StatusOK, q)
}

// AdminDeleteQuestion deletes a question.
func (h *Handler) AdminDeleteQuestion(c *gin.Context) {
	slug := strings.TrimSpace(c.Param("slug"))
	qid := strings.TrimSpace(c.Param("id"))
	quiz, err := h.getQuizBySlug(c, slug)
	if err == sql.ErrNoRows {
		writeError(c, http.StatusNotFound, "quiz not found")
		return
	}
	if err != nil {
		writeError(c, http.StatusInternalServerError, "could not load quiz")
		return
	}
	res, err := h.db.ExecContext(c.Request.Context(), `
		DELETE FROM questions WHERE id = $1 AND quiz_id = $2
	`, qid, quiz.ID)
	if err != nil {
		writeError(c, http.StatusInternalServerError, "could not delete question")
		return
	}
	n, _ := res.RowsAffected()
	if n == 0 {
		writeError(c, http.StatusNotFound, "question not found")
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}
