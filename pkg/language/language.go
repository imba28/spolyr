package language

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/pemistahl/lingua-go"
	"strings"
)

var ErrUnknownLanguage = errors.New("unknown language")

type language string

func (l language) linguaLanguage() (lingua.Language, error) {
	var language lingua.Language
	v := fmt.Sprintf("\"%s\"", strings.ToUpper(string(l)))
	err := json.Unmarshal([]byte(v), &language)
	if err != nil {
		return language, err
	}

	return language, nil
}

// use defaultLanguages supported by mongodb:
// https://www.mongodb.com/docs/manual/reference/text-search-languages/#std-label-text-search-languages
var defaultLanguages = []lingua.Language{
	lingua.Danish,
	lingua.Dutch,
	lingua.English,
	lingua.Finnish,
	lingua.French,
	lingua.French,
	lingua.German,
	lingua.Hungarian,
	lingua.Italian,
	lingua.Nynorsk,
	lingua.Bokmal,
	lingua.Portuguese,
	lingua.Romanian,
	lingua.Russian,
	lingua.Spanish,
	lingua.Swedish,
	lingua.Turkish,
}

type Detector struct {
	d lingua.LanguageDetector
}

func New() Detector {
	d := lingua.NewLanguageDetectorBuilder().
		FromLanguages(defaultLanguages...).
		Build()
	return Detector{d: d}
}
func WithLanguages(languages []string) (Detector, error) {
	var l []lingua.Language
	for i := range languages {
		ll, err := language(languages[i]).linguaLanguage()
		if err != nil {
			return Detector{}, errors.New("unknown ISO 639-1 code " + string(languages[i]))
		}
		l = append(l, ll)
	}

	d := lingua.NewLanguageDetectorBuilder().
		FromLanguages(l...).
		Build()
	return Detector{d: d}, nil
}

func (d Detector) Detect(text string) (string, error) {
	if language, exists := d.d.DetectLanguageOf(text); exists {
		return strings.ToLower(language.String()), nil
	}
	return "", ErrUnknownLanguage
}
