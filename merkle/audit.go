package merkle

type Direction int

const (
	left Direction = iota
	right
	oldRoot
)

type AuditNode struct {
	branch   Direction
	hashInfo string
}

func makeAuditNode(hash string, branch Direction) *AuditNode {
	a := AuditNode{hashInfo: hash, branch: branch}
	return &a
}
