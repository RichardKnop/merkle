package merkle

import (
	"crypto/sha1"
	"errors"
	"io"
)

// NewLeafNode returns a new node from a data block using the provided
// crypto.Hash, and calculates the block's checksum
func NewLeafNode(b []byte) (*Node, error) {
	n := new(Node)
	// TODO - make hashing algorithm configurable
	h := sha1.New()
	if _, err := io.WriteString(h, string(b)); err != nil {
		return nil, err
	}
	// Sum appends the current hash to b and returns the resulting slice.
	// It does not change the underlying hash state.
	n.checksum = h.Sum(nil)
	return n, nil
}

// Node is a basic unit of a tree
type Node struct {
	Parent, Left, Right *Node
	checksum            []byte
}

// IsLeaf returns true if this is a leaf node (has no children)
func (n Node) IsLeaf() bool {
	return n.Left == nil && n.Right == nil
}

// Checksum returns the checksum of a data block if it is a leaf node,
// or the concatenated checksum of this node's children
func (n Node) Checksum() ([]byte, error) {
	if n.IsLeaf() {
		if n.checksum != nil {
			return n.checksum, nil
		}
		return nil, errors.New("Leaf node has no checksum")
	}

	leftChecksum, err := n.Left.Checksum()
	if err != nil {
		return nil, err
	}
	rightChecksum, err := n.Right.Checksum()
	if err != nil {
		return nil, err
	}

	h := sha1.New()
	// TODO - run these in parallel
	if _, err := io.WriteString(h, string(leftChecksum)); err != nil {
		return nil, err
	}
	if _, err := io.WriteString(h, string(rightChecksum)); err != nil {
		return nil, err
	}

	return h.Sum(nil), nil
}

// BuildTree builds a tree from leaf nodes and returns a root node
func BuildTree(theNodes ...*Node) *Node {
	var nodes []*Node
	for i := 0; i < len(theNodes)-1; i = i + 2 {
		parentNode := new(Node)
		parentNode.Left = theNodes[i]
		parentNode.Right = theNodes[i+1]
		parentNode.Left.Parent = parentNode
		parentNode.Right.Parent = parentNode
		nodes = append(nodes, parentNode)
	}
	if len(nodes) == 1 {
		return nodes[0]
	}
	return BuildTree(nodes...)
}
