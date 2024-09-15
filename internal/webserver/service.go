package webserver

import (
	"fmt"
	"log"
	"net/http"

	manager "microinformer/internal/maanger"
	"microinformer/internal/settings"
)

type Service struct {
	Settings *settings.Service
	Manager  *manager.Service
}

type request struct {
	Info   []manager.Info `json:"info"`
	Reboot bool           `json:"reboot"`
}

type response struct {
	Info     *manager.Info    `json:"info,omitempty"`
	Settings *settings.Public `json:"settings,omitempty"`
}

func NewService(
	settingsService *settings.Service,
	managerService *manager.Service) *Service {
	return &Service{
		Settings: settingsService,
		Manager:  managerService,
	}
}

func (s Service) Start() {
	fmt.Println("Start server by :8081")
	http.HandleFunc("/api.json", s.api)
	http.HandleFunc("/update", s.update)
	http.HandleFunc("/settings.json", s.setting)
	http.Handle("/page/", http.StripPrefix("/page", http.FileServer(http.Dir("./frontend"))))
	log.Fatal(http.ListenAndServe(":8081", nil))
}
