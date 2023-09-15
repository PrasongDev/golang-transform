package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"transform/excel/model"
	"transform/excel/pkg/longman"

	"github.com/xuri/excelize/v2"
)

func main() {
	http.HandleFunc("/convert", convertHandler)
	port := 8080
	fmt.Printf("Starting server on port %d...\n", port)
	http.ListenAndServe(fmt.Sprintf(":%d", port), nil)
}

func convertHandler(w http.ResponseWriter, r *http.Request) {
	// Parse the multipart form data to get the uploaded file
	err := r.ParseMultipartForm(10 << 20) // 10 MB limit
	if err != nil {
		http.Error(w, "Error parsing form data", http.StatusBadRequest)
		return
	}

	file, _, err := r.FormFile("excelFile") // "excelFile" is the name of the file input field in the HTML form
	if err != nil {
		http.Error(w, "Error getting file from form data", http.StatusBadRequest)
		return
	}
	defer file.Close()

	// Create a new XLSX file from the uploaded file
	xlsxFile, err := excelize.OpenReader(file)
	if err != nil {
		http.Error(w, "Error opening Excel file", http.StatusInternalServerError)
		return
	}

	// Process the Excel file (similar to the previous example)
	// Define your struct, extract data, and append it to jsonData

	// For example, you can define a struct and append data as shown previously

	// Define a struct to hold your data (assuming Excel has headers)
	// type Data struct {
	// 	Header1 string `json:"header1"`
	// 	Header2 string `json:"header2"`
	// 	Header3 string `json:"header3"`
	// 	// Add more fields as needed based on your Excel columns
	// }

	var jsonData []model.Word

	// Iterate through the sheets in the XLSX file
	for _, sheetName := range xlsxFile.GetSheetList() {
		rows, err := xlsxFile.GetRows(sheetName)
		if err != nil {
			http.Error(w, "Error reading sheet", http.StatusInternalServerError)
			return
		}

		// Assuming the first row contains headers
		// var headers []string
		// if len(rows) > 0 {
		// 	headers = rows[0]
		// }

		jsonData = longman.LongmanCommunication9000(rows)
	}

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
