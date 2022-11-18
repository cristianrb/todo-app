package handlers

import (
	"collector-service/pkg/models"
	"collector-service/utils"
	"encoding/json"
	"fmt"
	"net/http"
)

type handlerConfig struct {
	client       *http.Client
	recordClient models.RecordClient
}

func NewHandlerConfig(client *http.Client, recordClient models.RecordClient) *handlerConfig {
	return &handlerConfig{
		client:       client,
		recordClient: recordClient,
	}
}

func (h *handlerConfig) getAllItems(w http.ResponseWriter, r *http.Request) {
	recordURL := fmt.Sprintf("%s://%s/items", h.recordClient.Protocol, h.recordClient.Hostname)
	fmt.Println(recordURL)
	req, err := http.NewRequest("GET", recordURL, nil)
	if err != nil {
		utils.ErrorJSON(w, err, http.StatusBadRequest)
		return
	}

	response, err := h.client.Do(req)
	if err != nil {
		utils.ErrorJSON(w, err)
		return
	}
	defer response.Body.Close()

	var todoItems []models.TodoItem
	err = json.NewDecoder(response.Body).Decode(&todoItems)
	if err != nil {
		utils.ErrorJSON(w, err)
		return
	}
	utils.WriteJSON(w, http.StatusOK, todoItems)
}
