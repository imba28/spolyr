package language

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestDetector_Detect(t *testing.T) {
	d := New()

	testCases := []struct {
		query            string
		expectedLanguage string
	}{
		{
			query:            "Die Freiheit spielt auf allen Geigen",
			expectedLanguage: "german",
		},
		{
			query:            "Parti de rien, j’arrive au top",
			expectedLanguage: "french",
		},
		{
			query:            "There's vomit on his sweater already, mom's spaghetti",
			expectedLanguage: "english",
		},
	}

	for _, testCase := range testCases {
		language, err := d.Detect(testCase.query)
		assert.Equal(t, language, testCase.expectedLanguage)
		assert.Nil(t, err)
	}
}

func TestDetector_Detect__unkown_language(t *testing.T) {
	d := New()

	_, err := d.Detect("1234")
	assert.Error(t, err)
}

// make sure languages can be set using the language name
func TestWithLanguages(t *testing.T) {
	_, err := WithLanguages([]string{"german", "english", "french", "russian", "swedish"})

	assert.Nil(t, err)
}

func TestWithLanguages__invalid_language(t *testing.T) {
	_, err := WithLanguages([]string{"foobar"})

	assert.Error(t, err)
}
