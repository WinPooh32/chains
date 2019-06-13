package merkle

import (
	"fmt"
	"strconv"
)

type Tree struct {
	root     *node
	leaves   []*node
	data     [][]byte
	hashFunc func([]byte) string
}

func MakeTree(hasher func([]byte) string) *Tree {
	t := Tree{hashFunc: hasher}
	return &t
}

func (t *Tree) makeLeafNode(data []byte) (*node, error) {
	if data == nil {
		return nil, fmt.Errorf("data can't be a nil")
	}
	n := node{hashInfo: t.hashFunc(data), data: data}
	return &n, nil
}

func (t *Tree) makeLeaves() ([]*node, error) {
	nodes := make([]*node, 0, 4)

	for _, v := range t.data {
		n, err := t.makeLeafNode(v)
		if err != nil {
			return nil, err
		}

		nodes = append(nodes, n)
	}

	return nodes, nil
}

func (t *Tree) makeNode(left, right *node) (*node, error) {
	if left == nil && right == nil {
		return nil, fmt.Errorf("at least one node should not be a nil")
	}

	var hash string
	n := node{}

	if right == nil {
		hash = left.hashInfo
		left.parent = &n
	} else {
		left.parent = &n
		right.parent = &n
		hash = left.hashInfo + right.hashInfo
	}

	n.hashInfo = t.hashFunc([]byte(hash))
	n.left = left
	n.right = right

	return &n, nil
}

func (t *Tree) print(root *node, indent int, line string) (out string) {
	if root == nil {
		return
	}

	leafData := ""

	if root.isLeaf() {
		leafData = " - " + fmt.Sprintf("%v", root.data)
	}

	format := "%" + strconv.Itoa((len(root.hashInfo)+len(line))*indent+len(leafData)) + "s\n"

	out = fmt.Sprint(
		t.print(root.left, indent+1, "/"),
		fmt.Sprintf(format, line+root.hashInfo+leafData),
		t.print(root.right, indent+1, "\\"),
	)

	return
}

func (t *Tree) String() string {
	return t.print(t.root, 1, "")
}

func (t *Tree) Hash() string {
	return t.root.hashInfo
}

func (t *Tree) Insert(datas [][]byte) error {
	t.data = append(t.data, datas...)

	leaves, err := t.makeLeaves()
	if err != nil {
		return fmt.Errorf("merkle tree insert error: %s", err)
	}

	t.leaves = leaves

	if err := t.build(t.leaves); err != nil {
		return err
	}

	return nil
}

func (t *Tree) AuditProof(leafHash string) ([]AuditNode, error) {
	auditTrail := make([]AuditNode, 0, 4)

	if leaf := t.findLeaf(leafHash); leaf != nil {
		if leaf.parent == nil {
			return nil, fmt.Errorf("expected leaf to have a parent")
		}
		parent := leaf.parent
		if err := t.buildAuditTrail(&auditTrail, parent, leaf); err != nil {
			return nil, err
		}
	}

	return auditTrail, nil
}

func (t *Tree) buildAuditTrail(auditTrail *[]AuditNode, parent *node, child *node) error {
	if parent != nil {
		if !parent.equals(child.parent) {
			return fmt.Errorf("parent of child is not expected parent")
		}

		var nextChild *node
		var branch Direction

		if child.equals(parent.left) {
			nextChild = parent.right
			branch = left
		} else {
			nextChild = parent.left
			branch = right
		}

		if nextChild != nil {
			*auditTrail = append(*auditTrail, *makeAuditNode(nextChild.hashInfo, branch))
		}

		t.buildAuditTrail(auditTrail, child.parent.parent, child.parent)
	}

	return nil
}

func (t *Tree) findLeaf(hash string) *node {
	for _, v := range t.leaves {
		if v.hashInfo == hash {
			return v
		}
	}
	return nil
}

func (t *Tree) build(nodes []*node) error {
	if len(nodes) == 1 {
		t.root = nodes[0]
	} else {
		parents := make([]*node, 0, 4)
		length := len(nodes)

		for i := 0; i < length; i += 2 {
			var right, parent *node

			if i+1 < length {
				right = nodes[i+1]
			} else {
				right = nil
			}

			n, err := t.makeNode(nodes[i], right)
			if err != nil {
				return err
			}

			parent = n
			parents = append(parents, parent)
		}

		if err := t.build(parents); err != nil {
			return err
		}
	}

	return nil
}

func (t *Tree) VerifyAudit(auditTrail []AuditNode, leafHash string) (bool, error) {
	if len(auditTrail) == 0 {
		return false, fmt.Errorf("audit trail cannot be empty")
	}

	testHash := leafHash

	for _, v := range auditTrail {
		switch v.branch {
		case left:
			testHash = t.hashFunc([]byte(testHash + v.hashInfo))
		case right:
			testHash = t.hashFunc([]byte(v.hashInfo + testHash))
		}
	}

	return t.root.hashInfo == testHash, nil
}
