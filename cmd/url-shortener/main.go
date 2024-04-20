package main

import (
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/rohit-bisht1993/URL-Shortener/internal/urlshortener"

	"github.com/gorilla/mux"
)

// Route ... Keep router name, method, pattern and handler func
type Route struct {
	Name        string
	Method      string
	Pattern     string
	HandlerFunc http.HandlerFunc
}

// Routes ...  List of Routers
type Routes []Route

func createRouter(routes Routes) *mux.Router {
	router := mux.NewRouter().StrictSlash(true)
	for _, route := range routes {
		router.Methods(route.Method).
			Path(route.Pattern).
			Name(route.Name).
			Handler(route.HandlerFunc)
	}
	return router
}

func startServer(router *mux.Router) {
	log.Fatal(http.ListenAndServe(":10000", router))
	fmt.Printf("Server Started localhost:10000")
}

func main() {

	urlShortenerCtx := urlshortener.NewUrlShortener()
	var routes = Routes{
		{
			"urlshortener",
			strings.ToUpper("POST"),
			"/api/v1/urlshortener",
			urlShortenerCtx.UrlShortenerAPI,
		},
	}

	router := createRouter(routes)
	startServer(router)
}