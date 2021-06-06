package node

import (
	"fmt"
	"strings"

	"github.com/go-logr/logr"
	"github.com/mt-inside/dagger/pkg/value"
	exprpb "google.golang.org/genproto/googleapis/api/expr/v1alpha1"
)

type FuncNode struct {
	BaseNode

	f        func(...value.Value) value.Value
	children []Node
}

func buildCallNode(log logr.Logger, parent Node, f *exprpb.Expr_Call) Node {
	log.V(2).Info("func", "name", f.Function)

	// TODO: better, tag or something?
	if strings.HasPrefix(f.Function, "stream") {
		return buildStreamNode(log, parent, f)
	}

	node := &FuncNode{}
	node.parent = parent

	switch f.Function {
	case "default":
		node.f = func(xs ...value.Value) value.Value {
			expr := xs[0]
			if expr.Mode == value.Pending {
				return xs[1]
			}
			return expr
		}
	case "_+_":
		node.f = liftM(func(xs ...int64) int64 { return xs[0] + xs[1] })
	case "_-_":
		node.f = liftM(func(xs ...int64) int64 { return xs[0] - xs[1] })
	case "_*_":
		node.f = liftM(func(xs ...int64) int64 { return xs[0] * xs[1] })
	case "floor10":
		node.f = liftM(func(xs ...int64) int64 { return xs[0] - xs[0]%10 })
	default:
		panic("bottom")
	}

	for i, a := range f.Args {
		log.V(2).Info("arg", "ordinal", i, "type", fmt.Sprintf("%T", a.ExprKind), "expr", a.String())
		node.children = append(node.children, BuildNode(log, node, a))
	}

	return node
}
func (n *FuncNode) GetParent() Node {
	return n.parent
}
func (n *FuncNode) Evaluate() {
	var vs []value.Value
	for _, c := range n.children {
		vs = append(vs, c.GetValue())
	}
	val := n.f(vs...)
	if val != n.curValue {
		n.curValue = val
		propagate(n)
	}
}

func liftM(f func(...int64) int64) func(...value.Value) value.Value {
	return func(xs ...value.Value) value.Value {
		var ns []int64
		for _, x := range xs {
			if x.Mode == value.Pending {
				return value.NewPending()
			}
			ns = append(ns, x.Val)
		}
		return value.NewAvailable(f(ns...))
	}
}
