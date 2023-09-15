package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"transform/excel/model"
	"transform/excel/pkg/longman"
)

func main() {
	http.HandleFunc("/convert", convertHandler)
	port := 8080
	fmt.Printf("Starting server on port %d...\n", port)
	http.ListenAndServe(fmt.Sprintf(":%d", port), nil)
}

func convertHandler(w http.ResponseWriter, req *http.Request) {
	// Parse the multipart form data to get the uploaded file
	// err := req.ParseMultipartForm(32 << 20) // 32MB is the default used by FormFile
	err := req.ParseMultipartForm(10 << 20) // 10 MB limit
	if err != nil {
		http.Error(w, "Error parsing form data", http.StatusBadRequest)
		return
	}

	var jsonData []model.Word
	var totalWords []model.Word
	fhs := req.MultipartForm.File["excelFile"] // "excelFile" is the name of the file input field in the HTML form
	for _, fh := range fhs {
		file, err := fh.Open()
		if err != nil {
			http.Error(w, "Error getting file from form data", http.StatusBadRequest)
			return
		}

		// defer file.Close()

		totalWord, err := longman.GetExcelData(fh.Filename, file)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		totalWords = append(totalWords, totalWord...)
	}

	jsonData = longman.MregeWords(totalWords)

	// Marshal jsonData to a JSON response
	jsonBytes, err := json.MarshalIndent(jsonData, "", "  ")
	if err != nil {
		http.Error(w, "Error marshaling JSON", http.StatusInternalServerError)
		return
	}

	// Set the response content type
	w.Header().Set("Content-Type", "application/json")

	// Write the JSON data to the response
	_, err = w.Write(jsonBytes)
	if err != nil {
		http.Error(w, "Error writing JSON response", http.StatusInternalServerError)
		return
	}
}
