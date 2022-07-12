package api

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestSpaFileHandler(t *testing.T) {
	h := spaFileHandler("testdata")

	assert.HTTPBodyContainsf(t, h, "GET", "/", nil, "index.html says hello", "should return index.html")
	assert.HTTPBodyContainsf(t, h, "GET", "./", nil, "index.html says hello", "should return index.html")
	assert.HTTPBodyContainsf(t, h, "GET", "/non-existing-file", nil, "index.html says hello", "should return index.html")
	assert.HTTPBodyContainsf(t, h, "GET", "/foo/bar/non-existing-file", nil, "index.html says hello", "should return index.html")

	assert.HTTPBodyContainsf(t, h, "GET", "/style.css", nil, "body {}", "/style.css should return style.css")

	assert.HTTPRedirect(t, h, "GET", "/index.html", nil, "index.html says hello", "/index.html should redirect to /")
	assert.HTTPError(t, h, "GET", "../", nil, "invalid paths should return an error")

}
