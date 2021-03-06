package gocube

import (
	"errors"
	"strconv"
)

// CornerIndexes contains 8 sets of 3 values which corresponds to the x, y, and
// z sticker indexes for each corner piece.
var CornerIndexes = []int{
	51, 15, 35,
	44, 17, 33,
	45, 0, 29,
	38, 2, 27,
	53, 9, 24,
	42, 11, 26,
	47, 6, 18,
	36, 8, 20,
}

// CornerPieces contains 8 sets of 3 values which correspond to the x, y, and
// z stickers for each corner piece.
var CornerPieces = []int{
	6, 2, 4,
	5, 2, 4,
	6, 1, 4,
	5, 1, 4,
	6, 2, 3,
	5, 2, 3,
	6, 1, 3,
	5, 1, 3,
}

// EdgeIndexes contains 12 pairs of values which correspond to the sticker
// indexes of each edge.
var EdgeIndexes = []int{
	7, 19,
	23, 39,
	10, 25,
	21, 50,
	3, 46,
	5, 37,
	1, 28,
	30, 41,
	16, 34,
	32, 48,
	12, 52,
	14, 43,
}

// EdgePieces contains 12 pairs of values which correspond to the stickers of
// each edge.
var EdgePieces = []int{
	1, 3,
	3, 5,
	2, 3,
	3, 6,
	1, 6,
	1, 5,
	1, 4,
	4, 5,
	2, 4,
	4, 6,
	2, 6,
	2, 5,
}

// StickerCube converts a CubieCube to a StickerCube
func (c *CubieCube) StickerCube() StickerCube {
	res := SolvedStickerCube()

	// Insert the edge pieces.
	for i, piece := range c.Edges {
		pieceIdx := piece.Piece * 2
		s1, s2 := EdgePieces[pieceIdx], EdgePieces[pieceIdx+1]
		if piece.Flip {
			s1, s2 = s2, s1
		}
		destIdx := i * 2
		res[EdgeIndexes[destIdx]] = s1
		res[EdgeIndexes[destIdx+1]] = s2
	}

	// Insert the corner pieces.
	for i, piece := range c.Corners {
		idx := piece.Piece * 3
		s1, s2, s3 := CornerPieces[idx], CornerPieces[idx+1],
			CornerPieces[idx+2]

		// Transform corner piece to move to its current position.
		// If an odd number of quarter turns were needed to move it to this
		// position, the corner's permutation is in the odd-parity coset.
		difference := (piece.Piece ^ i) & 7
		if difference == 1 || difference == 2 || difference == 4 ||
			difference == 7 {
			s1, s3 = s3, s1
		}

		// Twist the corner piece
		if piece.Orientation == 2 {
			s1, s2, s3 = s3, s1, s2
		} else if piece.Orientation == 0 {
			s1, s2, s3 = s2, s3, s1
		}

		destIdx := i * 3
		res[CornerIndexes[destIdx]] = s1
		res[CornerIndexes[destIdx+1]] = s2
		res[CornerIndexes[destIdx+2]] = s3
	}

	return res
}

// CubieCube converts a StickerCube to a CubieCube.
func (s *StickerCube) CubieCube() (*CubieCube, error) {
	var result CubieCube

	// Translate corner pieces.
	for i := 0; i < 8; i++ {
		idx := i * 3
		stickers := [3]int{s[CornerIndexes[idx]], s[CornerIndexes[idx+1]],
			s[CornerIndexes[idx+2]]}
		piece, orientation, err := findCorner(stickers)
		if err != nil {
			return nil, err
		}
		result.Corners[i].Piece = piece
		result.Corners[i].Orientation = orientation
	}

	// Translate edge pieces.
	for i := 0; i < 12; i++ {
		idx := i * 2
		stickers := [2]int{s[EdgeIndexes[idx]], s[EdgeIndexes[idx+1]]}
		piece, flip, err := findEdge(stickers)
		if err != nil {
			return nil, err
		}
		result.Edges[i].Piece = piece
		result.Edges[i].Flip = flip
	}

	return &result, nil
}

// findCorner finds the physical corner given its three colors.
func findCorner(stickers [3]int) (idx int, orientation int, err error) {
	for i := 0; i < 8; i++ {
		start := i * 3
		if !setsEqual(stickers[:], CornerPieces[start:start+3]) {
			continue
		}
		orientation = listIndex(stickers[:], 1)
		if orientation == -1 {
			orientation = listIndex(stickers[:], 2)
		}
		return i, orientation, nil
	}
	return 0, 0, errors.New("unrecognized corner: " +
		strconv.Itoa(stickers[0]) + "," + strconv.Itoa(stickers[1]) + "," +
		strconv.Itoa(stickers[2]))
}

// findEdge finds the physical edge given its two colors.
func findEdge(stickers [2]int) (idx int, flip bool, err error) {
	for i := 0; i < 12; i++ {
		start := i * 2
		if !setsEqual(stickers[:], EdgePieces[start:start+2]) {
			continue
		}

		// Using the EO rules, we can tell if the edge is good or bad.
		flip = false
		if stickers[1] == 1 || stickers[1] == 2 {
			// Top/bottom color in the wrong direction.
			flip = true
		} else if stickers[1] == 3 || stickers[1] == 4 {
			if stickers[0] != 1 && stickers[0] != 2 {
				// E-Slice edge with left/right color facing wrong direction.
				flip = true
			}
		}
		return i, flip, nil
	}
	return 0, false, errors.New("unrecognized edge: " +
		strconv.Itoa(stickers[0]) + "," + strconv.Itoa(stickers[1]))
}

func listContains(list []int, num int) bool {
	for _, x := range list {
		if x == num {
			return true
		}
	}
	return false
}

func listIndex(list []int, num int) int {
	for i, x := range list {
		if x == num {
			return i
		}
	}
	return -1
}

func setsEqual(set1 []int, set2 []int) bool {
	if len(set1) != len(set2) {
		return false
	}
	for _, x := range set1 {
		if !listContains(set2, x) {
			return false
		}
	}
	return true
}
