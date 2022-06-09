package main

import (
	"os"
	"time"

	"github.com/google/cel-go/cel"
	"github.com/google/cel-go/checker/decls"
	"github.com/mt-inside/dagger/pkg/node"
	"github.com/mt-inside/go-usvc"
	exprpb "google.golang.org/genproto/googleapis/api/expr/v1alpha1"
)

func main() {
	log := usvc.GetLogger(true)

	env, err := cel.NewEnv(
		cel.Declarations(
			decls.NewFunction("default",
				decls.NewOverload(
					"default_int",
					[]*exprpb.Type{decls.Int, decls.Int}, // args
					decls.Int,                            // return
				),
			),
			decls.NewFunction("timestamp",
				decls.NewOverload(
					"timestamp_int",
					[]*exprpb.Type{decls.Int}, // args
					decls.Int,                 // return
				),
			),
			decls.NewFunction("streamRange",
				decls.NewOverload(
					"stream_range_int",
					[]*exprpb.Type{decls.Int, decls.Int}, // args
					decls.Int,                            // return
				),
			),
			decls.NewFunction("streamTcp",
				decls.NewOverload(
					"stream_tcp_int",
					[]*exprpb.Type{decls.String}, // args
					decls.Int,                    // return
				),
			),
			decls.NewFunction("streamHttp",
				decls.NewOverload(
					"stream_http_int",
					[]*exprpb.Type{decls.String}, // args
					decls.Int,                    // return
				),
			),
			decls.NewFunction("floor10",
				decls.NewOverload(
					"round_10_int",
					[]*exprpb.Type{decls.Int}, // args
					decls.Int,                 // return
				),
			),
		),
	)
	if err != nil {
		log.Error(err, "execution building execution environment")
		os.Exit(1)
	}

	// https://github.com/google/cel-spec/blob/master/doc/langdef.md
	ast, issues := env.Compile(`2 * streamTcp('localhost:1234') + streamHttp('http://localhost:8080')`)
	//ast, issues := env.Compile(`2 * default(streamTcp('localhost:1234'),1) + default(streamTcp('localhost:1235'),1)`)
	//ast, issues := env.Compile(`2 + streamHttp('http://localhost:8080') - streamHttp('http://localhost:8080')`) // Flips, would need a debounce mode (eg settle for 1s after an Update())
	//ast, issues := env.Compile(`floor10(streamRange(1,1))`)
	if issues != nil && issues.Err() != nil {
		log.Error(issues.Err(), "compile error")
		os.Exit(1)
	}

	tree := node.BuildNode(log, nil, ast.Expr())
	log.Info("Compiled", "return type", ast.ResultType(), "issues", issues)

	tree.GetValue().Print()
	time.Sleep(1000 * time.Second)
}
