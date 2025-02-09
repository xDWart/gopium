package fmtio

import (
	"context"
	"go/ast"
	"go/format"
	"go/token"
	"io"
)

// Gofmt implements printer by
// using canonical ast go fmt printer
type Gofmt struct{} // struct size: 0 bytes; struct align: 1 bytes; struct aligned size: 0 bytes; - 🌺 gopium @1pkg

// Print gofmt implementation
func (p Gofmt) Print(ctx context.Context, w io.Writer, fset *token.FileSet, node ast.Node) error {
	// manage context actions
	// in case of cancelation
	// stop execution
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
	}
	// use gofmt node
	return format.Node(w, fset, node)
}
