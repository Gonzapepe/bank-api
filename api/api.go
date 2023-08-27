package api

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/Gonzapepe/bank-api/helper"
	"github.com/Gonzapepe/bank-api/middleware"
	"github.com/Gonzapepe/bank-api/storage"
	"github.com/Gonzapepe/bank-api/types"
	"github.com/gorilla/mux"
)

func WriteJSON(w http.ResponseWriter, status int, v any) error {
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(status)
	return json.NewEncoder(w).Encode(v)
}

type apiFunc func(http.ResponseWriter,*http.Request) error

type ApiError struct {
	Error string
}

func makeHTTPHandlerFunc(f apiFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := f(w, r); err != nil {
			WriteJSON(w, http.StatusBadRequest, ApiError{Error: err.Error()})
		}
	}
}

type APIServer struct {
	listenAddr string
	store storage.Storage
}

func NewAPIServer(listenAddr string, store storage.Storage) *APIServer {
	return &APIServer{
		listenAddr: listenAddr,
		store: store,
	}
}

func (s *APIServer) Run() {
	router := mux.NewRouter()

	router.HandleFunc("/account", makeHTTPHandlerFunc(s.handleCreateAccount)).Methods("POST")
	router.HandleFunc("/account/login", makeHTTPHandlerFunc(s.handleLogin)).Methods("POST")
	router.HandleFunc("/accounts", makeHTTPHandlerFunc(s.handleGetAccounts)).Methods("GET")
	router.HandleFunc("/account/{id}", makeHTTPHandlerFunc(s.handleGetAccount)).Methods("GET")
	router.HandleFunc("/account/{id}", makeHTTPHandlerFunc(s.handleDeleteAccount)).Methods("DELETE")
	router.HandleFunc("/account/{id}", makeHTTPHandlerFunc(s.handleUpdateAccount)).Methods("PUT")
	

	log.Println("JSON Api server running on port: ", s.listenAddr)

	http.ListenAndServe(s.listenAddr, router)
}

func (s *APIServer) handleGetAccount(w http.ResponseWriter, r *http.Request) error {
	id, err := strconv.Atoi(mux.Vars(r)["id"])
	
	if err != nil {
		return err
	}
	account, err := s.store.GetAccountByID(id)
	if err != nil {
		return err
	}

	return WriteJSON(w, http.StatusOK, account)
}

func (s *APIServer) handleCreateAccount(w http.ResponseWriter, r *http.Request) error {

	createAccountReq := &types.CreateAccountRequest{}

	if err := json.NewDecoder(r.Body).Decode(&createAccountReq); err != nil {
		return err
	}

	hashedPassword, err := helper.EncryptPassword(createAccountReq.Password)

	if err != nil {
		fmt.Println("There's an error in password encryption: ", err)
		return err
	}

	account := types.NewAccount(createAccountReq.FirstName, createAccountReq.LastName, createAccountReq.Gender, createAccountReq.Dni, hashedPassword)

	if err := s.store.CreateAccount(account); err != nil {
		return err
	}

	return WriteJSON(w, http.StatusOK, account)
}

func (s *APIServer) handleGetAccounts(w http.ResponseWriter, r *http.Request) error {

	accounts, err := s.store.GetAccounts()
	if err != nil {
		return err
	}

	return WriteJSON(w, http.StatusOK, accounts)
}

func (s *APIServer) handleDeleteAccount(w http.ResponseWriter, r *http.Request) error {

	id, err := strconv.Atoi(mux.Vars(r)["id"])

	if err != nil {
		return err
	}

	err = s.store.DeleteAccount(id)
	if err != nil {
		return err
	}

	response := fmt.Sprintf("Account with id %d deleted successfully", id)
	jsonData := map[string]string{"response": response}

	return WriteJSON(w, http.StatusOK, jsonData)
}

func (s *APIServer) handleUpdateAccount(w http.ResponseWriter, r *http.Request) error {

	id, err := strconv.Atoi(mux.Vars(r)["id"])
	if err != nil {
		return err
	}
	
	updatedAccount := &types.Account{}

	if err := json.NewDecoder(r.Body).Decode(&updatedAccount); err != nil {
		return err
	}


	err = s.store.UpdateAccount(id, updatedAccount)
	
	if err != nil {
		return err
	}

	response := fmt.Sprintf("Account with id %d updated successfully", id)
	jsonData := map[string]string{"response": response}

	return WriteJSON(w, http.StatusOK, jsonData)
}

func (s *APIServer) handleTransfer(w http.ResponseWriter, r *http.Request) error {
	return nil
}

func (s *APIServer) handleLogin(w http.ResponseWriter, r *http.Request) error {
	login := &types.Login{}

	if err := json.NewDecoder(r.Body).Decode(&login); err != nil {
		return err
	}

	account, err := s.store.GetAccountByDni(login.Dni)

	if err != nil {
		return err
	}

	match := helper.DecryptPassword(login.Password, account.Password)

	if match != true {
		return err
	}

	jwt, err := middleware.GenerateJWT(account.Dni)

	if err != nil {
		return err
	}

	jsonData := map[string]string{"response": jwt}

	return WriteJSON(w, http.StatusOK, jsonData)
}