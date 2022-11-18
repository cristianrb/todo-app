package handlers

import (
	"net/http"
	"record-service/internal/services"
	"record-service/pkg/utils"
)

type handlerConfig struct {
	client       *http.Client
	itemsService services.ItemsService
}

func NewHandlerConfig(client *http.Client, itemsService services.ItemsService) *handlerConfig {
	return &handlerConfig{
		client:       client,
		itemsService: itemsService,
	}
}

func (u *handlerConfig) getTodoItems(w http.ResponseWriter, r *http.Request) {
	items, err := u.itemsService.GetAll()
	if err != nil {
		utils.ErrorJSON(w, err, http.StatusBadRequest)
		return
	}

	utils.WriteJSON(w, http.StatusOK, items)
}
