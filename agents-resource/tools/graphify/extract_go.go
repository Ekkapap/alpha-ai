package main

// extract_go.go — Go AST extraction using go/ast stdlib.

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"path/filepath"
	"strings"
)

func extractGo(path, rel string) (ExtractedFile, error) {
	ef := ExtractedFile{RelPath: rel, Language: "go"}

	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, path, nil, 0)
	if err != nil {
		return ef, err
	}

	fileLabel := fileNodeLabel(rel)
	ef.Nodes = append(ef.Nodes, RawNode{Label: fileLabel, Location: "L1", Kind: "file"})

	// Track declared function names for call resolution
	declaredFuncs := map[string]bool{}

	ast.Inspect(f, func(n ast.Node) bool {
		switch v := n.(type) {

		case *ast.FuncDecl:
			name := v.Name.Name
			if v.Recv != nil && len(v.Recv.List) > 0 {
				recv := receiverType(v.Recv.List[0].Type)
				name = recv + "." + name
			}
			label := name + "()"
			loc := fmt.Sprintf("L%d", fset.Position(v.Pos()).Line)
			ef.Nodes = append(ef.Nodes, RawNode{Label: label, Location: loc, Kind: "func"})
			ef.Edges = append(ef.Edges, RawEdge{FromLabel: fileLabel, ToLabel: label, Relation: "contains", Location: loc})
			declaredFuncs[name] = true

		case *ast.TypeSpec:
			loc := fmt.Sprintf("L%d", fset.Position(v.Pos()).Line)
			kind := "type"
			if _, ok := v.Type.(*ast.InterfaceType); ok {
				kind = "interface"
			} else if _, ok := v.Type.(*ast.StructType); ok {
				kind = "struct"
			}
			ef.Nodes = append(ef.Nodes, RawNode{Label: v.Name.Name, Location: loc, Kind: kind})
			ef.Edges = append(ef.Edges, RawEdge{FromLabel: fileLabel, ToLabel: v.Name.Name, Relation: "contains", Location: loc})
		}
		return true
	})

	// Second pass: collect call edges
	ast.Inspect(f, func(n ast.Node) bool {
		fd, ok := n.(*ast.FuncDecl)
		if !ok {
			return true
		}
		callerName := fd.Name.Name
		if fd.Recv != nil && len(fd.Recv.List) > 0 {
			callerName = receiverType(fd.Recv.List[0].Type) + "." + callerName
		}
		callerLabel := callerName + "()"

		ast.Inspect(fd.Body, func(inner ast.Node) bool {
			call, ok := inner.(*ast.CallExpr)
			if !ok {
				return true
			}
			callee := calleeLabel(call)
			if callee == "" || callee == callerLabel {
				return true
			}
			loc := fmt.Sprintf("L%d", fset.Position(call.Pos()).Line)
			ef.Edges = append(ef.Edges, RawEdge{
				FromLabel: callerLabel,
				ToLabel:   callee,
				Relation:  "calls",
				Location:  loc,
			})
			return true
		})
		return true
	})

	// Import references
	for _, imp := range f.Imports {
		if imp.Path == nil {
			continue
		}
		importPath := strings.Trim(imp.Path.Value, `"`)
		parts := strings.Split(importPath, "/")
		pkgName := parts[len(parts)-1]
		loc := fmt.Sprintf("L%d", fset.Position(imp.Pos()).Line)
		ef.Edges = append(ef.Edges, RawEdge{
			FromLabel: fileLabel,
			ToLabel:   pkgName,
			Relation:  "references",
			Location:  loc,
		})
	}

	return ef, nil
}

func receiverType(expr ast.Expr) string {
	switch t := expr.(type) {
	case *ast.StarExpr:
		return receiverType(t.X)
	case *ast.Ident:
		return t.Name
	case *ast.IndexExpr:
		return receiverType(t.X)
	}
	return "Unknown"
}

func calleeLabel(call *ast.CallExpr) string {
	switch fn := call.Fun.(type) {
	case *ast.Ident:
		if fn.Name == "" {
			return ""
		}
		return fn.Name + "()"
	case *ast.SelectorExpr:
		return fn.Sel.Name + "()"
	}
	return ""
}

func fileNodeLabel(rel string) string {
	// e.g. "alpha/main.go" → "main.go"
	return strings.TrimPrefix(rel, strings.TrimSuffix(rel, "/"+lastSegment(rel))+"/")
}

func lastSegment(rel string) string {
	parts := strings.Split(filepath.ToSlash(rel), "/")
	return parts[len(parts)-1]
}
