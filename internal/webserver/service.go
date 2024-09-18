package webserver

import (
	"fmt"
	"log"
	"net/http"
	"os"

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
	Items    []manager.Info   `json:"items,omitempty"`
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
	err := os.MkdirAll("./uploads", os.ModePerm)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Start server by :8081")
	http.HandleFunc("/api/items/data.json", s.api)
	http.HandleFunc("/api/items/list", s.list)
	http.HandleFunc("/api/items/add", s.addItem)
	http.HandleFunc("/api/items/del", s.delItem)
	http.HandleFunc("/api/items/update", s.update)
	http.HandleFunc("/api/files/upload", s.FileUploadHandler)
	http.HandleFunc("/api/settings/data.json", s.setting)
	http.HandleFunc("/api/page/reboot", s.pageReboot)
	http.Handle("/page/", http.StripPrefix("/page", http.FileServer(http.Dir("./frontend/panel"))))
	http.Handle("/admin/", http.StripPrefix("/admin", http.FileServer(http.Dir("./frontend/admin"))))
	http.Handle("/uploads/", http.StripPrefix("/uploads", http.FileServer(http.Dir("./uploads"))))
	log.Fatal(http.ListenAndServe(":8081", nil))
}
