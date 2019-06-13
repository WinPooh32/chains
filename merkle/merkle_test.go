package merkle

import (
	"crypto/md5"
	"fmt"
	"runtime/debug"
	"testing"
)

func assertPanic(t *testing.T, testCase string, f func()) {
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("The code did not panic in case: %s", testCase)
		}
	}()
	f()
}

func md5Wrap(data []byte) string {
	sum := md5.Sum(data)
	return fmt.Sprintf("%X", sum)
}

func checkErr(e error, t *testing.T) {
	if e != nil {
		t.Fatal(e, string(debug.Stack()))
	}
}

func makeTestTree() *Tree {
	tree := MakeTree(md5Wrap)
	return tree
}

func makeTestLeaf(testData []byte, t *testing.T) (*Tree, *node) {
	tree := makeTestTree()
	n, err := tree.makeLeafNode(testData)
	checkErr(err, t)
	return tree, n
}

func TestMakeNilLeaf(t *testing.T) {
	tree := makeTestTree()
	_, err := tree.makeLeafNode(nil)
	if err == nil {
		t.Fatal("tree was successfully created with nil data array")
	}
}

func TestMakeTree(t *testing.T) {
	tree := MakeTree(md5Wrap)

	if tree == nil {
		t.Fatal()
	}

	if tree.data != nil {
		t.Fatal()
	}

	if tree.leaves != nil {
		t.Fatal()
	}

	if tree.root != nil {
		t.Fatal()
	}

	if tree.hashFunc([]byte("test")) != md5Wrap([]byte("test")) {
		t.Fatal()
	}
}

func testArrEq(a, b []byte) bool {

	// If one is nil, the other must also be nil.
	if (a == nil) != (b == nil) {
		return false
	}

	if len(a) != len(b) {
		return false
	}

	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}

	return true
}

func TestMakeLeafNode(t *testing.T) {

	cases := [][]byte{
		[]byte{},
		[]byte{1},
		[]byte{1, 2},
		[]byte{1, 2, 3},
	}

	for _, v := range cases {
		tree, n := makeTestLeaf(v, t)

		if tree == nil {
			t.Fatal()
		}

		if n == nil {
			t.Fatal()
		}

		if n.data == nil {
			t.Fatal()
		}

		if !testArrEq(n.data, v) {
			t.Fatal()
		}

		if n.hashInfo == "" {
			t.Fatal()
		}

		if tree.hashFunc(v) != n.hashInfo {
			t.Fatal()
		}
	}
}

func TestIsLeaf(t *testing.T) {
	tree := makeTestTree()
	n, err := tree.makeLeafNode([]byte{})

	if err != nil {
		t.Fatal(err)
	}

	if !n.isLeaf() {
		t.Fatal()
	}
}

func leafEqualsCases(t *testing.T) {
	tree := makeTestTree()
	n1, err := tree.makeLeafNode([]byte{1, 2, 3})
	checkErr(err, t)
	n2, err := tree.makeLeafNode([]byte{})
	checkErr(err, t)
	n3, err := tree.makeLeafNode([]byte{1, 2, 3})
	checkErr(err, t)

	if !n1.equals(n1) {
		t.Fatal()
	}

	if n1.equals(n2) {
		t.Fatal()
	}

	if !n1.equals(n3) || !n3.equals(n1) {
		t.Fatal()
	}
}

func TestEquals(t *testing.T) {
	leafEqualsCases(t)
}

func TestMakeLeaves(t *testing.T) {
	testLeaves := func(leaves []*node, dataSet [][]byte) {
		for i, v := range leaves {
			if !v.isLeaf() {
				t.Fatal()
			}

			if !testArrEq(v.data, dataSet[i]) {
				t.Fatal()
			}
		}
	}

	func() {
		tree := makeTestTree()
		leaves, err := tree.makeLeaves()
		checkErr(err, t)
		testLeaves(leaves, [][]byte{})
	}()

	cases := []([][]byte){
		{},
		{{1}},
		{{1}, {2}},
		{{1}, {2}, {3}},
	}

	for _, v := range cases {
		tree := makeTestTree()
		tree.data = v
		leaves, err := tree.makeLeaves()
		checkErr(err, t)
		testLeaves(leaves, v)
	}
}

func TestAuditProof(t *testing.T) {

	data := [][]byte{
		[]byte{6, 6, 6},
		[]byte{3, 2, 1},
		[]byte{1, 2},
		[]byte{4, 4},
		[]byte{5, 5},
		[]byte{7, 7},
	}

	tree := makeTestTree()

	err := tree.Insert(data)
	checkErr(err, t)

	t.Log("Tree:\n", tree)

	trail, err := tree.AuditProof("0CB988D042A7F28DD5FE2B55B3F5AC7A")
	checkErr(err, t)

	t.Log("Audit trail:", trail)

	verify, err := tree.VerifyAudit(trail, "0CB988D042A7F28DD5FE2B55B3F5AC7A")
	checkErr(err, t)

	t.Log("Verified:", verify)

	// t.Fatal()
}

func TestTree_VerifyAudit(t *testing.T) {

	data := [][]byte{
		[]byte{6, 6, 6},
		[]byte{3, 2, 1},
		[]byte{1, 2},
		[]byte{4, 4},
		[]byte{5, 5},
		[]byte{7, 7},
	}

	tree := makeTestTree()

	err := tree.Insert(data)
	checkErr(err, t)

	type args struct {
		auditTrail []AuditNode
		leafHash   string
	}
	tests := []struct {
		name    string
		t       *Tree
		args    args
		want    bool
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.t.VerifyAudit(tt.args.auditTrail, tt.args.leafHash)
			if (err != nil) != tt.wantErr {
				t.Errorf("Tree.VerifyAudit() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("Tree.VerifyAudit() = %v, want %v", got, tt.want)
			}
		})
	}
}
