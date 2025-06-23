package engine

import (
	"encoding/json"
	"io"
	"log"
	"net/http"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/shafnybuilds/car_management_sys/models"
	"github.com/shafnybuilds/car_management_sys/service"
)

type EngineHandler struct {
	service service.EngineServiceInterface
}

func NewEngineHandler(service service.EngineServiceInterface) *EngineHandler {
	return &EngineHandler{
		service: service,
	}
}

// methods
func (e *EngineHandler) GetEngineByID(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	vars := mux.Vars(r)
	id := vars["id"]

	resp, err := e.service.GetEngineById(ctx, id)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println(err)
		return
	}

	body, err := json.Marshal(resp)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println(err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	// write response body
	_, err = w.Write(body)
	if err != nil {
		log.Println("Error Writing Response: ", err)
	}
}

func (e *EngineHandler) CreteEngine(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	body, err := io.ReadAll(r.Body)
	if err != nil {
		log.Println("Error Reading Response Body: ", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	var engineReq models.EngineRequest
	err = json.Unmarshal(body, &engineReq)
	if err != nil {
		log.Println("Error while Unmarshalling Engine Request body: ", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// sending the req to service layer
	createdEngine, err := e.service.CreateEngine(ctx, &engineReq)
	if err != nil {
		log.Println("Error Creating Engine: ", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	responseBody, err := json.Marshal(createdEngine)
	if err != nil {
		log.Println("Error while marshalling: ", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

	// write teh response body
	_, _ = w.Write(responseBody)
}

func (e *EngineHandler) UpdateEngine(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	params := mux.Vars(r)
	id := params["id"]

	body, err := io.ReadAll(r.Body)
	if err != nil {
		log.Println("Error Reading Request Body: ", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	var engineReq models.EngineRequest
	err = json.Unmarshal(body, &engineReq)
	if err != nil {
		log.Println("Error while Unmarshalling Request body: ", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// send the req to service layer
	updatedEngine, err := e.service.UpdateEngine(ctx, id, &engineReq)
	if err != nil {
		log.Println("Error While Updating the Engine: ", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	resBody, err := json.Marshal(updatedEngine)
	if err != nil {
		log.Println("Error while Marshalling Body: ", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	// write the response body
	_, _ = w.Write(resBody)
}

func (e *EngineHandler) DeleteEngine(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	params := mux.Vars(r)
	id := params["id"]

	deletedEngine, err := e.service.DeleteEngine(ctx, id)
	if err != nil {
		log.Println("Error While Deleting the Engine: ", err)
		w.WriteHeader(http.StatusInternalServerError)
		response := map[string]string{"error": "Invalid ID or Engine not Found"}
		jsonResponse, _ := json.Marshal(response)

		// write the response body
		w.Write(jsonResponse)
		return
	}

	// check engine was deleted properly or not
	if deletedEngine.EngineID == uuid.Nil {
		w.WriteHeader(http.StatusNotFound)
		response := map[string]string{"error": "Engine Not Found"}
		jsonResponse, _ := json.Marshal(response)

		// write
		_, _ = w.Write(jsonResponse)
		return
	}

	// return the deleted engine details
	jsonResponse, err := json.Marshal(deletedEngine)
	if err != nil {
		log.Println("Error while marshalling deleted engine response: ", err)
		w.WriteHeader(http.StatusInternalServerError)
		response := map[string]string{"error": "Internal Server Error"}
		jsonResponse, _ := json.Marshal(response)
		// write
		_, _ = w.Write(jsonResponse)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	_, _ = w.Write(jsonResponse)
}
