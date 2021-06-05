package node

import (
	"github.com/go-logr/logr"
	"github.com/mt-inside/dagger/pkg/value"
	exprpb "google.golang.org/genproto/googleapis/api/expr/v1alpha1"
)

type Node interface {
	GetValue() value.Value
	Evaluate()
	GetParent() Node
	//DumpDot()
}

func BuildNode(log logr.Logger, parent Node, expr *exprpb.Expr) Node {
	switch expr.ExprKind.(type) {
	case *exprpb.Expr_CallExpr:
		return buildCallNode(log, parent, expr.GetCallExpr())
	case *exprpb.Expr_ConstExpr:
		return buildConstNode(log, parent, expr.GetConstExpr())
	default:
		panic("bottom")
	}
}

func propagate(n Node) {
	if n.GetParent() == nil {
		n.GetValue().Print()
		return
	}

	n.GetParent().Evaluate()
}
