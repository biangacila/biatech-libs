package domaindriven

import (
	"encoding/json"
	"github.com/biangacila/biatech-libs/utils"
	"github.com/gorilla/mux"
	"net/http"
)

type GenericController[T any] interface {
	Create(w http.ResponseWriter, r *http.Request)
	Get(w http.ResponseWriter, r *http.Request)
	Find(w http.ResponseWriter, r *http.Request)
	Update(w http.ResponseWriter, r *http.Request)
	Delete(w http.ResponseWriter, r *http.Request)
}

type GenericControllerImpl[T any] struct {
	service GenericService[T] // Generic service
}

// NewGenericController creates a new instance of GenericControllerImpl
func NewGenericController[T any](service GenericService[T]) GenericController[T] {
	return &GenericControllerImpl[T]{service: service}
}

// Create method to create a new record
func (c *GenericControllerImpl[T]) Create(w http.ResponseWriter, r *http.Request) {
	var record T
	err := json.NewDecoder(r.Body).Decode(&record)
	if err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	// Use the generic service to save the record
	err = c.service.Save("generic", record, record)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
}

// Get method to retrieve records
func (c *GenericControllerImpl[T]) Get(w http.ResponseWriter, r *http.Request) {
	var vars = mux.Vars(r)
	entity, _ := vars["entity"]

	// Example: fetch records using field values from query params (simple example)
	fieldValues, _ := utils.ExtractQueryParams(r)

	// Use the generic service to fetch the records
	records, err := c.service.Get(entity, fieldValues, *new(T))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(records)
}

// Get method to retrieve records
func (c *GenericControllerImpl[T]) Find(w http.ResponseWriter, r *http.Request) {
	var vars = mux.Vars(r)
	entity, _ := vars["entity"]

	// Example: fetch records using field values from query params (simple example)
	fieldValues, _ := utils.ExtractQueryParams(r) /*map[string]interface{}{
		"org": r.URL.Query().Get("org"),
	}*/

	// Use the generic service to fetch the records
	records, err := c.service.Get(entity, fieldValues, *new(T))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(records)
}

// Update method to update an existing record
func (c *GenericControllerImpl[T]) Update(w http.ResponseWriter, r *http.Request) {
	var vars = mux.Vars(r)
	entity, _ := vars["entity"]

	// Assuming we use query parameters for conditions and body for field values to update
	var fieldValues map[string]interface{}
	err := json.NewDecoder(r.Body).Decode(&fieldValues)
	if err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	conditions, _ := utils.ExtractQueryParams(r)

	// Use the generic service to update the record
	err = c.service.Update(entity, conditions, fieldValues, *new(T))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

// Delete method to delete a record
func (c *GenericControllerImpl[T]) Delete(w http.ResponseWriter, r *http.Request) {
	// Example: delete by org and id from query params
	fieldValues := map[string]interface{}{
		"org":  r.URL.Query().Get("org"),
		"code": r.URL.Query().Get("code"),
	}

	// Use the generic service to delete the record
	err := c.service.Delete("generic", fieldValues, *new(T))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}
