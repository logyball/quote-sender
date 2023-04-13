package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBuildTriviaMessage(t *testing.T) {
	testTriviaObject := TriviaObject{
		Category: "stuff",
		Question: "Question",
		Answer:   "Answer",
	}
	testTriviaObjectNoCat := TriviaObject{
		Question: "Question",
		Answer:   "Answer",
	}

	tests := map[string]struct {
		input TriviaObject
		want  string
	}{
		"simple":      {input: testTriviaObject, want: "🏆❓ It's Trivia Tuesday! ❓🏆\n\nToday's category is: stuff!\n\nQuestion: Question?\n\n...\n...\n...\n...\n\nAnswer: Answer"},
		"no category": {input: testTriviaObjectNoCat, want: "🏆❓ It's Trivia Tuesday! ❓🏆\n\nToday's category is: Free for all!\n\nQuestion: Question?\n\n...\n...\n...\n...\n\nAnswer: Answer"},
	}

	for name, tst := range tests {
		t.Run(name, func(t *testing.T) {
			res := MakeTriviaTwilioMessage(&tst.input)
			assert.Equal(t, tst.want, res, "not right message")
		})
	}
}

func TestTriviaDecode(t *testing.T) {
	resByteArr := []byte("[{\"category\": \"general\", \"question\": \"What is dendrochronology\", \"answer\": \"Tree Ring Dating\"}]")
	res, err := parseTriviaJsonResponse(resByteArr)

	if err != nil {
		assert.FailNowf(t, "err was not nil: ", "%s", err)
	}

	assert.Equal(t, "general", res.Category)
	assert.Equal(t, "What is dendrochronology", res.Question)
	assert.Equal(t, "Tree Ring Dating", res.Answer)
}
