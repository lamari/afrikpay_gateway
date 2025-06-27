package handlers

import (
    "encoding/json"
    "net/http"

    "github.com/afrikpay/gateway/services/crud/internal/models"
    "github.com/afrikpay/gateway/services/crud/internal/services"
    "github.com/gorilla/mux"
)

type UserHandler struct {
    Service *services.UserService
}

func (h *UserHandler) Register(r *mux.Router) {
    r.HandleFunc("/users", h.create).Methods(http.MethodPost)
    r.HandleFunc("/users/{id}", h.get).Methods(http.MethodGet)
    r.HandleFunc("/users/{id}", h.update).Methods(http.MethodPut)
    r.HandleFunc("/users/{id}", h.delete).Methods(http.MethodDelete)
}

func (h *UserHandler) create(w http.ResponseWriter, r *http.Request) {
    var req models.User
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }
    if err := h.Service.Create(r.Context(), &req); err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }
    w.WriteHeader(http.StatusCreated)
    _ = json.NewEncoder(w).Encode(req)
}

func (h *UserHandler) get(w http.ResponseWriter, r *http.Request) {
    id := mux.Vars(r)["id"]
    user, err := h.Service.Get(r.Context(), id)
    if err != nil {
        http.Error(w, err.Error(), http.StatusNotFound)
        return
    }
    _ = json.NewEncoder(w).Encode(user)
}

func (h *UserHandler) update(w http.ResponseWriter, r *http.Request) {
    id := mux.Vars(r)["id"]
    var req models.User
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }
    req.ID = id
    if err := h.Service.Update(r.Context(), &req); err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }
    _ = json.NewEncoder(w).Encode(req)
}

func (h *UserHandler) delete(w http.ResponseWriter, r *http.Request) {
    id := mux.Vars(r)["id"]
    if err := h.Service.Delete(r.Context(), id); err != nil {
        http.Error(w, err.Error(), http.StatusNotFound)
        return
    }
    w.WriteHeader(http.StatusNoContent)
}
