package main

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/httplog"
	"github.com/go-chi/render"
	"github.com/rs/zerolog"
)

var (
	oplog       zerolog.Logger
	appName     = "myapp"
	servicePort = os.Getenv("PORT")
)

func main() {

	oplog = httplog.LogEntry(context.Background())

	/* jsonify logging */
	httpLogger := httplog.NewLogger(
		appName,
		httplog.Options{
			JSON:           true,
			LogLevel:       slog.LevelInfo.String(),
			LevelFieldName: "severity",
			Concise:        true,
		})

	oplog.Debug().Msg("logger initialized")

	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(60 * time.Second))
	r.Use(httplog.RequestLogger(httpLogger))

	r.Get("/ping", pingPong)

	r.Get("/hostname", func(w http.ResponseWriter, r *http.Request) {
		oplog.Info().Str("path", "/hostname").Send()
		host, err := os.Hostname()
		if err != nil {
			errorRender(w, r, http.StatusInternalServerError, err)
			return
		}
		render.HTML(w, r, fmt.Sprintf("<h1>My name is %s</h1>", host))
	})

	if servicePort == "" {
		servicePort = "8080"
	}

	oplog.Debug().Msg("starting api listening on port " + servicePort)
	if err := http.ListenAndServe(":"+servicePort, r); err != nil {
		oplog.Err(err)
	}

}

var errorRender = func(w http.ResponseWriter, r *http.Request, httpCode int, err error) {
	oplog.Error().Str("path", r.URL.Path).Send()
	render.Status(r, httpCode)
	render.JSON(w, r, map[string]any{"ERROR": err.Error()})
}

func pingPong(w http.ResponseWriter, r *http.Request) {
	oplog.Info().Str("path", "/ping").Send()
	render.Status(r, http.StatusOK)
	render.PlainText(w, r, "Pong\n")
}
