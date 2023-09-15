package model

type PartsOfSpeech struct {
	Type          string `json:"type"`
	TranslationTH string `json:"translationTH"` // คำแปล
}

type Word struct {
	ID            int             `json:"id"`
	Vocabulary    string          `json:"vocabulary"`    // คำศัพท์
	Pronunciation string          `json:"pronunciation"` // คำอ่าน, การออกเสียง
	Phonetics     string          `json:"phonetics"`     // สัทศาสตร์
	PartsOfSpeech []PartsOfSpeech `json:"partsOfSpeech"`
	Tag           []string        `json:"tag"`
	Syn           []string        `json:"syn"`
	Source        []string        `json:"source"`
	Example       []string        `json:"example"`
}
