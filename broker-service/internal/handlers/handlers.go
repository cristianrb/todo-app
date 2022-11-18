package handlers

import (
	"broker-service/internal/events"
	"broker-service/pkg/models"
	"broker-service/pkg/utils"
	"net/http"
)

type handlerConfig struct {
	client  *http.Client
	emitter events.Emitter
}

func NewHandlerConfig(client *http.Client, emitter events.Emitter) *handlerConfig {
	return &handlerConfig{
		client:  client,
		emitter: emitter,
	}
}

func (u *handlerConfig) pushItemToQueue(w http.ResponseWriter, r *http.Request) {
	todoItem := models.TodoItem{}

	err := utils.ReadJSON(w, r, &todoItem)
	if err != nil {
		utils.ErrorJSON(w, err, http.StatusBadRequest)
		return
	}

	if err := u.emitter.PushToQueue(todoItem); err != nil {
		utils.ErrorJSON(w, err, http.StatusBadRequest)
		return
	}
	utils.WriteJSON(w, http.StatusAccepted, nil)
}
