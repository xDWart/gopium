package walkers

import (
	"context"
	"go/ast"
	"go/printer"
	"regexp"

	"1pkg/gopium"
	"1pkg/gopium/astutil"
	"1pkg/gopium/collections"
	"1pkg/gopium/fmtio"

	"golang.org/x/sync/errgroup"
)

// list of wast presets
var (
	fsptnstd = wast{
		apply: astutil.Sync,
		wgen:  fmtio.Stdout,
	}
	fsptngo = wast{
		apply: astutil.Sync,
		wgen:  fmtio.FileGo,
	}
	fsptngopium = wast{
		apply: astutil.Sync,
		wgen:  fmtio.FileGopium,
	}
)

// wast defines packages walker sync ast implementation
// that uses pkgs.Parser to parse packages types data
// astutil to update ast to results from strategy
type wast struct {
	parser  gopium.Parser
	exposer gopium.Exposer
	apply   astutil.Apply
	wgen    fmtio.WriterGen
	backref bool
}

// With erich wast walker with parser, exposer, and ref instance
func (w wast) With(parser gopium.Parser, exposer gopium.Exposer, backref bool) wast {
	w.parser = parser
	w.exposer = exposer
	w.backref = backref
	return w
}

// Visit wast implementation
func (w wast) Visit(ctx context.Context, regex *regexp.Regexp, stg gopium.Strategy, deep bool) error {
	// uses gopium.Visit and gopium.VisitFunc helpers
	// to go through all structs decls inside the package
	// and apply strategy to them to get results
	// then overrides os.Files with updated ast
	// builded from strategy results

	// use parser to parse types pkg data
	// we don't care about fset
	pkg, loc, err := w.parser.ParseTypes(ctx)
	if err != nil {
		return err
	}
	// create govisit func
	// using visit helper
	// and run it on pkg scope
	ch := make(appliedCh)
	gvisit := visit(
		regex,
		stg,
		w.exposer,
		loc,
		ch,
		deep,
		w.backref,
	)
	// run visiting in separate goroutine
	go gvisit(ctx, pkg.Scope())
	// prepare struct storage
	h := make(collections.Hierarchic)
	for applied := range ch {
		// in case any error happened just return error
		// it will cancel context automatically
		if applied.Error != nil {
			return applied.Error
		}
		// push struct to storage
		h.Push(applied.ID, applied.Loc, applied.Result)
	}
	// run sync write
	// with collected strategies results
	return w.write(ctx, h)
}

// write wast helps apply
// sync and persist to format strategy results
// by updating os.Files
func (w wast) write(ctx context.Context, h collections.Hierarchic) error {
	// use parser to parse ast pkg data
	pkg, loc, err := w.parser.ParseAst(ctx)
	if err != nil {
		return err
	}
	// run ast apply with strategy result
	// to update ast.Package on the parsed ast.Package
	// in case any error happened just return error back
	pkg, err = w.apply(ctx, pkg, loc, h)
	if err != nil {
		return err
	}
	// run async persist helper
	return w.persist(ctx, pkg, loc)
}

// persist wast helps to update os.Files
// accordingly to updated ast.Package
// concurently or return error otherwise
func (w wast) persist(ctx context.Context, pkg *ast.Package, loc gopium.Locator) error {
	// create sync error group
	// with cancelation context
	group, gctx := errgroup.WithContext(ctx)
loop:
	// go through all files in package
	// and update them to concurently
	for name, file := range pkg.Files {
		// manage context actions
		// in case of cancelation
		// stop execution
		select {
		case <-gctx.Done():
			break loop
		default:
		}
		// capture name and file copies
		name := name
		file := file
		// run error group write call
		group.Go(func() error {
			// manage context actions
			// in case of cancelation
			// stop execution and return error
			select {
			case <-gctx.Done():
				return gctx.Err()
			default:
			}
			// generate relevant writer
			writer, err := w.wgen(name, loc.Loc(file.Pos()))
			// in case any error happened just return error
			// it will cancel context automatically
			if err != nil {
				return err
			}
			// grab the latest file fset
			fset, _ := loc.Fset(file.Pos(), nil)
			// write updated ast.File to related os.File
			// use original toke.FileSet to keep format
			// in case any error happened just return error
			// it will cancel context automatically
			return printer.Fprint(
				writer,
				fset,
				file,
			)
		})
	}
	// wait until all writers
	// resolve their jobs and
	return group.Wait()
}
