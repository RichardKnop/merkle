package merkle_test

import (
	"crypto/sha1"
	"fmt"
	"io"
	"testing"

	"github.com/RichardKnop/merkle"
)

func TestMerkleTree(t *testing.T) {
	var (
		dataBlock1 = "data block 1"
		dataBlock2 = "data block 2"
		dataBlock3 = "data block 3"
		dataBlock4 = "data block 4"
	)

	// Create leaf nodes
	leafNode1, err := merkle.NewLeafNode([]byte(dataBlock1))
	if err != nil {
		t.Error(err)
	}
	leafNode2, err := merkle.NewLeafNode([]byte(dataBlock2))
	if err != nil {
		t.Error(err)
	}
	leafNode3, err := merkle.NewLeafNode([]byte(dataBlock3))
	if err != nil {
		t.Error(err)
	}
	leafNode4, err := merkle.NewLeafNode([]byte(dataBlock4))
	if err != nil {
		t.Error(err)
	}

	// Build the merkle tree
	rootNode := merkle.BuildTree(leafNode1, leafNode2, leafNode3, leafNode4)

	// Assert parent checksums
	c1, err := leafNode1.Parent.Checksum()
	if err != nil {
		t.Error(err)
	}
	c2, err := leafNode2.Parent.Checksum()
	if err != nil {
		t.Error(err)
	}
	c3, err := leafNode3.Parent.Checksum()
	if err != nil {
		t.Error(err)
	}
	c4, err := leafNode4.Parent.Checksum()
	if err != nil {
		t.Error(err)
	}
	if string(c1) != string(c2) {
		t.Errorf("%s != %s, but they should be equal", c1, c2)
	}
	if string(c3) != string(c4) {
		t.Errorf("%s != %s, but they should be equal", c3, c4)
	}
	if string(c1) == string(c3) {
		t.Errorf("%s == %s, but they should be different", c1, c3)
	}

	// Assert root checksum
	expectedRootChecksum := "3a7ee98226fdcbd8006266271289ac6e9b0f16ce"
	c, err := rootNode.Checksum()
	if err != nil {
		t.Error(err)
	}
	rootChecksum := fmt.Sprintf("%x", c)
	if rootChecksum != expectedRootChecksum {
		t.Errorf("expected checksum %q, got %q", expectedRootChecksum, rootChecksum)
	}

	// The main advantage of a merkle tree is that one branch of the hash tree
	// can be downloaded at a time and the integrity of each branch can be checked
	// immediately, even though the whole tree is not available yet.
	// For example, the integrity of dataBlock2 can be verified immediately
	// if the tree already contains leafNode1 and leafNode3.Parent checksum
	// by hashing the data block and iteratively combining the result with leafNode1
	// and then leafNode3.Parent and finally comparing the result with the top hash.
	c1, err = leafNode1.Checksum()
	if err != nil {
		t.Error(err)
	}
	h := sha1.New()
	if _, err = io.WriteString(h, dataBlock2); err != nil {
		t.Error(err)
	}
	c2 = h.Sum(nil)
	h.Reset()
	if _, err = io.WriteString(h, string(c1)+string(c2)); err != nil {
		t.Error(err)
	}
	c3 = h.Sum(nil)
	h.Reset()
	c4, err = leafNode3.Parent.Checksum()
	if err != nil {
		t.Error(err)
	}
	if _, err := io.WriteString(h, string(c3)+string(c4)); err != nil {
		t.Error(err)
	}
	rootChecksum = fmt.Sprintf("%x", h.Sum(nil))
	if rootChecksum != expectedRootChecksum {
		t.Errorf("expected checksum %q, got %q", expectedRootChecksum, rootChecksum)
	}
}
