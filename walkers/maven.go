package walkers

import (
	"go/types"
	"sync"

	"github.com/1pkg/gopium/collections"
	"github.com/1pkg/gopium/gopium"
)

// sizealign defines data transfer
// object that holds type pair
// of size and align vals
type sizealign struct {
	size  int64 `gopium:"filter_pads,memory_pack,cache_rounding_cpu_l1_discrete,struct_annotate_comment,add_tag_group_force"`
	align int64 `gopium:"filter_pads,memory_pack,cache_rounding_cpu_l1_discrete,struct_annotate_comment,add_tag_group_force"`
} // struct size: 16 bytes; struct align: 8 bytes; struct aligned size: 16 bytes; - 🌺 gopium @1pkg

// maven defines visiting helper
// that aggregates some useful
// operations on underlying facilities
type maven struct {
	store sync.Map               `gopium:"filter_pads,memory_pack,cache_rounding_cpu_l1_discrete,struct_annotate_comment,add_tag_group_force"`
	exp   gopium.Exposer         `gopium:"filter_pads,memory_pack,cache_rounding_cpu_l1_discrete,struct_annotate_comment,add_tag_group_force"`
	loc   gopium.Locator         `gopium:"filter_pads,memory_pack,cache_rounding_cpu_l1_discrete,struct_annotate_comment,add_tag_group_force"`
	ref   *collections.Reference `gopium:"filter_pads,memory_pack,cache_rounding_cpu_l1_discrete,struct_annotate_comment,add_tag_group_force"`
	_     [48]byte               `gopium:"filter_pads,memory_pack,cache_rounding_cpu_l1_discrete,struct_annotate_comment,add_tag_group_force"`
} // struct size: 128 bytes; struct align: 8 bytes; struct aligned size: 128 bytes; - 🌺 gopium @1pkg

// has defines struct store id helper
// that uses locator to build id
// for a structure and check that
// builded id has not been stored already
func (m *maven) has(tn *types.TypeName) (id string, loc string, ok bool) {
	// build id for the structure
	id = m.loc.ID(tn.Pos())
	// build loc for the structure
	loc = m.loc.Loc(tn.Pos())
	// in case id of structure
	// has been already stored
	if _, ok := m.store.Load(id); ok {
		return id, loc, true
	}
	// mark id of structure as stored
	m.store.Store(id, struct{}{})
	return id, loc, false
}

// enum defines struct enumerating converting helper
// that goes through all structure fields
// and uses exposer to expose field DTO
// for each field and puts them back
// to resulted struct object
func (m *maven) enum(name string, st *types.Struct) gopium.Struct {
	// set structure name
	r := gopium.Struct{}
	r.Name = name
	// get number of struct fields
	nf := st.NumFields()
	// prefill Fields
	r.Fields = make([]gopium.Field, 0, nf)
	for i := 0; i < nf; i++ {
		// get field
		f := st.Field(i)
		// get size and align for field
		sa := m.refsa(f.Type())
		// fill field structure
		r.Fields = append(r.Fields, gopium.Field{
			Name:     f.Name(),
			Type:     m.exp.Name(f.Type()),
			Size:     sa.size,
			Align:    sa.align,
			Tag:      st.Tag(i),
			Exported: f.Exported(),
			Embedded: f.Embedded(),
		})
	}
	return r
}

// refsa defines size and align getter
// with reference helper that uses reference
// if it has been provided
// or uses exposer to expose type size
func (m *maven) refsa(t types.Type) sizealign {
	// in case we don't have a reference
	// just use default exposer size
	if m.ref == nil {
		return sizealign{
			size:  m.exp.Size(t),
			align: m.exp.Align(t),
		}
	}
	// for refsize only named structures
	// and arrays should be calculated
	// not with default exposer size
	switch tp := t.(type) {
	case *types.Array:
		// note: copied from `go/types/sizes.go`
		n := tp.Len()
		if n <= 0 {
			return sizealign{}
		}
		// n > 0
		sa := m.refsa(tp.Elem())
		sa.size = collections.Align(sa.size, sa.align)*(n-1) + sa.size
		return sa
	case *types.Named:
		// in case it's not a struct skip it
		if _, ok := tp.Underlying().(*types.Struct); !ok {
			break
		}
		// get id for named structures
		id := m.loc.ID(tp.Obj().Pos())
		// get size of the structure from ref
		if sa, ok := m.ref.Get(id).(sizealign); ok {
			return sa
		}
	}
	// just use default exposer size
	return sizealign{
		size:  m.exp.Size(t),
		align: m.exp.Align(t),
	}
}

// refst helps to create struct
// size refence for provided key
// by preallocating the key and then
// pushing total struct size to ref with closure
func (m *maven) refst(name string) func(gopium.Struct) {
	// preallocate the key
	m.ref.Alloc(name)
	// return the pushing closure
	return func(st gopium.Struct) {
		// calculate structure align and aligned size
		stsize, stalign := collections.SizeAlign(st)
		// set ref key size and align
		m.ref.Set(name, sizealign{size: stsize, align: stalign})
	}
}
