package longman

import (
	"fmt"
	"log"
	"transform/excel/model"
	"transform/excel/util"
)

// Source >> Longman Communication 9000
func LongmanCommunication9000(rows [][]string) []model.Word {
	// header >> 1: "High Frequency", 2: "Medium Frequency", 3: "Low Frequency"

	countID := 0
	countErr := 0
	arrErr := []string{}
	var words []model.Word

	for idr, row := range rows[1:] {
		for i, cellValue := range row {
			countID++
			var dataRow model.Word
			matchedPartsOfSpeech, unmatchedParts := util.PrepareInformation(cellValue)

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

			if cellValue != "" && unmatchedParts != "" {
				dataRow.ID = countID
				dataRow.Vocabulary = util.RemoveNumbersFromWord(unmatchedParts)

				for idm := 0; idm < len(matchedPartsOfSpeech); idm++ {
					dataRow.PartsOfSpeech = append(dataRow.PartsOfSpeech, model.PartsOfSpeech{
						Type:          matchedPartsOfSpeech[idm],
						TranslationTH: "",
					})
				}

				// Map cell values to struct fields based on headers
				switch i {
				case 0:
					dataRow.Tag = append(dataRow.Tag, "High Frequency")
				case 1:
					dataRow.Tag = append(dataRow.Tag, "Medium Frequency")
				case 2:
					dataRow.Tag = append(dataRow.Tag, "Low Frequency")
				}

				dataRow.Source = "Longman Communication 9000"
				words = append(words, dataRow)
			}
		}
	}

	if countErr > 0 {
		fmt.Println("--------------------------------------------")
	}

	res, err := util.PrettyStruct(arrErr)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("total error >> ", res)
	fmt.Println("--------------------------------------------")

	return words
}
