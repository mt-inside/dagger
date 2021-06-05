package node

import "github.com/mt-inside/dagger/pkg/value"

type BaseNode struct {
	curValue value.Value
	parent   Node
}

func (n *BaseNode) GetValue() value.Value {
	return n.curValue
}
func (n *BaseNode) GetParent() Node {
	return n.parent
}
