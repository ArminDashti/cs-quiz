package main

import (
	"context"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/ArminDashti/cs-quiz-api/internal/auth"
	"github.com/ArminDashti/cs-quiz-api/internal/config"
	appdb "github.com/ArminDashti/cs-quiz-api/internal/db"
	"github.com/ArminDashti/cs-quiz-api/internal/handlers"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	config.LoadDotEnv(".env")

	cfg := config.Load()
	if cfg.DatabaseURL == "" {
		log.Fatal("DATABASE_URL is required")
	}

	sqlDB, err := appdb.Open(cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("database: %v", err)
	}
	defer sqlDB.Close()

	if err := appdb.Migrate(sqlDB, cfg.MigrationsDir); err != nil {
		log.Fatalf("migrate: %v", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := appdb.SeedInviteCode(ctx, sqlDB, cfg.InviteCode); err != nil {
		log.Fatalf("seed invite code: %v", err)
	}

	defaultHash, err := auth.HashPassword(appdb.DefaultPassword)
	if err != nil {
		log.Fatalf("hash default password: %v", err)
	}
	if err := appdb.SeedDefaultUser(ctx, sqlDB, defaultHash); err != nil {
		log.Fatalf("seed default user: %v", err)
	}

	if err := os.MkdirAll(filepath.Join(cfg.UploadDir, "avatars"), 0o755); err != nil {
		log.Fatalf("upload dir: %v", err)
	}

	h := handlers.New(sqlDB, cfg)
	r := gin.Default()

	r.Use(cors.New(cors.Config{
		AllowOriginFunc: func(origin string) bool {
			return true
		},
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: false,
		MaxAge:           12 * time.Hour,
	}))

	r.Static("/uploads", cfg.UploadDir)
	r.GET("/health", h.Health)

	api := r.Group("/api/v1")
	{
		api.POST("/auth/register", h.Register)
		api.POST("/auth/login", h.Login)
		api.GET("/profiles/:username", h.GetProfile)
		api.GET("/quizzes", h.ListQuizzes)
		api.GET("/quizzes/:slug", h.GetQuiz)
		api.POST("/external/quizzes", h.ExternalIngest)

		authed := api.Group("")
		authed.Use(auth.Middleware(cfg.JWTSecret))
		{
			authed.GET("/auth/me", h.Me)
			authed.PATCH("/account", h.UpdateAccount)
			authed.POST("/account/avatar", h.UploadAvatar)
			authed.POST("/account/password", h.ChangePassword)
			authed.DELETE("/account", h.DeleteAccount)

			authed.GET("/quizzes/:slug/play", h.PlayQuiz)
			authed.POST("/quizzes/:slug/submit", h.SubmitQuiz)

			admin := authed.Group("/admin")
			admin.Use(auth.RequireAdmin())
			{
				admin.GET("/invite-code", h.GetInviteCode)
				admin.PUT("/invite-code", h.UpdateInviteCode)
				admin.GET("/quizzes", h.AdminListQuizzes)
				admin.POST("/quizzes", h.AdminCreateQuiz)
				admin.PATCH("/quizzes/:slug", h.AdminUpdateQuiz)
				admin.DELETE("/quizzes/:slug", h.AdminDeleteQuiz)
				admin.GET("/quizzes/:slug/questions", h.AdminListQuestions)
				admin.POST("/quizzes/:slug/questions", h.AdminCreateQuestion)
				admin.PATCH("/quizzes/:slug/questions/:id", h.AdminUpdateQuestion)
				admin.DELETE("/quizzes/:slug/questions/:id", h.AdminDeleteQuestion)
			}
		}
	}

	log.Printf("cs-quiz-api listening on %s", cfg.Addr)
	if err := r.Run(cfg.Addr); err != nil {
		log.Fatalf("server: %v", err)
	}
}
