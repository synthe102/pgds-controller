package router

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/google/uuid"
	"github.com/synthe102/pgds-controller/internal/model"
)

var items = make(map[string]model.Item)
var newItems = make(chan model.Item, 10)

type Router struct {
	*chi.Mux
}

func New() Router {
	r := chi.NewRouter()

	r.Use(middleware.Logger)

	r.Route("/items", func(r chi.Router) {
		r.Get("/", GetItems)
		r.Post("/", AddItem)
		r.Route("/{id}", func(r chi.Router) {
			r.Get("/", GetItem)
			r.Put("/", UpdateItem)
			r.Delete("/", DeleteItem)
		})
		r.Get("/watch", WatchItems)
	})

	return Router{r}
}

func (r Router) Run() error {
	return http.ListenAndServe(":3000", r.Mux)
}

func WatchItems(w http.ResponseWriter, r *http.Request) {
	select {
	case <-time.After(5 * time.Second):
		http.Error(w, "{}", http.StatusRequestTimeout)
	case <-r.Context().Done():
		http.Error(w, "{}", http.StatusRequestTimeout)
	case item := <-newItems:
		json.NewEncoder(w).Encode(item)
	}
}

// GetItems returns all items
func GetItems(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode(items)
}

// AddItem adds a new item
func AddItem(w http.ResponseWriter, r *http.Request) {
	var item model.Item
	json.NewDecoder(r.Body).Decode(&item)
	uuid, err := uuid.NewV7()
	if err != nil {
		http.Error(w, "Error generating UUID", http.StatusInternalServerError)
		return
	}
	item.ID = uuid.String()
	items[item.ID] = item
	newItems <- item
	json.NewEncoder(w).Encode(item)
}

// GetItem returns a single item
func GetItem(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	item, ok := items[id]
	if !ok {
		http.Error(w, "Item not found", http.StatusNotFound)
		return
	}
	json.NewEncoder(w).Encode(item)
}

// UpdateItem updates an existing item
func UpdateItem(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	var item model.Item
	json.NewDecoder(r.Body).Decode(&item)
	items[id] = item
	json.NewEncoder(w).Encode(item)
}

// DeleteItem deletes an item
func DeleteItem(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	delete(items, id)
	w.WriteHeader(http.StatusNoContent)
}
