package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"url-shortener/internal/storage"

	"github.com/jackc/pgx/v5"
)

type ShortRequest struct {
	OriginUrl string `json:"OriginUrl"`
}

type ShortResponce struct {
	ShortUrl string `json:"ShortUrl"`
}

type Handlers struct {
	Conn *pgx.Conn
}

func (h *Handlers) CreateUrlShort(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	var req ShortRequest
	json.NewDecoder(r.Body).Decode(&req)

	if req.OriginUrl == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	Short_code := storage.GenerateShortCode(6)

	sql_query := `INSERT INTO url_short (original_url,short_code) VALUES ($1,$2)`

	_, err := h.Conn.Exec(context.Background(), sql_query, req.OriginUrl, Short_code)
	if err != nil {
		fmt.Println("Error", err)
	}

	resp := ShortResponce{
		ShortUrl: "http://localhost:8060/" + Short_code,
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(resp)
}

func (h *Handlers) RedirectUrl(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	Url := strings.TrimPrefix(r.URL.Path, "/")

	if Url == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	var orig_url string

	sql_query := `SELECT original_url FROM url_short WHERE short_code = $1`
	err := h.Conn.QueryRow(context.Background(), sql_query, Url).Scan(&orig_url)
	if err != nil {
		http.Error(w, "Link not found", http.StatusNotFound)
		return
	}

	updateQuery := `UPDATE url_short SET clicks = clicks + 1 WHERE short_code = $1`
	_, _ = h.Conn.Exec(context.Background(), updateQuery, Url)

	http.Redirect(w, r, orig_url, http.StatusFound)
}

func (h *Handlers) GetStat(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	Url := strings.TrimPrefix(r.URL.Path, "/stats/")
	if Url == "" {
		http.Error(w, "Short code is required", http.StatusBadRequest)
		return
	}

	var clicks int
	var orig_url string
	sql_query := `SELECT clicks, original_url FROM url_short WHERE short_code = $1`
	err := h.Conn.QueryRow(context.Background(), sql_query, Url).Scan(&clicks, &orig_url)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	resp := map[string]interface{}{
		"origin_url": orig_url,
		"short_code": Url,
		"clicks":     clicks,
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(resp)

}

func (h *Handlers) UrlDelete(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	Url := strings.TrimPrefix(r.URL.Path, "/delete/")
	if Url == "" {
		http.Error(w, "Short code is required", http.StatusBadRequest)
		return
	}

	sql_query := `DELETE FROM url_short WHERE short_code = $1`
	_, err := h.Conn.Exec(context.Background(), sql_query, Url)
	if err != nil {
		http.Error(w, "Failed to delete", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
