package node

import (
	"github.com/go-logr/logr"
	"github.com/mt-inside/dagger/pkg/value"
	exprpb "google.golang.org/genproto/googleapis/api/expr/v1alpha1"
)

type ConstNode struct {
	BaseNode
}

func buildConstNode(log logr.Logger, parent Node, c *exprpb.Constant) Node {
	val := value.NewAvailable(c.GetInt64Value())
	log.V(2).Info("const", "val", val)
	return &ConstNode{BaseNode{parent: parent, curValue: val}}
}
func (n *ConstNode) GetParent() Node {
	return n.parent
}
func (n *ConstNode) Evaluate() {
	// tempus fugit
}
