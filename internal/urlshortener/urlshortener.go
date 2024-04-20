package urlshortener

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
	"strings"
	"time"

	"github.com/gorilla/mux"
	"github.com/rohit-bisht1993/URL-Shortener/internal/constant"
	"github.com/rohit-bisht1993/URL-Shortener/internal/utils"
)

type urlInfo struct {
	url string
}

type UrlShortnerContext struct {
	urls map[string]string
}

// NewUrlShortener
func NewUrlShortener() *UrlShortnerContext {
	urls := make(map[string]string)
	return &UrlShortnerContext{urls: urls}
}

// UrlShortenerAPI short url api
func (urlCntx *UrlShortnerContext) UrlShortenerAPI(w http.ResponseWriter, r *http.Request) {

	// Read data from URL
	originalURL := strings.TrimSpace(r.FormValue("url"))
	if len(originalURL) == 0 {
		//Parsing Request body if url is not in url path
		reqData := urlInfo{}
		reqBody, _ := ioutil.ReadAll(r.Body)
		_ = json.Unmarshal(reqBody, &reqData)

		originalURL = strings.TrimSpace(reqData.url)

		if len(originalURL) == 0 {
			http.Error(w, "URL parameter is missing", http.StatusBadRequest)
			return
		}
	}

	shortKey, urlExist := utils.IsValueExist(urlCntx.urls, originalURL)
	if !urlExist {
		// Generate a unique shortened key for the original URL
		shortKey = urlCntx.generateShortKey()
		urlCntx.urls[shortKey] = originalURL
	}

	// Construct the full shortened URL
	shortenedURL := fmt.Sprintf("http://localhost:"+constant.PORT+"/api/v1/%s", shortKey)

	w.WriteHeader(http.StatusCreated)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(shortenedURL)
}

// generateShortKey create short keys
func (urlCntx *UrlShortnerContext) generateShortKey() string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	const keyLength = 6

	rand.Seed(time.Now().UnixNano())
	shortKey := make([]byte, keyLength)
	for i := range shortKey {
		shortKey[i] = charset[rand.Intn(len(charset))]
	}
	return string(shortKey)
}

// Redirect Url handler
func (urlCntx *UrlShortnerContext) RedirectAPI(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	shortKey := strings.TrimSpace(vars["urlshortenerkey"])
	if shortKey == "" {
		http.Error(w, "Shortened key is missing", http.StatusBadRequest)
		return
	}

	// Retrieve the original URL from the `urls` map using the shortened key
	originalURL, found := urlCntx.urls[shortKey]
	if !found {
		http.Error(w, "Shortened key not found", http.StatusNotFound)
		return
	}

	// Redirect the user to the original URL
	http.Redirect(w, r, originalURL, http.StatusSeeOther)
}
