package main

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"image"
	"image/color"
	"image/png"
	"log"
	"net/http"
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

func server(address string, port int) {
	http.HandleFunc("/api/hide", hideTextInImageHandler)
	http.HandleFunc("/api/extract", extractTextFromImageHandler)
	http.HandleFunc("/", serveIndexFile)
	http.HandleFunc("/index.html", serveIndexFile)

	addr := fmt.Sprintf("%v:%v", address, port)

	log.Printf("Server listening on %v", addr)
	err := http.ListenAndServe(addr, nil)
	if err != nil {
		log.Fatalf("Error starting server: %s\n", err)
	}
}

func serveIndexFile(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprintf(w, "Page not found")
		return
	}

	http.ServeFile(w, r, "assets/index.html")
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

	hiddenText := hideReq.Msg

	bounds := img.Bounds()
	width, height := bounds.Max.X, bounds.Max.Y

	data := append([]byte(hiddenText), 0)
	if len(data)*8 > width*height {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "Input image is too small to hide the text")
		return
	}

	imgWithText := image.NewRGBA(bounds)
	counter := 0

	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			c := img.At(x, y)
			rgba := color.RGBAModel.Convert(c).(color.RGBA)

			if counter < len(data)*8 {
				b := data[counter/8]
				bit := (b >> uint(counter%8)) & 1
				rgba.R = (rgba.R & 0xFE) | bit
				counter++
			}

			imgWithText.Set(x, y, rgba)
		}
	}

	buffer := new(bytes.Buffer)
	err = png.Encode(buffer, imgWithText)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Error encoding the image")
		return
	}

	w.Header().Set("Content-Type", "image/png")
	w.Write(buffer.Bytes())
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

	bounds := img.Bounds()
	width, height := bounds.Max.X, bounds.Max.Y

	var data bytes.Buffer
	var currentByte byte
	bitIndex := 0

	for counter := 0; counter < width*height; counter++ {
		x, y := counter%width, counter/width

		c := color.RGBAModel.Convert(img.At(x, y)).(color.RGBA)
		bit := c.R & 1

		currentByte |= (bit << uint(bitIndex))
		bitIndex++

		if bitIndex == 8 {
			if currentByte == 0 {
				break
			}
			data.WriteByte(currentByte)
			currentByte = 0
			bitIndex = 0
		}
	}

	response := ExtractResponse{
		Msg: data.String(),
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
