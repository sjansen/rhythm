package api

import (
	"net/http"

	"github.com/go-chi/chi"
)

type Config struct {
	Secret string `ignored:"true"`
}

func New(cfg *Config) *chi.Mux {
	r := chi.NewRouter()
	r.Get("/secret", cfg.secret)
	return r
}

func (cfg *Config) secret(w http.ResponseWriter, r *http.Request) {
	h := w.Header()
	h.Add("Cache-Control", "no-cache, no-store, must-revalidate")
	h.Add("Content-Type", "text/text; charset=utf-8")
	w.Write([]byte(cfg.Secret))
}
