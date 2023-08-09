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

func TestGenerateFakeMulti(t *testing.T) {
	generateFakeTests(t, 4, false)
}

func TestGenerateFakeSingle(t *testing.T) {
	generateFakeTestOutput(t)
}

func TestGenerateFakeTree(t *testing.T) {
	generateFakeTests(t, 4, true)
}

func generateFakeTests(t *testing.T, multiplier int, tree bool) {
	testCount := faker.IntRange(2*multiplier, 5*multiplier)
	for i := 0; i < testCount; i++ {
		t.Run(faker.BuzzWord(), func(t *testing.T) {
			generateFakeTestOutput(t)
			if tree && multiplier > 0 && faker.Bool() {
				generateFakeTests(t, multiplier-1, true)
			}
		})
	}
}

func generateFakeTestOutput(t *testing.T) {
	lineCount := faker.IntRange(10, 100)
	for j := 0; j < lineCount; j++ {
		words := faker.IntRange(5, 20)
		fmt.Fprintln(os.Stderr, faker.Sentence(words))
	}
	if faker.IntRange(1, 10) == 10 {
		t.Error("Random failure")
	}
}
