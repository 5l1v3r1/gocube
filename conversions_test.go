package gocube

import (
	"testing"
)

func TestStickerToCubieIdentity(t *testing.T) {
	stickers, err := ParseStickerCube("111111111 222222222 333333333 " + 
		"444444444 555555555 666666666")
	if err != nil {
		t.Error(err)
		return
	}
	cubies, err := StickerToCubie(*stickers)
	if err != nil {
		t.Error(err)
		return
	}
	for i := 0; i < 8; i++ {
		if cubies.Corners[i].Piece != i || cubies.Corners[i].Orientation != 0 {
			t.Error("Invalid corner at index", i)
		}
	}
	
	for i := 0; i < 12; i++ {
		if cubies.Edges[i].Piece != i || cubies.Edges[i].Flip {
			t.Error("Invalid edge at index", i)
		}
	}
}

func TestStickerToCube(t *testing.T) {
	// I did the algorithm B U D B' L2 D' R' F2 L F D2 R2 F' U2 R B2 L' U'
	stickers, err := ParseStickerCube("OGBYWWOOY OWOGYGGBR WBBGGOBRB " +
		"RWYYBRWYR RBWWROWYG GRGBORYOY")
	if err != nil {
		t.Error(err)
		return
	}
	cubies, err := StickerToCubie(*stickers)
	if err != nil {
		t.Error(err)
		return
	}
	
	// Run algorithm for comparison
	answer := SolvedCubieCube()
	moves, _ := ParseMoves("B U D B' L2 D' R' F2 L F D2 R2 F' U2 R B2 L' U'")
	for _, move := range moves {
		answer.Move(move)
	}
	
	// Make sure the cubes are equal
	for i, x := range answer.Corners {
		c := cubies.Corners[i]
		if x.Piece != c.Piece || x.Orientation != c.Orientation {
			t.Error("Invalid corner at index", i)
		}
	}
	for i, x := range answer.Edges {
		e := cubies.Edges[i]
		if x.Piece != e.Piece || x.Flip != e.Flip {
			t.Error("Invalid edge at index", i)
		}
	}
}
