package models

import "time"

// User is a persisted account.
type User struct {
	ID           string    `json:"id"`
	Email        string    `json:"email"`
	PasswordHash string    `json:"-"`
	FirstName    string    `json:"first_name"`
	LastName     string    `json:"last_name"`
	Username     *string   `json:"username"`
	AvatarPath   *string   `json:"avatar_path,omitempty"`
	AvatarURL    *string   `json:"avatar_url,omitempty"`
	IsAdmin      bool      `json:"is_admin"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

// Quiz is a named quiz topic.
type Quiz struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"` // title in the UI
	Slug        string    `json:"slug"`
	Category    string    `json:"category"`
	Description string    `json:"description"`
	Attachments []string  `json:"attachments"`
	Enabled     bool      `json:"enabled"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// Question is a multiple-choice item (admin view includes correct_index).
type Question struct {
	ID           string    `json:"id"`
	QuizID       string    `json:"quiz_id"`
	Prompt       string    `json:"prompt"`
	OptionA      string    `json:"option_a"`
	OptionB      string    `json:"option_b"`
	OptionC      string    `json:"option_c"`
	OptionD      string    `json:"option_d"`
	CorrectIndex *int      `json:"correct_index,omitempty"`
	QuestionType string    `json:"question_type"`
	Difficulty   string    `json:"difficulty"`
	SortOrder    int       `json:"sort_order"`
	CreatedAt    time.Time `json:"created_at,omitempty"`
	UpdatedAt    time.Time `json:"updated_at,omitempty"`
}

// Score is a saved quiz attempt.
type Score struct {
	ID        string    `json:"id"`
	UserID    string    `json:"user_id"`
	QuizID    string    `json:"quiz_id"`
	QuizSlug  string    `json:"quiz_slug,omitempty"`
	Correct   int       `json:"correct"`
	Total     int       `json:"total"`
	CreatedAt time.Time `json:"created_at"`
}
