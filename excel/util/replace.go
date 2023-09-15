package util

import (
	"regexp"
	"strings"
)

func insertSpaceBeforeNumbers(sentence string) string {
	// Use a regular expression to insert a space before numbers
	re := regexp.MustCompile(`(\D)(\d)`)
	sentence = re.ReplaceAllString(sentence, "$1 $2")

	return sentence
}

func insertSpaceAfterNumbers(sentence string) string {
	// Use a regular expression to insert a space after numbers
	re := regexp.MustCompile(`(\d)`)
	sentence = re.ReplaceAllString(sentence, "$1 ")

	return sentence
}

func replaceCommasWithSpaces(input string) string {
	// Replace all commas with spaces
	result := strings.ReplaceAll(input, ",", " ")
	return result
}

func RemoveNumbersFromWord(word string) string {
	// Define a regular expression pattern to match numbers
	pattern := "[0-9]"

	// Compile the regular expression
	re := regexp.MustCompile(pattern)

	// Replace all matched numbers with an empty string
	cleanedWord := re.ReplaceAllString(word, "")

	return cleanedWord
}

func removeExtraSpaces(sentence string) string {
	// Use a regular expression to replace consecutive spaces with a single space
	re := regexp.MustCompile(`\s+`)
	sentence = re.ReplaceAllString(sentence, " ")

	// Trim leading and trailing spaces
	sentence = strings.TrimSpace(sentence)

	return sentence
}

func isOnlySpaces(input string) bool {
	// Remove all spaces from the input string
	withoutSpaces := strings.ReplaceAll(input, " ", "")

	// If the resulting string is empty, it means there were only spaces
	return withoutSpaces == ""
}
