package handler

import (
	"encoding/json"
	"log/slog"
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
	slog.Info("Validando URL", "url", raw)
	u, err := url.ParseRequestURI(raw)
	if err != nil {
		slog.Error("Erro ao validar URL", "error", err)
		return false
	}

	slog.Info("URL validada", "url", u.String())

	return (strings.HasPrefix(u.Scheme, "http"))
}

type CreateLinkRequest struct {
	OriginalURL string `json:"original_url"`
}

func (h *LinkHandler) CreateLinkHandler(w http.ResponseWriter, r *http.Request) {

	slog.InfoContext(r.Context(), "Criando link encurtado")

	var originalURL string

	slog.InfoContext(r.Context(), "Tentando pegar URL da query string")
	originalURL = r.URL.Query().Get("shorten")
	slog.InfoContext(r.Context(), "URL da query string", "url", originalURL)

	if originalURL == "" && r.Header.Get("Content-Type") == "application/json" {
		slog.InfoContext(r.Context(), "Tentando pegar URL do corpo da requisição")
		var req CreateLinkRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			slog.ErrorContext(r.Context(), "Erro ao decodificar JSON", "error", err)
			http.Error(w, "Invalid JSON body", http.StatusBadRequest)
			return
		}
		slog.InfoContext(r.Context(), "URL do corpo da requisição", "url", req.OriginalURL)
		originalURL = req.OriginalURL
	}

	if originalURL == "" {
		slog.ErrorContext(r.Context(), "URL original não informada")
		http.Error(w, "Missing original URL", http.StatusBadRequest)
		return
	}

	if !isValidURL(originalURL) {
		slog.ErrorContext(r.Context(), "URL inválida", "url", originalURL)
		http.Error(w, "Invalid URL", http.StatusBadRequest)
		return
	}

	slog.InfoContext(r.Context(), "URL válida", "url", originalURL)
	shortID := utils.GenerateShortID()

	slog.InfoContext(r.Context(), "Gerando ID curto", "short_id", shortID)
	link := &model.Link{
		OriginalURL: originalURL,
		ShortURL:    shortID,
		Clicks:      0,
		OwnerID:     0,
		ExpiredAt:   time.Now().Add(10 * 24 * time.Hour),
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	slog.InfoContext(r.Context(), "Criando link", "link", link)

	existsLink, err := h.Repo.GetyByOriginalURL(originalURL)
	if err != nil {
		slog.ErrorContext(r.Context(), "Erro ao buscar link", "error", err)
		http.Error(w, "Failed to check existing link", http.StatusInternalServerError)
		return
	}

	if existsLink != nil {
		slog.InfoContext(r.Context(), "Link já existe", "link", existsLink)
		slog.InfoContext(r.Context(), "Retornando link existente", "link", existsLink)
		resp := map[string]string{"short_url": existsLink.ShortURL}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(resp)
		return
	}

	if err := h.Repo.CreateLink(link); err != nil {
		slog.ErrorContext(r.Context(), "Erro ao criar link", "error", err)
		http.Error(w, "Failed to save link", http.StatusInternalServerError)
		return
	}

	slog.InfoContext(r.Context(), "Link criado com sucesso", "link", link)
	resp := map[string]string{"short_url": shortID}
	w.Header().Set("Content-Type", "application/json")

	json.NewEncoder(w).Encode(resp)
}

func (h *LinkHandler) GetLinkHandler(w http.ResponseWriter, r *http.Request) {
	slog.InfoContext(r.Context(), "Redirecionando para link encurtado")
	shortID := chi.URLParam(r, "short_url")
	slog.InfoContext(r.Context(), "ID curto", "short_id", shortID)
	if shortID == "" {
		slog.ErrorContext(r.Context(), "ID curto não informado")
		http.Error(w, "Missing short URL", http.StatusBadRequest)
		return
	}

	slog.InfoContext(r.Context(), "Buscando link no banco de dados", "short_id", shortID)
	link, err := h.Repo.GetLinkByShortURL(shortID)
	if err != nil || link == nil {
		slog.ErrorContext(r.Context(), "Erro ao buscar link", "error", err)
		http.Error(w, "Link not found", http.StatusNotFound)
		return
	}
	slog.InfoContext(r.Context(), "Link encontrado", "link", link)
	slog.InfoContext(r.Context(), "Verificando se o link está expirado", "expired_at", link.ExpiredAt)

	if !link.ExpiredAt.IsZero() && link.ExpiredAt.Before(time.Now()) {
		slog.ErrorContext(r.Context(), "Link expirado", "expired_at", link.ExpiredAt)
		http.Error(w, "Link expired", http.StatusGone)
		return
	}

	slog.InfoContext(r.Context(), "Link não expirado", "link", link)

	slog.InfoContext(r.Context(), "Incrementando contagem de cliques", "link_id", link.ID)
	_ = h.Repo.IncrementClickCount(link.ID)

	slog.InfoContext(r.Context(), "Redirecionando para URL original", "original_url", link.OriginalURL)
	http.Redirect(w, r, link.OriginalURL, http.StatusFound)
}
