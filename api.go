package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
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
	app := http.NewServeMux()
	app.HandleFunc("DELETE /account/{id}", HandleFunc(s.handleDeleteAccount))
	app.HandleFunc("GET /account", HandleFunc(s.handleGetAccount))
	app.HandleFunc("GET /account/{id}", HandleFunc(s.getAccountById))
	app.HandleFunc("POST /account", HandleFunc(s.handleCreateAccount))
	app.HandleFunc("PUT /account/{id}", HandleFunc(s.handleUpdateAccount))

	fmt.Println("JSON api running on th port :", s.listnerAddr)
	//http.ListenAndServe(s.listnerAddr, &router)

	if err := http.ListenAndServe("localhost:3000", app); err != nil {
		fmt.Println(err.Error())
	}
}

func (s *APIServer) handleGetAccount(w http.ResponseWriter, r *http.Request) error {

	result, err := s.store.GetAllAccounts()

	if err != nil {

		return err
	}
	return WriteJson(w, http.StatusOK, result)

}
func (s *APIServer) handleUpdateAccount(w http.ResponseWriter, r *http.Request) error {

	idstr := r.PathValue("id")
	id, err := strconv.Atoi(idstr)
	if err != nil {
		return err
	}
	accRequest := new(AccountUdpateRequest)
	if err := json.NewDecoder(r.Body).Decode(accRequest); err != nil {
		return err
	}
	account := Account{
		ID:        id,
		FirstName: accRequest.FirstName,
		LastName:  accRequest.LastName,
		Number:    accRequest.Number,
		Balance:   accRequest.Balance,
	}

	return s.store.UpdateAccount(&account)
}
func (s *APIServer) getAccountById(w http.ResponseWriter, r *http.Request) error {

	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		return err
	}
	account, err := s.store.GetAccountById(id)

	if err != nil {
		return err
	}
	return WriteJson(w, http.StatusOK, account)
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

	idstr := r.PathValue("id")
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
