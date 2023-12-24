package game

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGenerateWordsInRandomOrder(t *testing.T) {
	res1 := generateWordsInRandomOrder()
	res2 := generateWordsInRandomOrder()

	l := len(res1)
	sameWordsCount := 0
	for i := 0; i < l; i++ {
		if res1[i] == res2[i] {
			sameWordsCount++
		}
	}
	assert.Less(t, sameWordsCount, l)
}
