package main

import (
	"context"
	"database/sql"
	"errors"
	"log/slog"
	"net/http"
	"os/signal"
	"skillsapi/skill"
	"syscall"
	"time"

	"log"
	"os"

	"github.com/gin-gonic/gin"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	var err error
	db, err := sql.Open("postgres", os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Panic(err)
	}

	defer db.Close()

	if err = db.Ping(); err != nil {
		log.Panic(err)
	}

	h := &skill.Handler{Db: db}
	r := gin.Default()
	skill.SetRouter(r, h)

	srv := http.Server{
		Addr:              ":" + os.Getenv("PORT"),
		Handler:           r,
		ReadHeaderTimeout: 3 * time.Second,
	}

	go func() {
		<-ctx.Done()

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		slog.Info("Shutting down server")
		if err := srv.Shutdown(ctx); err != nil {
			log.Fatal(err)
		}
	}()

	if err := srv.ListenAndServe(); err != nil {
		if !errors.Is(err, http.ErrServerClosed) {
			log.Panic(err)
		}
	}

	slog.Info("Server exiting")
}
