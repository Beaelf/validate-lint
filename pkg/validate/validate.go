package validate

import (
	"go/ast"
	"golang.org/x/tools/go/analysis"
	"strings"
)

var Validate = &analysis.Analyzer{
	Name: "validateTagRule",
	Doc:  "check if a struct type with validate tag has a corresponding test",
	Run:  run,
}

func run(pass *analysis.Pass) (interface{}, error) {
	inspect := func(n ast.Node) bool {
		if ts, ok := n.(*ast.TypeSpec); ok {
			if st, ok := ts.Type.(*ast.StructType); ok {
				for _, field := range st.Fields.List {
					if field.Tag != nil && strings.Contains(field.Tag.Value, "validate") {
						structName := ts.Name.Name
						//// struct with validate tag must has test

						if /*(args.Option == "default" || args.Option == args.RequireStruct) &&*/ !hasTestForStruct(pass, structName) {
							pass.Reportf(ts.Pos(), "struct %s has validate tag but no corresponding test", structName)
						}
						//fmt.Println("next")
						if /*(args.Option == "default" || args.Option == args.RequireDive) && */ strings.Contains(field.Tag.Value, "dive") {
							if !isSliceOrMapType(field.Type) {
								pass.Reportf(field.Pos(),
									"validate tag \"dive\" can't dive on a non slice or map, %s.%s is %s type",
									structName, field.Names[0], field.Type)
							}
						}

						// can't use comma as separator in oneof tag
						if /*(args.Option == "default" || args.Option == args.RequireOneof) &&*/ strings.Contains(field.Tag.Value, "oneof=") {
							//fmt.Println("==>", field.Tag.Value)
							tagValue := strings.Split(field.Tag.Value, "oneof=")[1]
							//fmt.Println("==>", tagValue)
							if strings.Contains(tagValue, ",") {
								pass.Reportf(field.Pos(),
									"validate tag \"oneof\" can't use comma as separator, %s.%s has invalid tag value: %s",
									structName, field.Names[0], tagValue)
							}
						}

					}

				}
			}
		}
		return true
	}
	for _, file := range pass.Files {
		ast.Inspect(file, inspect)
	}
	return nil, nil
}

func isSliceOrMapType(expr ast.Expr) bool {
	switch expr.(type) {
	case *ast.ArrayType:
		return true
	case *ast.MapType:
		return true
	default:
		return false
	}
}

func hasTestForStruct(pass *analysis.Pass, structName string) bool {
	for _, file := range pass.Files {
		if strings.HasSuffix(pass.Fset.File(file.Pos()).Name(), "_test.go") {
			for _, decl := range file.Decls {
				if fn, ok := decl.(*ast.FuncDecl); ok {
					if strings.Contains(fn.Name.Name, structName) {
						return true
					}
				}
			}
		}
	}
	return false
}
