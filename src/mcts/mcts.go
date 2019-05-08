package mcts

import (
	"fmt"
	"math"
	"math/rand"
	"utick"
)

type Node struct {
	Pos          utick.Position
	LastMove     *utick.Coord
	Result       utick.Cell
	Visits       float64
	Wins         float64
	UntriedMoves []utick.Coord
	Parent       *Node
	ChildNodes   []*Node
}

func NewStartNode(pos utick.Position) (n Node) {
	n.Pos = pos
	n.LastMove = nil
	n.Result = utick.NONE
	n.Visits = 0
	n.Wins = 0
	n.UntriedMoves = pos.LegalMoves()
	n.Parent = nil
	n.ChildNodes = make([]*Node, 0)
	return n
}

func (n *Node) AddChild(mv utick.Coord) (nn *Node) {
	nn = new(Node)
	// mv must be in n.UntriedMoves
	ix := -1
	for i, c := range n.UntriedMoves {
		if c == mv {
			ix = i
			break
		}
	}
	if ix < 0 {
		panic("Move not in UntriedMoves")
	}
	n.UntriedMoves = append(n.UntriedMoves[:ix], n.UntriedMoves[ix+1:]...)

	newPos := n.Pos.Clone()
	if err := newPos.PlayCoord(mv); err != nil {
		panic("Play illegal move")
	}
	nn.Pos = newPos
	nn.LastMove = &mv
	nn.Result = newPos.Result()
	nn.Visits = 0
	nn.Wins = 0
	nn.UntriedMoves = newPos.LegalMoves()
	nn.Parent = n
	nn.ChildNodes = make([]*Node, 0)
	n.ChildNodes = append(n.ChildNodes, nn)
	return nn
}

func (n *Node) UpdateResult(r utick.Cell) {
	n.Visits += 1
	switch r {
	case utick.DRAW:
		n.Wins += 0.5
	case utick.PLAYER1:
		if n.Pos.NextPlayer == utick.PLAYER2 {
			n.Wins += 1
		}
	case utick.PLAYER2:
		if n.Pos.NextPlayer == utick.PLAYER1 {
			n.Wins += 1
		}
	}
}

// Select a child of a fully explored node
func (n *Node) SelectChild() *Node {
	if len(n.UntriedMoves) > 0 {
		panic("SelectChild called when UntriedMoves not empty")
	}
	if len(n.ChildNodes) == 0 {
		panic("SelectChild called when no child nodes")
	}
	bestIndex := -1
	bestScore := -1e99
	for ix, child := range n.ChildNodes {
		score := child.Wins/child.Visits + math.Sqrt(2*math.Log(n.Visits)/child.Visits)
		if score > bestScore {
			bestIndex, bestScore = ix, score
		}
	}
	return n.ChildNodes[bestIndex]
}

func SelectPath(start *Node) (*Node, []*Node) {
	t := []*Node{start}
	cur := start
	for len(cur.UntriedMoves) == 0 && cur.Result == utick.NONE {
		cur = cur.SelectChild()
		t = append(t, cur)
	}
	// t contains trace from start to cur
	return cur, t
}

func RandomPlayOut(r *rand.Rand, p utick.Position) utick.Cell {
	cur := p.Clone()
	res := cur.Result()
	for res == utick.NONE {
		mvs := cur.LegalMoves()
		mv := mvs[r.Intn(len(mvs))]
		cur.PlayCoord(mv)
		res = cur.Result()
	}
	return res
}

func NaiveMCTS(r *rand.Rand, startPos utick.Position, numTries int) utick.Coord {
	start := NewStartNode(startPos)
	for t := 0; t < numTries; t++ {
		// Select
		n, trace := SelectPath(&start)
		result := n.Result
		if result == utick.NONE {
			// Expand
			mv := n.UntriedMoves[r.Intn(len(n.UntriedMoves))]
			nn := n.AddChild(mv)
			trace = append(trace, nn)
			// Play out
			result = RandomPlayOut(r, nn.Pos)
		}
		// Propagate
		for _, node := range trace {
			node.UpdateResult(result)
		}
	}
	// Select move with max number of tries
	var best utick.Coord
	bestScore := 0.0
	for _, child := range start.ChildNodes {
		fmt.Printf("%v : %g\n", child.LastMove, child.Visits)
		if child.Visits > bestScore {
			best, bestScore = *child.LastMove, child.Visits
		}
	}
	return best
}

type NaiveMCTSPlayer struct {
	r        *rand.Rand
	numTries int
}

func (p NaiveMCTSPlayer) NextMove(pos utick.Position) utick.Coord {
	return NaiveMCTS(p.r, pos, p.numTries)
}

func NewNaiveMCTSPlayer() *NaiveMCTSPlayer {
	return &NaiveMCTSPlayer{rand.New(rand.NewSource(42)), 10000}
}
