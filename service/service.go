// The Service module provides RESTFul interfaces for users to operate the key-value storage
// system, and the data consensus is guaranteed by the Store module, which is based on the
// Raft Consensus Algorithm.
package service

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

// An interface for Storer to implement.
type Storer interface {
	Get(key string) (string, error)
	Set(key string, val string) error
	Delete(key string) error
	Join(nodeId string, nodeAddr string) error
}

// Key-value storage services provider.
type Service struct {
	listenAddr string  // Service bind address.
	store      Storer  // Storage system that can guarantee the consensus.
	router     *Router // Service router.
}

// Create a new service provider without router initialization.
func NewService(listenAddr string, store Storer) *Service {
	return &Service{
		listenAddr: listenAddr,
		store:      store,
	}
}

// Start providing service.
// Before starting, the Store module must be opened at first.
func (s *Service) Open() error {
	// Initialize service router
	router := NewRouter(s)
	router.InitRouter()
	s.router = router

	go func() {
		log.Fatal(http.ListenAndServe(s.listenAddr, s.router.muxRouter))
	}()

	return nil
}

func (s *Service) HandleGet(w http.ResponseWriter, r *http.Request) {
	log.Printf("%s\t%s", r.Method, r.RequestURI)
	vars := mux.Vars(r)
	key := vars["key"]
	val, err := s.store.Get(key)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(val)
}

func (s *Service) HandleSet(w http.ResponseWriter, r *http.Request) {
	log.Printf("%s\t%s", r.Method, r.RequestURI)
	pairs := map[string]string{}
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&pairs); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	defer r.Body.Close()

	for k, v := range pairs {
		if err := s.store.Set(k, v); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}

	w.WriteHeader(http.StatusOK)
}

func (s *Service) HandleDelete(w http.ResponseWriter, r *http.Request) {
	log.Printf("%s\t%s", r.Method, r.RequestURI)
	vars := mux.Vars(r)
	key := vars["key"]
	if err := s.store.Delete(key); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func (s *Service) HandleJoin(w http.ResponseWriter, r *http.Request) {
	log.Printf("%s\t%s", r.Method, r.RequestURI)
	nodes := map[string]string{}
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&nodes); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	defer r.Body.Close()

	for id, addr := range nodes {
		if err := s.store.Join(id, addr); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}

	w.WriteHeader(http.StatusOK)
}
