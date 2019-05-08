// "Ultimate tick-tack-toe"
package utick

import (
	"errors"
	"fmt"
	"math/rand"
	"strings"
	"time"
)

type Cell int32

const (
	NONE Cell = iota
	PLAYER1
	PLAYER2
	DRAW // this is for a meta-cell result
)

type MetaCell int32 // MetaCell represents nine cells packed in 18 bits

type UnpackedMetaCell [9]Cell

const (
	ANY_CELL = -1
)

type Position struct {
	Cells      [9]MetaCell
	NextPlayer Cell // PLAYER1 or PLAYER2
	NextCell   int  // 0-8 or ANY_CELL
}

func (p Position) Clone() (q Position) {
	q.Cells = p.Cells
	q.NextPlayer = p.NextPlayer
	q.NextCell = p.NextCell
	return q
}

func (m MetaCell) Unpack() (u UnpackedMetaCell) {
	for i := 0; i < 9; i++ {
		u[i] = Cell((m >> (2 * uint(i))) & 3)
	}
	return u
}

func (u UnpackedMetaCell) Pack() MetaCell {
	m := int32(0)
	for i := 0; i < 9; i++ {
		m += (int32)(u[i]) << (2 * uint(i))
	}
	return MetaCell(m)
}

var MetaCellResult [1 << 18]Cell
var metaCellResultInited bool

func InitMetaCellResult() {
	if metaCellResultInited {
		return
	}
	rows := [][]int{
		{0, 1, 2},
		{3, 4, 5},
		{6, 7, 8},
		{0, 3, 6},
		{1, 4, 7},
		{2, 5, 8},
		{0, 4, 8},
		{2, 4, 6},
	}
	for m := 0; m < 1<<18; m++ {
		u := MetaCell(m).Unpack()
		res := NONE
		draw := true
		for _, r := range rows {
			n1, n2, nd := 0, 0, 0
			for _, x := range r {
				switch u[x] {
				case PLAYER1:
					n1++
				case PLAYER2:
					n2++
				case DRAW:
					nd++
				}
			}
			if n1 == 3 {
				res = PLAYER1
				break
			} else if n2 == 3 {
				res = PLAYER2
				break
			} else if (n1 == 0 || n2 == 0) && nd == 0 {
				draw = false
			}
		}
		if res == NONE && draw {
			res = DRAW
		}
		MetaCellResult[m] = res
	}
	metaCellResultInited = true
}

func (m MetaCell) Result() Cell {
	InitMetaCellResult()
	return MetaCellResult[m]
}

func (p Position) Result() Cell {
	InitMetaCellResult()
	var meta UnpackedMetaCell
	for i, mc := range p.Cells {
		meta[i] = MetaCellResult[mc]
	}
	return MetaCellResult[meta.Pack()]
}

func (p Position) Get(i, j int) Cell {
	if i < 0 || i > 8 || j < 0 || j > 8 {
		panic("Out of range")
	}
	a := (i/3)*3 + (j / 3)
	b := (i%3)*3 + (j % 3)
	return p.Cells[a].Unpack()[b]
}

func (p Position) Dump() string {
	var s strings.Builder
	for j := 0; j < 9; j++ {
		for i := 0; i < 9; i++ {
			c := p.Get(i, j)
			switch c {
			case NONE:
				s.WriteRune('.')
			case PLAYER1:
				s.WriteRune('X')
			case PLAYER2:
				s.WriteRune('O')
			case DRAW: // can't happen
				panic("Bad pos")
			}
			if i == 2 || i == 5 {
				s.WriteRune(' ')
			}

		}
		s.WriteRune('\n')
		if j == 2 || j == 5 {
			s.WriteRune('\n')
		}
	}
	s.WriteString(fmt.Sprintf("NextPlayer=%d NextCell=", p.NextPlayer))
	if p.NextCell == ANY_CELL {
		s.WriteString("ANY")
	} else {
		s.WriteString(fmt.Sprintf("(%d,%d)", p.NextCell/3, p.NextCell%3))
	}
	return s.String()
}

func InitialPosition() (p Position) {
	p.NextPlayer = PLAYER1
	p.NextCell = ANY_CELL
	return p
}

type Coord struct {
	i, j int
}

func (p *Position) Play(i, j int) error {
	if i < 0 || i > 8 || j < 0 || j > 8 {
		panic("Out of range")
	}
	a := (i/3)*3 + (j / 3)
	b := (i%3)*3 + (j % 3)
	if p.NextCell != ANY_CELL && p.NextCell != a {
		return errors.New("Disallowed metacell")
	}
	u := p.Cells[a].Unpack()
	if u[b] != NONE {
		return errors.New("Cell is used")
	}
	u[b] = p.NextPlayer
	p.Cells[a] = u.Pack()
	if p.Cells[b].Result() == NONE {
		p.NextCell = b
	} else {
		p.NextCell = ANY_CELL
	}
	if p.NextPlayer == PLAYER1 {
		p.NextPlayer = PLAYER2
	} else {
		p.NextPlayer = PLAYER1
	}
	return nil
}

func (p *Position) PlayCoord(c Coord) error {
	return p.Play(c.i, c.j)
}

func (p Position) LegalMoves() []Coord {
	var v []Coord
	if p.Result() != NONE {
		return v
	}
	scanMC := func(mc int) {
		ii, jj := mc/3, mc%3
		u := p.Cells[mc].Unpack()
		for i := 0; i < 3; i++ {
			for j := 0; j < 3; j++ {
				if u[i*3+j] == NONE {
					v = append(v, Coord{ii*3 + i, jj*3 + j})
				}
			}
		}
	}
	if p.NextCell != ANY_CELL {
		scanMC(p.NextCell)
	} else {
		for ij, c := range p.Cells {
			if c.Result() == NONE {
				scanMC(ij)
			}
		}
	}
	return v
}

func (p Position) IsLegalMove(mv Coord) bool {
	if p.NextCell != ANY_CELL {
		ii, jj := mv.i/3, mv.j/3
		if p.NextCell != ii*3+jj {
			return false
		}
	}
	return p.Get(mv.i, mv.j) == NONE
}

type Player interface {
	NextMove(p Position) Coord
}

type RandomPlayer struct {
	r *rand.Rand
}

func (p RandomPlayer) NextMove(pos Position) Coord {
	mvs := pos.LegalMoves()
	return mvs[p.r.Intn(len(mvs))]
}

func NewRandomPlayerWithSeed(seed int64) *RandomPlayer {
	return &RandomPlayer{
		r: rand.New(rand.NewSource(seed)),
	}
}

func NewRandomPlayer() *RandomPlayer {
	return NewRandomPlayerWithSeed(time.Now().UnixNano())
}
