package longman

import (
	"fmt"
	"log"
	"mime/multipart"
	"transform/excel/model"
	"transform/excel/util"

	"github.com/pkg/errors"
	"github.com/xuri/excelize/v2"
)

func GetExcelData(fileName string, fileMeta multipart.File) ([]model.Word, error) {
	// Create a new XLSX file from the uploaded file
	xlsxFile, err := excelize.OpenReader(fileMeta)
	if err != nil {
		return nil, errors.Wrap(err, "Error opening Excel file")
	}

	var jsonData []model.Word

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

	// Iterate through the sheets in the XLSX file
	for _, sheetName := range xlsxFile.GetSheetList() {
		rows, err := xlsxFile.GetRows(sheetName)
		if err != nil {
			return nil, errors.Wrap(err, "Error reading sheet")
		}

		// Assuming the first row contains headers
		// var headers []string
		// if len(rows) > 0 {
		// 	headers = rows[0]
		// }

		switch fileName {
		case "longman3135.xlsx":
			jsonData = LongmanCommunication3135(rows)
		case "longman9000.xlsx":
			jsonData = LongmanCommunication9000(rows)
		}
	}

	return jsonData, nil
}

func MregeWords(words []model.Word) []model.Word {
	// Create a map to track unique Vocabulary values
	seen := make(map[string]model.Word)

	// Initialize a slice to store duplicates
	duplicates := make([]model.Word, 0)

	// Iterate over the slice and check for duplicates
	for _, word := range words {
		doingWords, found := seen[word.Vocabulary]

		// If the Vocabulary value is already in the map, it's a duplicate
		if !found {
			seen[word.Vocabulary] = word
		} else {
			if word.Pronunciation == "" {
				word.Pronunciation = doingWords.Pronunciation
			}
			if word.Phonetics == "" {
				word.Phonetics = doingWords.Phonetics
			}

			word.PartsOfSpeech = append(word.PartsOfSpeech, doingWords.PartsOfSpeech...)
			word.Tag = removeDuplicates(append(word.Tag, doingWords.Tag...))
			word.Syn = removeDuplicates(append(word.Syn, doingWords.Syn...))
			word.Source = removeDuplicates(append(word.Source, doingWords.Source...))
			word.Example = removeDuplicates(append(word.Example, doingWords.Example...))

			seen[word.Vocabulary] = word
		}
	}

	duplicates = WordMapToSlice(seen)

	return duplicates
}

func WordMapToSlice(wordMap map[string]model.Word) []model.Word {
	words := make([]model.Word, 0, len(wordMap))

	for _, word := range wordMap {
		words = append(words, word)
	}

	return words
}

func removeDuplicates(slice []string) []string {
	seen := make(map[string]bool)
	result := []string{}
	for _, value := range slice {
		if !seen[value] {
			seen[value] = true
			result = append(result, value)
		}
	}
	return result
}

func LongmanCommunication3135(rows [][]string) []model.Word {
	countID := 0
	countErr := 0
	arrErr := []string{}
	var words []model.Word

	// row >> 1: "NO.", 2: "Vocabulary", 3: "TranslationTH"
	for _, row := range rows {
		countID++

		var dataRow model.Word
		dataRow.ID = countID

		if (row[1] != "") && (row[2] != "") {
			words = append(words, model.Word{
				ID:            countID,
				Vocabulary:    row[1],
				Pronunciation: "",
				Phonetics:     "",
				PartsOfSpeech: []model.PartsOfSpeech{{
					Type:          "",
					TranslationTH: row[2],
				}},
				Tag:     []string{},
				Syn:     []string{},
				Source:  []string{"Longman Communication 3135"},
				Example: []string{},
			})
		} else {
			countErr++
			fmt.Println("--------------------------------------------")
			fmt.Printf(">> row >> %d: %v \n", countErr, row)
			arrErr = append(arrErr, row...)
		}
	}

	if countErr > 0 {
		fmt.Println("--------------------------------------------")
	}
	log.Println("Source >> Longman Communication 3135")

	res, err := util.PrettyStruct(arrErr)
	if err != nil {
		log.Fatal(err)
	}

	log.Println("total error >> ", res)
	fmt.Println("--------------------------------------------")

	return words
}

func LongmanCommunication9000(rows [][]string) []model.Word {
	// header >> 1: "High Frequency", 2: "Medium Frequency", 3: "Low Frequency"

	countID := 0
	countErr := 0
	arrErr := []string{}
	var words []model.Word

	for idr, row := range rows[1:] {
		for i, cellValue := range row {
			countID++
			matchedPartsOfSpeech, unmatchedParts := util.PrepareInformation(cellValue)

			// UPDATE
			if cellValue != "" && unmatchedParts == "" {
				if !util.ContainsComma(cellValue) && len(matchedPartsOfSpeech) > 1 {
					countErr++
					fmt.Println("--------------------------------------------")
					fmt.Printf(">> meaningless >> %d: %v \n", countErr, cellValue)
					arrErr = append(arrErr, cellValue)
				} else {
					used := false
					insertData := rows[idr][i]
					_, unmatchedParts2 := util.PrepareInformation(insertData)

					for idj, v := range words {
						if v.Vocabulary == util.RemoveNumbersFromWord(unmatchedParts2) {
							// Found!
							for idx := 0; idx < len(matchedPartsOfSpeech); idx++ {
								v.PartsOfSpeech = append(v.PartsOfSpeech, model.PartsOfSpeech{
									Type:          matchedPartsOfSpeech[idx],
									TranslationTH: "",
								})
							}

							used = true
							words[idj] = v
						}
					}

					if used == false {
						countErr++
						fmt.Println("--------------------------------------------")
						fmt.Printf(">> not used >> %d: %v \n", countErr, cellValue)
					}
				}
			}

			// NEW
			if cellValue != "" && unmatchedParts != "" {
				partsOfSpeechs := []model.PartsOfSpeech{}
				for idm := 0; idm < len(matchedPartsOfSpeech); idm++ {
					partsOfSpeechs = append(partsOfSpeechs, model.PartsOfSpeech{
						Type:          matchedPartsOfSpeech[idm],
						TranslationTH: "",
					})
				}

				// Map cell values to struct fields based on headers
				tag := []string{}
				switch i {
				case 0:
					tag = append(tag, "High Frequency")
				case 1:
					tag = append(tag, "Medium Frequency")
				case 2:
					tag = append(tag, "Low Frequency")
				}

				words = append(words, model.Word{
					ID:            countID,
					Vocabulary:    util.RemoveNumbersFromWord(unmatchedParts),
					Pronunciation: "",
					Phonetics:     "",
					PartsOfSpeech: partsOfSpeechs,
					Tag:           tag,
					Syn:           []string{},
					Source:        []string{"Longman Communication 9000"},
					Example:       []string{},
				})
			}
		}
	}

	if countErr > 0 {
		fmt.Println("--------------------------------------------")
	}
	log.Println("Source >> Longman Communication 9000")

	res, err := util.PrettyStruct(arrErr)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("total error >> ", res)
	fmt.Println("--------------------------------------------")

	return words
}
