package main

import (
	"fmt"
	"os"
	"time"

	"github.com/davecgh/go-spew/spew"
	"github.com/google/cel-go/cel"
	"github.com/google/cel-go/checker/decls"
	"github.com/google/cel-go/common/types"
	"github.com/google/cel-go/common/types/ref"
	"github.com/google/cel-go/common/types/traits"
	"github.com/google/cel-go/interpreter/functions"
	"github.com/kr/pretty"
	"github.com/mt-inside/go-usvc"
	exprpb "google.golang.org/genproto/googleapis/api/expr/v1alpha1"
	rpcpb "google.golang.org/genproto/googleapis/rpc/context/attribute_context"
	structpb "google.golang.org/protobuf/types/known/structpb"
	tpb "google.golang.org/protobuf/types/known/timestamppb"
)

func main() {
	log := usvc.GetLogger(true)

	// Types for "contains" function
	typeParamA := decls.NewTypeParamType("A")
	typeParamB := decls.NewTypeParamType("B")
	mapAB := decls.NewMapType(typeParamA, typeParamB)

	env, err := cel.NewEnv(
		cel.Types(&rpcpb.AttributeContext_Request{}),
		cel.Declarations(
			decls.NewVar("request", decls.NewObjectType("google.rpc.context.AttributeContext.Request")),
			decls.NewFunction("contains",
				decls.NewParameterizedInstanceOverload(
					"map_contains_key_value",
					[]*exprpb.Type{mapAB, typeParamA, typeParamB},
					decls.Bool,
					[]string{"A", "B"},
				),
			),
		),
	)
	if err != nil {
		log.Error(err, "execution building execution environment")
		os.Exit(1)
	}

	ast, issues := env.Compile(`request.auth.claims.contains('group', 'admn')`) // Parse, Check
	if issues != nil && issues.Err() != nil {
		log.Error(issues.Err(), "compile error")
		os.Exit(1)
	}

	log.Info("Compiled", "return type", ast.ResultType(), "issues", issues)

	prg, err := env.Program( // Plan execution, bind functions, etc
		ast,
		cel.Functions(
			&functions.Overload{
				Operator: "map_contains_key_value",
				Function: mapContainsKeyValue,
			},
		),
	)
	if err != nil {
		log.Error(err, "execution planning error")
		os.Exit(1)
	}

	vars := map[string]interface{}{
		"request": &rpcpb.AttributeContext_Request{
			Auth: &rpcpb.AttributeContext_Auth{
				Principal: "matt",
				Claims:    &structpb.Struct{Fields: map[string]*structpb.Value{"group": structpb.NewStringValue("admin")}},
			},
			Time: &tpb.Timestamp{Seconds: time.Now().Unix()},
		},
	}
	fmt.Println("== input ==")
	pretty.Print(vars)
	out, details, err := prg.Eval(vars)
	if err != nil {
		log.Error(err, "execution error")
		os.Exit(1)
	}

	// Report() ?

	fmt.Println("== output ==")
	spew.Dump(details)
	fmt.Println(out)
}

// map.contains(key, value) bool
func mapContainsKeyValue(args ...ref.Val) ref.Val {
	if len(args) != 3 {
		return types.NewErr("no such overload")
	}
	obj := args[0]
	m, isMap := obj.(traits.Mapper)
	if !isMap {
		// The helper ValOrErr ensures that errors on input are propagated.
		return types.ValOrErr(obj, "no such overload")
	}

	// CEL has many interfaces for dealing with different type abstractions.
	// The traits.Mapper interface unifies field presence testing on proto
	// messages and maps.
	key := args[1]
	v, found := m.Find(key)
	// If not found and the value was non-nil, the value is an error per the
	// `Find` contract. Propagate it accordingly.
	if !found {
		if v != nil {
			return types.ValOrErr(v, "unsupported key type")
		}
		// Return CEL False if the key was not found.
		return types.False
	}
	// Otherwise whether the value at the key equals the value provided.
	return v.Equal(args[2])
}
