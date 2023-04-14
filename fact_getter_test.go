package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBuildFactMessage(t *testing.T) {
	testFactObject := FunFact{
		Fact: "heres a fact",
	}
	nilFactObject := FunFact{}

	tests := map[string]struct {
		input     FunFact
		want      string
		wantError bool
	}{
		"simple":      {input: testFactObject, want: "📣 It's Fun Fact Saturday! 📣\n\nToday's fun fact: heres a fact", wantError: false},
		"no category": {input: nilFactObject, want: "", wantError: true},
	}

	for name, tst := range tests {
		t.Run(name, func(t *testing.T) {
			res, err := BuildFactTwilioMessage(&tst.input)
			if tst.wantError {
				assert.Error(t, err, "error that should be there wasn't")
			} else {
				assert.Equal(t, tst.want, res, "fact results was wrong")
			}
		})
	}
}

func TestFactDecode(t *testing.T) {
	resByteArr := []byte("[{\"fact\": \"The first theatre to show motion pictures was the Nickelodeon on June 19, 1905 in Pittsburgh, Pennsylvania. It was opened by Harry Davis on Smithfield Street\"}]")
	res, err := parseFactJsonResponse(resByteArr)

	if err != nil {
		assert.FailNowf(t, "err was not nil: ", "%s", err)
	}

	assert.Equal(t, "The first theatre to show motion pictures was the Nickelodeon on June 19, 1905 in Pittsburgh, Pennsylvania. It was opened by Harry Davis on Smithfield Street", res.Fact)
}
