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

func TestUnregisterPlayer(t *testing.T) {
	h := newHub()
	r := h.AddRoomToHub()
	go r.Run()
	for i := 1; i < 5; i++ {
		p := &Player{
			Room:        r,
			Conn:        nil,
			MsgToPlayer: make(chan []byte),
			AckChan:     make(chan bool),
			ID:          i,
			Username:    "test1",
		}
		r.Register <- p
		<-p.AckChan
		close(p.AckChan)
	}

	wordsLength := len(r.Words)

	r.CurrentTurnIdx = 3
	r.unregisterPlayer(r.Players[4])
	assert.Equal(t, []int{1, 2, 3}, r.TurnOrder)
	assert.Equal(t, 0, r.CurrentTurnIdx)
	assert.Equal(t, len(r.Words), wordsLength-1)

	r.CurrentTurnIdx = 1
	r.unregisterPlayer(r.Players[2])
	assert.Equal(t, []int{1, 3}, r.TurnOrder)
	assert.Equal(t, 1, r.CurrentTurnIdx)
	assert.Equal(t, len(r.Words), wordsLength-2)

	r.unregisterPlayer(r.Players[1])
	assert.Equal(t, []int{3}, r.TurnOrder)
	assert.Equal(t, 0, r.CurrentTurnIdx)
	assert.Equal(t, len(r.Words), wordsLength-2)
}
