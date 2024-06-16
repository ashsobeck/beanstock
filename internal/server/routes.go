package server

import (
	"crypto/sha256"
	"encoding/json"
	"io"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/google/uuid"

	"beanstock/internal/parser"
	"beanstock/internal/types"
)

func (s *Server) RegisterRoutes() http.Handler {
	r := chi.NewRouter()
	r.Use(middleware.Logger)

	r.Get("/", s.HelloWorldHandler)

	r.Get("/health", s.healthHandler)
	r.Post("/register", s.RegisterSiteHandler)

	return r
}

func (s *Server) HelloWorldHandler(w http.ResponseWriter, r *http.Request) {
	resp := make(map[string]string)
	resp["message"] = "Hello World"

	jsonResp, err := json.Marshal(resp)
	if err != nil {
		log.Fatalf("error handling JSON marshal. Err: %v", err)
	}

	_, _ = w.Write(jsonResp)
}

func (s *Server) healthHandler(w http.ResponseWriter, r *http.Request) {
	jsonResp, _ := json.Marshal(s.db.Health())
	_, _ = w.Write(jsonResp)
}

func (s *Server) RegisterSiteHandler(w http.ResponseWriter, r *http.Request) {

	var site types.Website
	err := json.NewDecoder(r.Body).Decode(&site)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	site.Id = uuid.NewString()
	site.Json, err = GetSiteJson(&site)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	err = s.db.StoreWebsite(site)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func GetSiteJson(s *types.Website) (m map[string]interface{}, e error) {
	siteResp, getErr := http.Get(s.Url)
	if getErr != nil {
		return nil, getErr
	}
	defer siteResp.Body.Close()
	body, readErr := io.ReadAll(siteResp.Body)
	if readErr != nil {
		return nil, readErr
	}

	sortedJson, err := parser.SortJson(body)
	if err != nil {
		return nil, err
	}
	s.LastHash, err = sha256.New().Write(sortedJson)
	if err != nil {
		return nil, err
	}

	e = json.Unmarshal(sortedJson, &m)
	return m, e
}
