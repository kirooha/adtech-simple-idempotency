package main

import (
	"encoding/json"
	"log"
	"net/http"
	"sync"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

type CreateCampaignRequest struct {
	ID          uuid.UUID `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
}

type CreateCampaignResponse struct {
	ID uuid.UUID `json:"id"`
}

type AddViewRequest struct {
	EventID uuid.UUID `json:"event_id"`
	AdID    uuid.UUID `json:"ad_id"`
}

type ViewResponse struct {
	AdID  uuid.UUID `json:"ad_id"`
	Views int       `json:"views"`
}

var (
	campaignMap sync.Map
	viewMap     sync.Map

	m  = make(map[uuid.UUID]int)
	mu sync.RWMutex
)

func main() {
	r := mux.NewRouter()

	r.HandleFunc("/create-campaign", func(writer http.ResponseWriter, request *http.Request) {
		var (
			req  CreateCampaignRequest
			resp CreateCampaignResponse
		)

		if err := json.NewDecoder(request.Body).Decode(&req); err != nil {
			writer.WriteHeader(http.StatusInternalServerError)
			return
		}

		id, _ := campaignMap.LoadOrStore(req.Name, uuid.New())

		resp.ID = id.(uuid.UUID)
		req.ID = resp.ID

		log.Printf("received request to create campaign, request: %+v", req)

		writer.Header().Add("Content-Type", "application/json")
		err := json.NewEncoder(writer).Encode(resp)
		if err != nil {
			writer.WriteHeader(http.StatusInternalServerError)
			return
		}
	}).Methods(http.MethodPost)

	r.HandleFunc("/add-view", func(writer http.ResponseWriter, request *http.Request) {
		var (
			req AddViewRequest
		)

		if err := json.NewDecoder(request.Body).Decode(&req); err != nil {
			writer.WriteHeader(http.StatusInternalServerError)
			return
		}

		_, loaded := viewMap.LoadOrStore(req.EventID, 1)
		if loaded {
			return
		}

		mu.Lock()
		m[req.AdID]++
		mu.Unlock()
	}).Methods(http.MethodPost)

	r.HandleFunc("/view/{ad_id}", func(writer http.ResponseWriter, request *http.Request) {
		var (
			adIDStr = mux.Vars(request)["ad_id"]
			resp    ViewResponse
		)

		adID, err := uuid.Parse(adIDStr)
		if err != nil {
			writer.WriteHeader(http.StatusInternalServerError)
			return
		}

		mu.RLock()
		views := m[adID]
		mu.RUnlock()

		resp.AdID = adID
		resp.Views = views

		writer.Header().Add("Content-Type", "application/json")
		err = json.NewEncoder(writer).Encode(resp)
		if err != nil {
			writer.WriteHeader(http.StatusInternalServerError)
			return
		}
	}).Methods(http.MethodGet)

	log.Println("adserver is running on port 9091")
	log.Fatal(http.ListenAndServe(":9091", r))
}
