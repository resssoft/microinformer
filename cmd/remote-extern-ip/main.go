package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

func main() {
	http.HandleFunc("/api/ip/", info)
	log.Fatal(http.ListenAndServe(":8189", nil))
}

type result struct {
	Ip string `json:"ip"`
}

func info(w http.ResponseWriter, r *http.Request) {

	var data result
	data.Ip = ReadUserIP(r)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	err := json.NewEncoder(w).Encode(data)
	if err != nil {
		_, err := fmt.Fprintf(w, "setting error: %v", err)
		if err != nil {
			fmt.Println(err)
		}
	}
}

func ReadUserIP(r *http.Request) string {
	IPAddress := r.Header.Get("X-Real-Ip")
	if IPAddress == "" {
		IPAddress = r.Header.Get("X-Forwarded-For")
	}
	if IPAddress == "" {
		IPAddress = r.RemoteAddr
	}
	return IPAddress
}
