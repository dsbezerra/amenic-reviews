package providers

import (
	"fmt"
	"testing"
)

func TestReviewInfo(t *testing.T) {
	omelete := NewOmeleteProvider()

	actualResult, err := omelete.GetMovieReviewInfo("attack-on-titan-fim-do-mundo")
	if err != nil {
		t.Error(err)
	}

	expectedResult := &OmeleteReview{
		Headline: "Ataque dos Titãs: O Fim do Mundo | Crítica",
		Author:   "Gabriel Avila",
	}

	if !isEqualReview(actualResult, expectedResult) {
		t.Fatalf("Expected review: %s but got %s.\n", formatReview(expectedResult), formatReview(actualResult))
	}
}

func formatReview(review *OmeleteReview) string {
	if review == nil {
		return ""
	}

	s := fmt.Sprintf("\nReview:\n%s\n%s", review.Headline, review.Author)
	return s
}

func isEqualReview(reviewOne *OmeleteReview, reviewTwo *OmeleteReview) bool {
	return (reviewOne.Headline == reviewTwo.Headline &&
		reviewOne.Author == reviewTwo.Author)
}
