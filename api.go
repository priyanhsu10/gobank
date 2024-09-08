package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

type APIServer struct {
	listnerAddr string
	store       Storage
}

func WriteJson(w http.ResponseWriter, status int, v any) error {

	w.WriteHeader(status)
	w.Header().Add("Content-Type", "application/json")
	return json.NewEncoder(w).Encode(v)
}

type apiFunc func(http.ResponseWriter, *http.Request) error
type ApiError struct {
	Error string
}

func HandleFunc(f apiFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := f(w, r); err != nil {
			//hendle error
			WriteJson(w, http.StatusBadRequest, ApiError{Error: err.Error()})
		}
	}
}

func NewAPIServer(listnerAddres string, storage Storage) *APIServer {
	return &APIServer{
		listnerAddr: listnerAddres,
		store:       storage,
	}
}
func (s *APIServer) run() {
	router := mux.Router{}
	router.HandleFunc("/account", HandleFunc(s.hanldleAccount))
	router.HandleFunc("/account/{id}", HandleFunc(s.hanldleAccount))
	fmt.Println("JSON api running on th port :", s.listnerAddr)
	http.ListenAndServe(s.listnerAddr, &router)
}
func (s *APIServer) hanldleAccount(w http.ResponseWriter, r *http.Request) error {

	if r.Method == "GET" {
		return s.handleGetAccount(w, r)
	}
	if r.Method == "POST" {
		return s.handleCreateAccount(w, r)
	}
	if r.Method == "DELETE" {
		return s.handleDeleteAccount(w, r)
	}
	return fmt.Errorf("Method not suppored %s", r.Method)
}

func (s *APIServer) handleGetAccount(w http.ResponseWriter, r *http.Request) error {

	result, err := s.store.GetAllAccounts()

	if err != nil {

		return err
	}
	return WriteJson(w, http.StatusOK, result)

}
func (s *APIServer) handleCreateAccount(w http.ResponseWriter, r *http.Request) error {
	accRequest := new(AccountCreateReuqust)
	if err := json.NewDecoder(r.Body).Decode(accRequest); err != nil {
		return err
	}

	account := NewAccount(accRequest.FirstName, accRequest.LastName)
	if err := s.store.CreateAccount(account); err != nil {
		return err
	}

	return WriteJson(w, http.StatusOK, account)
}
func (s *APIServer) handleDeleteAccount(w http.ResponseWriter, r *http.Request) error {
	vars := mux.Vars(r)
	idstr := vars["id"]
	id, err := strconv.Atoi(idstr)
	if err != nil {
		return fmt.Errorf("invalid id given %s", idstr)
	}

	err = s.store.DeleteAccont(id)
	return err

}
func (s *APIServer) handleTransfer(w http.ResponseWriter, r *http.Request) error {
	return nil
}
