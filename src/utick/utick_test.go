package utick

import (
	"fmt"
	"testing"
)

func TestInitialPos(t *testing.T) {
	pos := InitialPosition()
	fmt.Println(pos.Dump())
	fmt.Println(pos.LegalMoves())
	if pos.NextPlayer != PLAYER1 {
		t.Errorf("Next player is %d, expected %d", pos.NextPlayer, PLAYER1)
	}
	if pos.NextCell != ANY_CELL {
		t.Errorf("Next cell is %d, expected %d", pos.NextCell, ANY_CELL)
	}
	if pos.Result() != NONE {
		t.Errorf("Expected no result")
	}
}

func TestSomePlay(t *testing.T) {
	pos := InitialPosition()
	if err := pos.Play(0, 4); err != nil {
		t.Errorf("Failed to play to (0,4)")
	}
	fmt.Println(pos.Dump())
	fmt.Println(pos.LegalMoves())
	if pos.NextPlayer != PLAYER2 {
		t.Errorf("Next player is %d, expected %d", pos.NextPlayer, PLAYER1)
	}
	if pos.NextCell != 1 {
		t.Errorf("Next cell is %d, expected %d", pos.NextCell, 1)
	}
	if pos.Result() != NONE {
		t.Errorf("Expected no result")
	}
	if err := pos.Play(0, 4); err == nil {
		t.Errorf("Playing to occupied cell should be disallowed")
	}
	if err := pos.Play(1, 1); err == nil {
		t.Errorf("Playing to wrong metacell should be disallowed")
	}
	if err := pos.Play(2, 3); err != nil {
		t.Errorf("Playing to (2,3) should be ok")
	}
	fmt.Println(pos.Dump())
	fmt.Println(pos.LegalMoves())
	if pos.NextPlayer != PLAYER1 {
		t.Errorf("Next player is %d, expected %d", pos.NextPlayer, PLAYER1)
	}
	if pos.NextCell != 2*3+0 {
		t.Errorf("Next cell is %d, expected %d", pos.NextCell, 2*3+0)
	}
	if pos.Result() != NONE {
		t.Errorf("Expected no result")
	}
}

func ExampleGame() {
	pos := InitialPosition()
	pos.Play(1, 3)
	pos.Play(4, 1)
	pos.Play(3, 3)
	pos.Play(0, 0)
	pos.Play(1, 0)
	fmt.Println(pos.Dump())
	// Output:
	// OX. ... ...
	// ... .O. ...
	// ... ... ...
	//
	// .X. X.. ...
	// ... ... ...
	// ... ... ...
	//
	// ... ... ...
	// ... ... ...
	// ... ... ...
	// NextPlayer=2 NextCell=(1,0)
}

func ExampleGame2() {
	pos := InitialPosition()
	pos.Play(4, 4) // X
	pos.Play(5, 4) // O
	pos.Play(7, 4) // X
	pos.Play(5, 3) // O
	pos.Play(7, 1) // X
	pos.Play(5, 5) // O
	pos.Play(7, 7) // X
	fmt.Println(pos.Dump())
	// Output:
	// ... ... ...
	// ... ... .X.
	// ... ... ...
	//
	// ... ..O ...
	// ... .XO .X.
	// ... ..O ...
	//
	// ... ... ...
	// ... ... .X.
	// ... ... ...
	// NextPlayer=2 NextCell=ANY
}

func undumpMC(s string) MetaCell {
	v := int32(0)
	i := uint(0)
	for _, c := range s {
		if c == 'X' {
			v += int32(PLAYER1) << i
		} else if c == 'O' {
			v += int32(PLAYER2) << i
		} else if c == '.' {
			v += int32(NONE) << i
		} else if c == '?' {
			v += int32(DRAW) << i
		} else {
			continue
		}
		i += 2
	}
	return MetaCell(v)
}

func TestMetaCellResult(t *testing.T) {
	cases := []struct {
		mc string
		r  Cell
	}{
		{"... ... ...", NONE},
		{"XXX ... ...", PLAYER1},
		{".O. .O. .O.", PLAYER2},
		{"XXO OXX XOO", DRAW},
		{"OXO XOX XOX", DRAW},
		{"..O OX. XXX", PLAYER1},
		{"..O .O. OXX", PLAYER2},
		// examples with draws
		{"??O XXX OO.", PLAYER1},
		{"O?O XOX ..O", PLAYER2},
		{"XXO OX? X?O", DRAW},
		{"?.? XX. O.O", NONE},
	}
	for _, c := range cases {
		mc := undumpMC(c.mc)
		r := mc.Result()
		if r != c.r {
			t.Errorf("Wrong result for %s, want %d got %d", c.mc, c.r, r)
		}
	}
}

func TestRandomGame(t *testing.T) {
	player1 := NewRandomPlayerWithSeed(123)
	player2 := NewRandomPlayerWithSeed(456)
	pos := InitialPosition()
	for pos.Result() == NONE {
		fmt.Println(pos.Dump())
		var mv Coord
		if pos.NextPlayer == PLAYER1 {
			mv = player1.NextMove(pos)
		} else {
			mv = player2.NextMove(pos)
		}
		fmt.Printf("Player %d move: %v\n", pos.NextPlayer, mv)
		pos.PlayCoord(mv)
	}
	fmt.Printf("%s\nResult: %d\n", pos.Dump(), pos.Result())
}
