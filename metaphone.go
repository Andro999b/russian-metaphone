package phonetics

import (
	"regexp"
	"strings"
)

const deafConsonantCriteria = "АЯОЫИЕЁЭЮЛМНР"
const deafConsonant = "БЗДВГ"
const allowedAlphabet = "ОЕАИУЭЮЯПСТРКЛМНБВГДЖЗЙФХЦЧШЩЁЫ"
const vowel = "ОЮЕЭЯЁЫ"

type phonemeReplacment struct {
	pattern *regexp.Regexp
	replace string
}

var phonemeReplacments []phonemeReplacment = nil

// EncodeRuMetaphone is a function to encode word with russian phinetics
func EncodeRuMetaphone(word string) string {
	word = strings.ToUpper(word)
	word = cleanUp(word) // remove all none allowed characters

	if len(word) == 1 {
		return word
	}

	word = replaceEndings(word) // do ending replacment
	word = replacePhonems(word)
	word = replaceVowel(word)
	word = reduceDeafConsonant(word)

	return word
}

func cleanUp(word string) string {
	var buff strings.Builder
	for _, ch := range word {
		if isInAlphabed(ch) {
			buff.WriteRune(ch)
		}
	}

	return buff.String()
}

func replacePhonems(word string) string {
	if phonemeReplacments == nil {
		phonemeReplacments = make([]phonemeReplacment, 3)
		phonemeReplacments = addPhonemeReplacment(phonemeReplacments, "(С?Т|С)Ч", "")
		phonemeReplacments = addPhonemeReplacment(phonemeReplacments, "(С?Т|С)Ч", "")
		phonemeReplacments = addPhonemeReplacment(phonemeReplacments, "(С?Т|С)Ч", "")
		phonemeReplacments = addPhonemeReplacment(phonemeReplacments, "(С?Т|С)Ч", "")
		phonemeReplacments = addPhonemeReplacment(phonemeReplacments, "(С?Т|С)Ч", "")
		phonemeReplacments = addPhonemeReplacment(phonemeReplacments, "(С?Т|С)Ч", "")
	}

	for _, pr := range phonemeReplacments {
		word = pr.pattern.ReplaceAllString(word, pr.replace)
	}

	return word
}

func addPhonemeReplacment(replacements []phonemeReplacment, pattern string, replace string) []phonemeReplacment {
	re, err := regexp.Compile(pattern)
	if err != nil {
		return replacements
	}

	return append(phonemeReplacments, phonemeReplacment{pattern: re, replace: replace})
}

func replaceEndings(word string) string {
	wordlen := len(word)
	ending := word[:6]
	switch ending {
	case "ОВСКИЙ":
		word = word[:wordlen-6] + "@"
	case "ЕВСКИЙ":
		word = word[:wordlen-6] + "#"
	case "ОВСКАЯ":
		word = word[:wordlen-6] + "$"
	case "ЕВСКАЯ":
		word = word[:wordlen-6] + "%"
	default:
		ending = word[:4]
		if ending == "ИЕВА" || ending == "ЕЕВА" {
			word = word[:wordlen-4] + "9"
		} else {
			ending = word[:3]
			switch ending {
			case "ОВА", "ЕВА":
				word = word[:wordlen-3] + "9"
			case "ИНА":
				word = word[:wordlen-3] + "1"
			case "ИЕВ", "ЕЕВ":
				word = word[:wordlen-3] + "4"
			case "НКО":
				word = word[:wordlen-3] + "3"
			default:
				ending = word[:2]
				switch ending {
				case "ОВ", "ЕВ":
					word = word[:wordlen-2] + "4"
				case "АЯ":
					word = word[:wordlen-2] + "6"
				case "ИЙ", "ЫЙ":
					word = word[:wordlen-2] + "7"
				case "ЫХ", "ИХ":
					word = word[:wordlen-2] + "5"
				case "ИН":
					word = word[:wordlen-2] + "8"
				case "ИК", "ЕК":
					word = word[:wordlen-2] + "2"
				case "УК", "ЮК":
					word = word[:wordlen-2] + "0"
				}
			}
		}
	}

	return word
}

func replaceVowel(word string) string {
	var buff strings.Builder
	for _, ch := range word {
		if isVowel(ch) {
			buff.WriteRune(getVowelReplacment(ch))
		} else {
			buff.WriteRune(ch)
		}
	}

	return buff.String()
}

func reduceDeafConsonant(word string) string {
	var buff strings.Builder

	wordlen := len(word)
	for i := 0; i < wordlen; i++ {
		ch := rune(word[i])
		if isDeafConsonant(ch) && (i == wordlen-1 || strings.ContainsRune(deafConsonantCriteria, rune(word[i+1]))) {
			buff.WriteRune(getDeafConsonantReplacment(ch))
		} else {
			buff.WriteRune(ch)
		}
	}

	return buff.String()
}

func getDeafConsonantReplacment(ch rune) rune {
	switch ch {
	case 'Б':
		return 'П'
	case 'З':
		return 'С'
	case 'Д':
		return 'Т'
	case 'В':
		return 'Ф'
	case 'Г':
		return 'К'
	}
	return ch
}

func isDeafConsonant(ch rune) bool {
	return strings.ContainsRune(deafConsonant, ch)
}

func getVowelReplacment(ch rune) rune {
	switch ch {
	case 'О', 'Ы', 'А', 'Я':
		return 'А'
	case 'Ю', 'У':
		return 'У'
	case 'Е', 'Ё', 'Э', 'И', 'Й':
		return 'И'
	}
	return ch
}

func isVowel(ch rune) bool {
	return strings.ContainsRune(vowel, ch)
}

func isInAlphabed(ch rune) bool {
	return strings.ContainsRune(allowedAlphabet, ch)
}
