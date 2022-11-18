package main

import (
	"collector-service/internal/handlers"
	"collector-service/pkg/models"
	"fmt"
	"net/http"
	"os"
)

const (
	PORT                   = "8082"
	RECORD_CLIENT_PROTOCOL = "RECORD_CLIENT_PROTOCOL"
	RECORD_CLIENT_HOSTNAME = "RECORD_CLIENT_HOSTNAME"
)

func main() {
	recordClient := models.RecordClient{
		Protocol: os.Getenv(RECORD_CLIENT_PROTOCOL),
		Hostname: os.Getenv(RECORD_CLIENT_HOSTNAME),
	}
	fmt.Printf("%v", recordClient)
	fmt.Println(recordClient.Protocol)
	fmt.Println(recordClient.Hostname)
	handlerConfig := handlers.NewHandlerConfig(&http.Client{}, recordClient)

	srv := &http.Server{
		Addr:    fmt.Sprintf(":%s", PORT),
		Handler: handlerConfig.Routes(),
	}
	err := srv.ListenAndServe()
	if err != nil {
		panic(err)
	}
}
