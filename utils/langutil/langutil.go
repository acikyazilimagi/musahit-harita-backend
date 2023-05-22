package langutil

import "strings"

func ConvertTurkishCharsToEnglish(str string) string {
	var turkishCharacters = map[string]string{
		"İ": "I",
		"ı": "i",
		"Ş": "S",
		"ş": "s",
		"Ğ": "G",
		"ğ": "g",
		"Ü": "U",
		"ü": "u",
		"Ö": "O",
		"ö": "o",
		"Ç": "C",
		"ç": "c",
	}

	var strBuilder strings.Builder
	// replace turkish characters with english characters with loop in the given string
	for _, ch := range str {
		if val, ok := turkishCharacters[string(ch)]; ok {
			str = strings.Replace(str, string(ch), val, -1)
		}

		strBuilder.WriteString(string(ch))
	}

	return str
}
