package stringutils

import (
	"strings"
)

var (
	tobeRemovedNeighborhoodSuffixes = []string{
		" MAHALLESİ",
		" MAH",
		" MAH.",
		" KÖYÜ",
	}

	tobeCheckedDistrictSuffixes = []string{
		"MERKEZ",
	}
)

func ParseParentheses(str string) []string {
	openParenthesis := strings.Index(str, "(")
	closeParenthesis := strings.LastIndex(str, ")")

	if openParenthesis != -1 && closeParenthesis != -1 && openParenthesis < closeParenthesis {
		// Remove trailing and leading spaces
		wordBefore := strings.TrimSpace(str[:openParenthesis])
		wordInside := strings.TrimSpace(str[openParenthesis+1 : closeParenthesis])

		return []string{wordBefore, wordInside}
	}

	return []string{str}
}

func ParseOvoNeighborhood(str string) string {
	// Remove trailing and leading spaces
	str = strings.TrimSpace(str)

	// Remove suffixes
	for _, suffix := range tobeRemovedNeighborhoodSuffixes {
		if strings.HasSuffix(str, suffix) {
			str = strings.TrimSuffix(str, suffix)
			break
		}
	}

	return str
}

func ParseOvoDistrict(str string) string {
	// Remove trailing and leading spaces
	str = strings.TrimSpace(str)

	for _, suffix := range tobeCheckedDistrictSuffixes {
		if strings.HasSuffix(str, suffix) {
			return suffix
		}
	}

	return str
}
