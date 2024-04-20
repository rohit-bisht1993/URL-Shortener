package urlshortener

import (
	"encoding/json"
	"net/http"
)

type UrlShortnerContext struct {
}

func NewUrlShortener() *UrlShortnerContext {
	return &UrlShortnerContext{}
}

func (urlCntx *UrlShortnerContext) UrlShortenerAPI(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusCreated)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode("Success Resp")
}
