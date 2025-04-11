package handler

import (
	"encoding/json"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/JeanGrijp/tinyliar/internal/model"
	"github.com/JeanGrijp/tinyliar/internal/repository"
	"github.com/JeanGrijp/tinyliar/internal/utils"
	"github.com/go-chi/chi/v5"
)

type LinkHandler struct {
	Repo *repository.LinkRepository
}

func isValidURL(raw string) bool {
	u, err := url.ParseRequestURI(raw)
	return err == nil && (strings.HasPrefix(u.Scheme, "http"))
}

type CreateLinkRequest struct {
	OriginalURL string `json:"original_url"`
}

func (h *LinkHandler) CreateLinkHandler(w http.ResponseWriter, r *http.Request) {
	var originalURL string

	// Primeiro tenta pegar da query string
	originalURL = r.URL.Query().Get("shorten")

	// Se não tiver na query, tenta pegar do corpo
	if originalURL == "" && r.Header.Get("Content-Type") == "application/json" {
		var req CreateLinkRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "Invalid JSON body", http.StatusBadRequest)
			return
		}
		originalURL = req.OriginalURL
	}

	if originalURL == "" {
		http.Error(w, "Missing original URL", http.StatusBadRequest)
		return
	}

	if !isValidURL(originalURL) {
		http.Error(w, "Invalid URL", http.StatusBadRequest)
		return
	}

	shortID := utils.GenerateShortID()

	link := &model.Link{
		OriginalURL: originalURL,
		ShortURL:    shortID,
		Clicks:      0,
		OwnerID:     0,
		ExpiredAt:   "",
		CreatedAt:   time.Now().Format(time.RFC3339),
		UpdatedAt:   time.Now().Format(time.RFC3339),
	}

	if err := h.Repo.CreateLink(link); err != nil {
		http.Error(w, "Failed to save link", http.StatusInternalServerError)
		return
	}

	resp := map[string]string{"short_url": shortID}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

func (h *LinkHandler) GetLinkHandler(w http.ResponseWriter, r *http.Request) {
	shortID := chi.URLParam(r, "short_url")
	if shortID == "" {
		http.Error(w, "Missing short URL", http.StatusBadRequest)
		return
	}

	link, err := h.Repo.GetLinkByShortURL(shortID)
	if err != nil || link == nil {
		http.Error(w, "Link not found", http.StatusNotFound)
		return
	}

	// Incrementa contagem de cliques (ignorar erro, mas ideal tratar)
	_ = h.Repo.IncrementClickCount(link.ID)

	http.Redirect(w, r, link.OriginalURL, http.StatusFound)
}
