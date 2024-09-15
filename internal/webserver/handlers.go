package webserver

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func (s Service) update(w http.ResponseWriter, r *http.Request) {
	var data response
	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}
	if data.Settings != nil {
		s.Settings.Set(data.Settings)
	}
	if data.Info != nil {
		s.Manager.AddItem(*data.Info)
		w.WriteHeader(http.StatusCreated)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func (s Service) setting(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	err := json.NewEncoder(w).Encode(s.Settings.Data)
	if err != nil {
		_, err := fmt.Fprintf(w, "Hi")
		if err != nil {
			fmt.Println(err)
		}
	}
}

func (s Service) api(w http.ResponseWriter, r *http.Request) {
	data := request{Info: s.Manager.GetInfo(), Reboot: s.Settings.Data.Reboot}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	err := json.NewEncoder(w).Encode(data)
	if err != nil {
		_, err := fmt.Fprintf(w, "Hi")
		if err != nil {
			fmt.Println(err)
		}
	}
	s.Settings.Taught()
}
