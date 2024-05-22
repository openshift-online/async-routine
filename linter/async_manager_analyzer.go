package linter

import (
	"go/ast"
	"go/token"

	"golang.org/x/tools/go/analysis"
)

var AsyncRoutineManagerAnalyzer = &analysis.Analyzer{
	Name: "AsyncManager",
	Doc:  "Checks that the async manager is used instead of the plain `go` keyword.",
	Run:  run,
}

func fileHasDisableComment(file *ast.File, fset *token.FileSet) bool {
	for _, commentGroup := range file.Comments {
		for _, comment := range commentGroup.List {
			if comment.Text == "// nolint:AsyncManager" {
				return true
			}
		}
	}
	return false
}

func isErrGroupGoCall(node ast.Node) bool {
	call, ok := node.(*ast.CallExpr)
	if !ok {
		return false
	}

	fun, ok := call.Fun.(*ast.SelectorExpr)
	if !ok {
		return false
	}

	if fun.Sel.Name != "Go" {
		return false
	}

	params := call.Args
	if len(params) != 1 {
		// errgroup.Go takes only one parameter
		// this is not an errgroup.Go call
		return false
	}

	funcLit, ok := params[0].(*ast.FuncLit)
	if !ok {
		// errgroup.Go parameter must be a function
		// this is not an errgroup.Go call
		return false
	}

	if funcLit.Type != nil &&
		funcLit.Type.TypeParams != nil &&
		len(funcLit.Type.TypeParams.List) != 0 {
		// errgroup.Go takes a `func() error` parameter.
		// This function in the inspected code is taking some parameters, so it is not an errgroup.Go
		// call
		return false
	}

	funcResults := funcLit.Type.Results.List
	if len(funcResults) != 1 {
		// errgroup.Go takes a `func() error` parameter.
		// This function in the inspected code is returning more than one result, so it is not an
		// errgroup.Go call
		return false
	}

	result, ok := funcResults[0].Type.(*ast.Ident)
	if !ok {
		// something weird happened. We should never enter here
		return false
	}
	if result.Name != "error" {
		// The function is not returning an `error`, so this is not an errgroup.Go call
		return false
	}

	return true
}

func run(pass *analysis.Pass) (interface{}, error) {
	inspect := func(node ast.Node) bool {
		if _, ok := node.(*ast.GoStmt); ok {
			pass.Reportf(node.Pos(), "Please replace the 'go' call with the equivalent "+
				"NewAsyncRoutine call")
			return true
		}

		if isErrGroupGoCall(node) {
			pass.Reportf(node.Pos(), "Please replace the errgroup.Group.Go call with the equivalent "+
				"NewAsyncRoutineWithErrGroup call")
		}

		return true
	}

	for _, f := range pass.Files {
		if fileHasDisableComment(f, pass.Fset) {
			continue
		}
		ast.Inspect(f, inspect)
	}
	return nil, nil
}
