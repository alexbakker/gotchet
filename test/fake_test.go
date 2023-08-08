package test

import (
	"fmt"
	"os"
	"testing"

	"github.com/brianvoe/gofakeit/v6"
)

var (
	faker = gofakeit.NewCrypto()
)

func generateFakeTests(t *testing.T, multiplier int) {
	testCount := faker.IntRange(2*multiplier, 5*multiplier)
	for i := 0; i < testCount; i++ {
		t.Run(faker.BuzzWord(), func(t *testing.T) {
			lineCount := faker.IntRange(10, 100)
			for j := 0; j < lineCount; j++ {
				words := faker.IntRange(5, 20)
				fmt.Fprintln(os.Stderr, faker.Sentence(words))
			}
			if faker.IntRange(1, 10) == 10 {
				t.Error("Random failure")
			}
			if multiplier > 0 && faker.Bool() {
				generateFakeTests(t, multiplier-1)
			}
		})
	}
}

func TestGenerateFakeTests(t *testing.T) {
	generateFakeTests(t, 4)
}
