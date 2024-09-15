package webserver

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"

	manager "microinformer/internal/maanger"
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
	s.Settings.SetReboot()
	w.WriteHeader(http.StatusOK)
}

func (s Service) addItem(w http.ResponseWriter, r *http.Request) {
	var data manager.Info
	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}
	err := s.Manager.AddItem(data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if !data.Once {
		s.Settings.SetReboot()
	}
	w.WriteHeader(http.StatusCreated)
	return
}

func (s Service) delItem(w http.ResponseWriter, r *http.Request) {
	var data manager.Info
	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}
	err := s.Manager.DelItem(data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	s.Settings.SetReboot()
	w.WriteHeader(http.StatusCreated)
	return
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
	s.Settings.NoReboot()
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
	s.Settings.NoReboot()
}

func (s Service) FileUploadHandler(w http.ResponseWriter, r *http.Request) {
	// file limit 10MB
	r.ParseMultipartForm(10 << 20)

	// get uploaded file
	file, handler, err := r.FormFile("uploadfile")
	if err != nil {
		http.Error(w, "Error retrieving the file", http.StatusBadRequest)
		fmt.Println("Error retrieving the file:", err)
		return
	}
	defer file.Close()

	// get file info
	fmt.Printf("Uploaded File: %+v\n", handler.Filename)
	fmt.Printf("File Size: %+v\n", handler.Size)
	fmt.Printf("MIME Header: %+v\n", handler.Header)

	// create file on the server
	dst, err := os.Create(filepath.Join("./uploads", handler.Filename)) //TODO: generate random name and return
	if err != nil {
		http.Error(w, "Error creating file on server", http.StatusInternalServerError)
		fmt.Println("Error creating file:", err)
		return
	}
	defer dst.Close()

	// copy uploaded file to new file
	if _, err := io.Copy(dst, file); err != nil {
		http.Error(w, "Error saving the file", http.StatusInternalServerError)
		fmt.Println("Error saving the file:", err)
		return
	}

	// Отправляем ответ клиенту
	fmt.Fprintf(w, "File uploaded successfully: %s\n", handler.Filename)
}
