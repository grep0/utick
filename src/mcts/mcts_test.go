package mcts

import (
	"fmt"
	"testing"
	"utick"
)

func TestGameAgainstRandomOpponent(t *testing.T) {
	player1 := NewNaiveMCTSPlayer()
	player2 := utick.NewRandomPlayerWithSeed(456)
	pos := utick.InitialPosition()
	for pos.Result() == utick.NONE {
		fmt.Println(pos.Dump())
		var mv utick.Coord
		if pos.NextPlayer == utick.PLAYER1 {
			mv = player1.NextMove(pos)
		} else {
			mv = player2.NextMove(pos)
		}
		fmt.Printf("Player %d move: %v\n", pos.NextPlayer, mv)
		pos.PlayCoord(mv)
	}
	fmt.Printf("%s\nResult: %d\n", pos.Dump(), pos.Result())
}

func TestGameAgainstSelf(t *testing.T) {
	player := NewNaiveMCTSPlayer()
	pos := utick.InitialPosition()
	for pos.Result() == utick.NONE {
		fmt.Println(pos.Dump())
		mv := player.NextMove(pos)
		fmt.Printf("Player %d move: %v\n", pos.NextPlayer, mv)
		pos.PlayCoord(mv)
	}
	fmt.Printf("%s\nResult: %d\n", pos.Dump(), pos.Result())
}
