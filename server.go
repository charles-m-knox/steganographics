package main

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"image"
	"image/png"
	"log"
	"net/http"
	"time"
)

type HideRequest struct {
	Msg string `json:"msg"`
	Png string `json:"png"`
}

type ExtractRequest struct {
	Png string `json:"png"`
}

type ExtractResponse struct {
	Msg string `json:"msg"`
}

func server(address string) {
	http.HandleFunc("/api/hide", hideTextInImageHandler)
	http.HandleFunc("/api/extract", extractTextFromImageHandler)

	log.Printf("server listening on %v", address)

	srv := &http.Server{
		ReadTimeout:  2 * time.Second,
		WriteTimeout: 3 * time.Second,
		Addr:         address,
	}

	var err error

	if flagCert == "" || flagKey == "" {
		log.Printf("listening on %v", srv.Addr)
		err = srv.ListenAndServe()
	} else {
		log.Printf("listening on %v with TLS", srv.Addr)
		err = srv.ListenAndServeTLS(flagCert, flagKey)
	}

	if err != nil {
		log.Fatalf("Error starting server: %s\n", err)
	}
}

func hideTextInImageHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		fmt.Fprintf(w, "Invalid method")

		return
	}

	var hideReq HideRequest

	err := json.NewDecoder(r.Body).Decode(&hideReq)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "Invalid request")

		return
	}

	decodedPng, err := base64.StdEncoding.DecodeString(hideReq.Png)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "Invalid base64 PNG data")

		return
	}

	img, _, err := image.Decode(bytes.NewReader(decodedPng))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "Invalid PNG data")

		return
	}

	imgWithText, err := HideTextInImage(img, []byte(hideReq.Msg))
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)

		log.Printf("failed to encode text in img: %v", err.Error())

		return
	}

	buffer := new(bytes.Buffer)

	err = png.Encode(buffer, imgWithText)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Error encoding the image")

		return
	}

	w.Header().Set("Content-Type", "image/png")

	_, err = w.Write(buffer.Bytes())
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Printf("failed to return response for hiding text in image: %v", err.Error())

		return
	}
}

func extractTextFromImageHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		fmt.Fprintf(w, "Invalid method")

		return
	}

	var extractReq ExtractRequest

	err := json.NewDecoder(r.Body).Decode(&extractReq)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "Invalid request")

		return
	}

	decodedPng, err := base64.StdEncoding.DecodeString(extractReq.Png)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "Invalid base64 PNG data")

		return
	}

	img, _, err := image.Decode(bytes.NewReader(decodedPng))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "Invalid PNG data")

		return
	}

	data := ExtractTextFromImage(img)

	response := ExtractResponse{
		Msg: string(data),
	}

	w.Header().Set("Content-Type", "application/json")

	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Printf("failed to return response for extracting text from image: %v", err.Error())

		return
	}
}
