package urlshortener

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/gorilla/mux"
	"github.com/rohit-bisht1993/URL-Shortener/internal/constant"
	"github.com/rohit-bisht1993/URL-Shortener/internal/utils"
)

type urlInfo struct {
	Url string `json:"url,omitempty"`
}

type MetricInfo struct {
	Hostname string
	Count    int
}

type UrlShortnerContext struct {
	urls       map[string]string
	MetricData []MetricInfo
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

		originalURL = strings.TrimSpace(reqData.Url)

		if len(originalURL) == 0 {
			http.Error(w, "URL parameter is missing", http.StatusBadRequest)
			return
		}
	}

	parsedUrl, err := url.Parse(originalURL)
	if err != nil {
		http.Error(w, "URL parsing Error", http.StatusInternalServerError)
		return
	}

	found := false
	for idx, data := range urlCntx.MetricData {
		if strings.EqualFold(data.Hostname, parsedUrl.Hostname()) {
			urlCntx.MetricData[idx].Count++
			found = true
			break
		}
	}

	if !found {
		var data MetricInfo
		data.Hostname = parsedUrl.Hostname()
		data.Count++
		urlCntx.MetricData = append(urlCntx.MetricData, data)
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

// UrlMetricAPI Url handler
func (urlCntx *UrlShortnerContext) UrlMetricAPI(w http.ResponseWriter, r *http.Request) {
	if len(urlCntx.MetricData) == 0 {
		w.WriteHeader(http.StatusOK)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode("No Data Found")
		return
	}

	first, second, third := 0, 0, 0
	firstIdx, secondIdx, thirdIdx := 0, 0, 0
	for idx, data := range urlCntx.MetricData {
		if data.Count > first {
			thirdIdx = secondIdx
			secondIdx = firstIdx
			third = second
			second = first
			first = data.Count
			firstIdx = idx
			continue
		}
		if data.Count > second {
			thirdIdx = secondIdx
			secondIdx = idx
			third = second
			second = data.Count
			continue
		}
		if data.Count > third {
			thirdIdx = idx
			third = data.Count
			continue
		}
	}
	resp := []MetricInfo{}
	resp = append(resp, urlCntx.MetricData[firstIdx])

	if secondIdx != firstIdx {
		resp = append(resp, urlCntx.MetricData[secondIdx])
	}
	if secondIdx != thirdIdx && firstIdx != thirdIdx {
		resp = append(resp, urlCntx.MetricData[thirdIdx])
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}
