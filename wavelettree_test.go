package wavelettree

import (
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestTreeContruction(t *testing.T) {
	Convey("Given a tree constructed from the string 'abracadabra'", t, func() {
		tree := New([]byte("abracadabra"))
		So(tree, ShouldNotBeNil)
		So(tree.bitVector.Any(), ShouldBeTrue)
	})
}

func TestTreeQuery(t *testing.T) {
	Convey("Given a tree constructed from the string 'abracadabra'", t, func() {
		tree := New([]byte("abracadabra"))

		Convey("The rank of 'a',5 should be 2", func() {
			rank := tree.Rank(5, byte('a'))
			So(rank, ShouldEqual, 2)
		})

		Convey("The rank of 'q', 10 should be 0", func() {
			rank := tree.Rank(10, byte('q'))
			So(rank, ShouldBeZeroValue)
		})

		Convey("When given a position query of 0, the answer must be 0", func() {
			So(tree.Rank(0, byte('a')), ShouldBeZeroValue)
			So(tree.Rank(0, byte('b')), ShouldBeZeroValue)
			So(tree.Rank(0, byte('c')), ShouldBeZeroValue)
			So(tree.Rank(0, byte('d')), ShouldBeZeroValue)
			So(tree.Rank(0, byte('Z')), ShouldBeZeroValue)
		})

		Convey("The ranking should handle position values larger than the original string length", func() {
			So(tree.Rank(100000, byte('a')), ShouldEqual, 5)
		})
	})
}
