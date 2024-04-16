package main

import (
	"encoding/json"
	"github.com/go-lets-go/validate"
	"log/slog"
	"net/http"
)

type PersonRequest struct {
	Name    string         `json:"name" validate:"required,blank,min=3,max=10"`
	CPF     string         `json:"cpf" validate:"cpf"`
	Company CompanyRequest `json:"company"`
}

type CompanyRequest struct {
	CNPJ string `validate:"cnpj"`
}

func main() {
	http.HandleFunc("/v1/person", handlePostPerson)
	slog.Error("", slog.String("failed when trying to run server",
		http.ListenAndServe(":8080", nil).Error()))
}

func handlePostPerson(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "POST permitted", http.StatusMethodNotAllowed)
		return
	}

	var person PersonRequest
	err := json.NewDecoder(r.Body).Decode(&person)
	if err != nil {
		http.Error(w, "Error decoder JSON", http.StatusBadRequest)
		return
	}

	validator := validate.NewValidate()
	validations, err := validator.Struct(person)

	if err != nil {
		http.Error(w, "Error validation", http.StatusInternalServerError)
		return
	}

	if len(validations) >= 1 {
		sendJSONResponse(w, http.StatusBadRequest, validations)
	} else {
		w.WriteHeader(http.StatusCreated)
	}
}

func sendJSONResponse(w http.ResponseWriter, statusCode int, validations []validate.FieldValidation) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(validations)
}
