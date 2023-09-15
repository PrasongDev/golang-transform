package util

import (
	"regexp"
	"strings"
)

const (
	noun          = "noun"         // A – คำนาม (Noun)
	verb          = "verb"         // B – คำกริยา (Verb)
	modalVerb     = "modal verb"   // B – คำกริยา (Verb)
	adjective     = "adjective"    // C – คำคุณศัพท์ (Adjective)
	adverb        = "adverb"       // D – คำกริยาวิเศษณ์ (Adverb)
	preposition   = "preposition"  // E – คำบุพบท (Preposition)
	pronoun       = "pronoun"      // F – คำสรรพนาม (Pronoun)
	conjunction   = "conjunction"  // G – คำเชื่อม (Conjunction)
	interjection  = "interjection" // H – คำอุทาน (Interjection)
	determiner    = "determiner"
	predeterminer = "predeterminer"
)

func PrepareInformation(sentence string) ([]string, string) {
	// Insert spaces after numbers in the input sentence
	inputWithSpaces := insertSpaceAfterNumbers(sentence)

	// Replace commas with spaces in the input string
	modifiedString := replaceCommasWithSpaces(inputWithSpaces)

	// Remove extra spaces from the input sentence
	cleanedSentence := removeExtraSpaces(modifiedString)

	// Find the parts of speech that match in the sentence
	matchedPartsOfSpeech, unmatchedParts := sentenceMatchesPattern(cleanedSentence)

	return matchedPartsOfSpeech, unmatchedParts
}

func sentenceMatchesPattern(sentence string) ([]string, string) {
	var partsOfSpeech = []string{noun, verb, modalVerb, adjective, adverb, preposition, pronoun, conjunction, interjection, determiner, predeterminer}

	// Convert the sentence to lowercase and split it into words
	words := strings.Fields(strings.ToLower(sentence))

	// Create a regular expression pattern
	pattern := strings.Join(partsOfSpeech, "|")
	re := regexp.MustCompile(pattern)

	// Find all matches in the sentence
	matches := re.FindAllString(strings.Join(words, " "), -1)

	// Create a string that contains the matched parts of speech
	matchedParts := strings.Join(matches, " ")

	// Create a string that contains the remaining words (not matched)
	unmatchedParts := ""
	for _, word := range words {
		if !strings.Contains(matchedParts, word) {
			unmatchedParts += word + " "
		}
	}
	unmatchedParts = strings.TrimSpace(unmatchedParts)

	return matches, unmatchedParts
}

func ContainsComma(input string) bool {
	return strings.Contains(input, ",")
}
