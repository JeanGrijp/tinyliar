package routes

import (
	"net/http"
	"time"

	"github.com/JeanGrijp/tinyliar/internal/handler"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func SetupRoutes(linkHandler *handler.LinkHandler) http.Handler {
	r := chi.NewRouter()

	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(10 * time.Second))

	r.Post("/", linkHandler.CreateLinkHandler)
	r.Get("/{short_url}", linkHandler.GetLinkHandler)

	return r
}
