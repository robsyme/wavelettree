// Package provides datastructures to provide fast rank queries over
// non-binary alphabets.
package wavelettree

import (
	"github.com/willf/bitset"
)

type WaveletTree struct {
	depth     uint8
	zeros     *WaveletTree
	ones      *WaveletTree
	bitVector *bitset.BitSet
}

// New returns a pointer to a wavelet tree derived from the given
// bytes. At the moment, the package assumes an alphabet size of 256.
// This means that queries have to traverse 8 trees for each query. In
// a future release, you will be able to specify an alphabet which
// will speed up Rank queries, particularly for small alphabets (AGTC,
// for example).
func New(data []byte) *WaveletTree {
	tree := makeTree(data, 0)

	return tree
}

// Rank returns the number of occurences of 'query' in the first
// 'position' of the original byte array.
func (tree *WaveletTree) Rank(position uint, query byte) uint {
	return tree.rank(position, query, 0)
}

func (tree *WaveletTree) rank(position uint, query byte, depth uint) uint {
	boolQuery := (query>>depth)%2 == 1

	var nextTree *WaveletTree

	if boolQuery {
		nextTree = tree.ones
	} else {
		nextTree = tree.zeros
	}

	if nextTree != nil {
		nextPosition := tree.binaryRank(position, boolQuery)
		return nextTree.rank(nextPosition, query, depth+1)
	}

	return tree.binaryRank(position, boolQuery)
}

// TODO This should be implemented with RRR datastructures. This is
// just a placeholder.
func (tree *WaveletTree) binaryRank(position uint, query bool) uint {
	count := uint(0)
	for i := uint(0); i < uint(position); i++ {
		if tree.bitVector.Test(i) {
			count += 1
		}
	}
	if query {
		return count
	}
	return uint(position) - count
}

func makeTree(data []byte, depth uint) *WaveletTree {
	tree := WaveletTree{
		bitVector: bitset.New(uint(len(data))),
	}

	for i, b := range data {
		if (b>>depth)%2 == 1 {
			tree.bitVector.Set(uint(i))
		}
	}

	if depth < 8 && len(data) > 1 {
		zeros, ones := divideData(data, tree.bitVector, depth)
		tree.zeros = makeTree(zeros, depth+1)
		tree.ones = makeTree(ones, depth+1)
	}

	return &tree
}

func divideData(data []byte, bitVector *bitset.BitSet, depth uint) (zeros, ones []byte) {
	zeros = make([]byte, uint(len(data))-bitVector.Count())
	ones = make([]byte, bitVector.Count())

	zerocount, onecount := 0, 0
	for _, b := range data {
		if (b>>depth)%2 == 0 {
			zeros[zerocount] = b
			zerocount += 1
		} else {
			ones[onecount] = b
			onecount += 1
		}
	}
	return zeros, ones
}
