package service

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
)

type Storer interface {
	Get(key string) (string, error)
	Set(key string, val string) error
	Delete(key string) error
}

type Service struct {
	listenAddr string
	store      Storer
	router     *Router
}

func NewService(listenAddr string, store Storer) *Service {
	return &Service{
		listenAddr: listenAddr,
		store:      store,
	}
}

func (s *Service) Open() error {
	// Initialize service's router
	router := NewRouter(s)
	router.InitRouter()
	s.router = router

	return http.ListenAndServe(s.listenAddr, s.router.muxRouter)
}

func (s *Service) HandleGet(w http.ResponseWriter, r *http.Request) {
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
	vars := mux.Vars(r)
	key := vars["key"]
	if err := s.store.Delete(key); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}
