// Package provides datastructures to provide fast rank queries over
// non-binary alphabets.
package wavelettree

import (
	"bytes"
	"fmt"
	"github.com/robsyme/succinctBitSet"
)

type WaveletTree struct {
	depth     uint8
	zeros     *WaveletTree
	ones      *WaveletTree
	bitVector *succinctBitSet.BitSet
}

// New returns a pointer to a wavelet tree derived from the given
// bytes. At the moment, the package assumes an alphabet size of 256.
// This means that queries have to traverse Log2(256) = 8 trees for
// each query. In a future release, you will be able to specify an
// alphabet which will speed up Rank queries, particularly for small
// alphabets (AGTC, for example).
func New(data []byte) *WaveletTree {
	fmt.Printf("MAKING NEW TREE: %v\n", data)
	tree := makeTree(data, 0)
	fmt.Println("DEBUG")
	fmt.Println(tree)
	return tree
}

func (tree *WaveletTree) String() string {
	return tree.string(0)
}

func (tree *WaveletTree) string(depth int) string {
	var buffer bytes.Buffer

	for i := 0; i < depth; i++ {
		buffer.WriteByte(' ')
	}

	fmt.Fprintln(&buffer, tree.bitVector)

	if tree.ones != nil {
		depth++
		buffer.WriteString(tree.ones.string(depth))
	}

	if tree.zeros != nil {
		depth++
		buffer.WriteString(tree.zeros.string(depth))
	}

	return buffer.String()
}

// Rank returns the number of occurences of 'query' in the first
// 'position' of the original byte array.
func (tree *WaveletTree) Rank(position uint, query byte) uint {
	return tree.rank(position, query, 0)
}

func (tree *WaveletTree) rank(position uint, query byte, depth uint) uint {
	fmt.Printf("Looking for %c (%08b) at depth %d:\n", query, query, depth)
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
	fmt.Printf("binaryRank(%d, %t) = ", position, query)
	count := tree.bitVector.Rank(position)

	if query {
		fmt.Println(count)
		return count
	}
	fmt.Println(position - count)
	return uint(position) - count
}

func makeTree(data []byte, depth uint) *WaveletTree {
	tree := WaveletTree{
		bitVector: succinctBitSet.New(),
	}

	var outBuffer bytes.Buffer

	popcount := 0
	bits := make(chan bool, len(data))
	go func() {
		fmt.Fprintf(&outBuffer, "Depth=%d ", depth)
		outBuffer.WriteString("\033[39m[ ")
		for _, b := range data {
			if (b>>depth)%2 == 1 {
				popcount++
				bits <- true
				// var formatBuffer bytes.Buffer
				// if depth > 0 {
				// 	fmt.Fprintf(&formatBuffer, "\033[39m%%0%db\033[32m%%c\033[39m%%0%db ", 7-depth, depth)
				// 	fmt.Fprintf(&outBuffer, formatBuffer.String(), b>>depth, b, b&(1<<(depth)-1))
				// } else {
				// 	fmt.Fprintf(&formatBuffer, "\033[39m%%0%db\033[32m%%c ", 7-depth)
				// 	fmt.Fprintf(&outBuffer, formatBuffer.String(), b>>depth, b)
				// }
			} else {
				bits <- false
				// var formatBuffer bytes.Buffer
				// if depth > 0 {
				// 	fmt.Fprintf(&formatBuffer, "\033[39m%%0%db\033[31m%%c\033[39m%%0%db ", 7-depth, depth)
				// 	fmt.Fprintf(&outBuffer, formatBuffer.String(), b>>depth, b, b&(1<<(depth)-1))
				// } else {
				// 	fmt.Fprintf(&formatBuffer, "\033[39m%%0%db\033[31m%%c ", 7-depth)
				// 	fmt.Fprintf(&outBuffer, formatBuffer.String(), b>>depth, b)
				// }
			}
		}
		outBuffer.WriteString("\033[39m]")
		fmt.Println(outBuffer.String())
		close(bits)
	}()

	tree.bitVector.AddFromBoolChan(bits)

	fmt.Println(tree.bitVector)

	if depth < 8 && len(data) > 1 {
		zeros, ones := divideData(data, depth, popcount)
		if len(zeros) > 0 {
			tree.zeros = makeTree(zeros, depth+1)
		}
		if len(ones) > 0 {
			tree.ones = makeTree(ones, depth+1)
		}
	}
	return &tree
}

func divideData(data []byte, depth uint, popcount int) (zeros, ones []byte) {

	zeros = make([]byte, len(data)-popcount)
	ones = make([]byte, popcount)

	zerosCount := 0
	onesCount := 0
	for _, b := range data {
		if (b>>depth)%2 == 0 {
			zeros[zerosCount] = b
			zerosCount++
		} else {
			ones[onesCount] = b
			onesCount++
		}
	}
	return zeros, ones
}
