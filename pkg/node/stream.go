package node

import (
	"time"

	"github.com/go-logr/logr"
	"github.com/mt-inside/dagger/pkg/stream"
	"github.com/mt-inside/dagger/pkg/value"
	exprpb "google.golang.org/genproto/googleapis/api/expr/v1alpha1"
)

type StreamNode struct {
	BaseNode
}

func buildStreamNode(log logr.Logger, parent Node, f *exprpb.Expr_Call) Node {
	log.V(2).Info("stream", "type", f.Function)

	node := &StreamNode{BaseNode{parent: parent}}

	switch f.Function {
	case "streamRange":
		go stream.Range(
			log,
			node,
			f.Args[0].GetConstExpr().GetInt64Value(),
			time.Duration(f.Args[1].GetConstExpr().GetInt64Value())*time.Second,
		)
	case "streamTcp":
		go stream.Tcp(
			log,
			node,
			f.Args[0].GetConstExpr().GetStringValue(),
		)
	case "streamHttp":
		go stream.Http(
			log,
			node,
			f.Args[0].GetConstExpr().GetStringValue(),
		)
	default:
		panic("bottom")
	}

	return node
}

func (n *StreamNode) Evaluate() {
	panic("stream nodes shouldn't have children")
}

func (n *StreamNode) Update(x value.Value) {
	if x != n.curValue {
		n.curValue = x
		propagate(n)
	}
}
