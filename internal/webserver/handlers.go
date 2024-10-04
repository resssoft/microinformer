package webserver

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"

	manager "microinformer/internal/maanger"
	gen "microinformer/pkg/generator"
)

func (s Service) list(w http.ResponseWriter, r *http.Request) {
	s.logRequest(r)
	data := s.Manager.ListItem()
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	err := json.NewEncoder(w).Encode(data)
	if err != nil {
		_, err := fmt.Fprintf(w, "list error: %v", err)
		if err != nil {
			fmt.Println(err)
		}
	}
}

func (s Service) update(w http.ResponseWriter, r *http.Request) {
	s.logRequest(r)
	var data response
	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
		http.Error(w, "Invalid input"+err.Error(), http.StatusBadRequest)
		return
	}
	if data.Settings != nil {
		s.Settings.Set(data.Settings)
	}
	if len(data.Items) > 0 {
		_ = s.Manager.Update(data.Items)
		w.WriteHeader(http.StatusCreated)
		return
	}
	s.Settings.SetReboot()
	w.WriteHeader(http.StatusOK)
}

func (s Service) importItems(w http.ResponseWriter, r *http.Request) {
	s.logRequest(r)
	var data manager.ImportItems
	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
		http.Error(w, "Invalid input"+err.Error(), http.StatusBadRequest)
		return
	}
	dataJson, errJson := json.MarshalIndent(data, "", "  ")
	fmt.Println(string(dataJson), errJson) // debug
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	result := s.Manager.AddItems(data)
	err := json.NewEncoder(w).Encode(result)
	if err != nil {
		_, err := fmt.Fprintf(w, "import error: %v", err)
		if err != nil {
			fmt.Println(err)
		}
	}
}

func (s Service) delItem(w http.ResponseWriter, r *http.Request) {
	s.logRequest(r)
	var data manager.Info
	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
		http.Error(w, "Invalid input"+err.Error(), http.StatusBadRequest)
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
	s.logRequest(r)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	err := json.NewEncoder(w).Encode(s.Settings.Data)
	if err != nil {
		_, err := fmt.Fprintf(w, "setting error: %v", err)
		if err != nil {
			fmt.Println(err)
		}
	}
	s.Settings.NoReboot()
}

func (s Service) api(w http.ResponseWriter, r *http.Request) {
	s.logRequest(r)
	data := request{Info: s.Manager.GetInfo(), Reboot: s.Settings.Data.Reboot}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	err := json.NewEncoder(w).Encode(data)
	if err != nil {
		_, err := fmt.Fprintf(w, "api data error: %v", err)
		if err != nil {
			fmt.Println(err)
		}
	}
	s.Settings.NoReboot()
}

func (s Service) pageReboot(w http.ResponseWriter, r *http.Request) {
	s.logRequest(r)
	s.Settings.SetReboot()
	w.WriteHeader(http.StatusOK)
}

func (s Service) FileUploadHandler(w http.ResponseWriter, r *http.Request) {
	s.logRequest(r)
	// file limit 10MB
	err := r.ParseMultipartForm(10 << 20)
	if err != nil {
		fmt.Println("Error ParseMultipartForm:", err)
	}

	// get uploaded file
	newFileName := r.PostFormValue("filename")
	file, handler, err := r.FormFile("uploadfile")
	if err != nil {
		http.Error(w, "Error retrieving the file: "+err.Error(), http.StatusBadRequest)
		fmt.Println("Error retrieving the file:", err)
		return
	}
	defer file.Close()

	// get file info
	fmt.Printf("Uploaded File: %+v\n", handler.Filename)
	fmt.Printf("File Size: %+v\n", handler.Size)
	fmt.Printf("MIME Header: %+v\n", handler.Header)

	// create file on the server
	if newFileName == "" {
		newFileName = gen.LatinStr(10) + filepath.Ext(handler.Filename)
	}
	fmt.Printf("new File name: %+v\n", newFileName)
	dst, err := os.Create(filepath.Join("./uploads", newFileName))
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
	fmt.Fprintf(w, "%s", newFileName)
}

func (s Service) logRequest(r *http.Request) {
	fmt.Println(r.URL)
}
