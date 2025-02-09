package strategies

import (
	"context"
	"fmt"
	"reflect"
	"strings"

	"github.com/1pkg/gopium/collections"
	"github.com/1pkg/gopium/gopium"
)

// list of tag presets
var (
	tags  = tag{force: false, discrete: false}
	tagf  = tag{force: true, discrete: false}
	tagsd = tag{force: false, discrete: true}
	tagfd = tag{force: true, discrete: true}
)

// tag defines strategy implementation
// that adds or updates fields tags annotation
// that could be processed by group strategy
type tag struct {
	tag      string   `gopium:"filter_pads,memory_pack,cache_rounding_cpu_l1_discrete,struct_annotate_comment,add_tag_group_force"`
	force    bool     `gopium:"filter_pads,memory_pack,cache_rounding_cpu_l1_discrete,struct_annotate_comment,add_tag_group_force"`
	discrete bool     `gopium:"filter_pads,memory_pack,cache_rounding_cpu_l1_discrete,struct_annotate_comment,add_tag_group_force"`
	_        [14]byte `gopium:"filter_pads,memory_pack,cache_rounding_cpu_l1_discrete,struct_annotate_comment,add_tag_group_force"`
} // struct size: 32 bytes; struct align: 8 bytes; struct aligned size: 32 bytes; - 🌺 gopium @1pkg

// Names erich tag strategy with strategy names tag
func (stg tag) Names(names ...gopium.StrategyName) tag {
	// convert strategies names to strings
	strs := make([]string, 0, len(names))
	for _, name := range names {
		strs = append(strs, string(name))
	}
	// concat them by comma
	stg.tag = strings.Join(strs, ",")
	return stg
}

// Apply tag implementation
func (stg tag) Apply(ctx context.Context, o gopium.Struct) (gopium.Struct, error) {
	// copy original structure to result
	r := collections.CopyStruct(o)
	// iterate through all fields
	for i := range r.Fields {
		f := &r.Fields[i]
		// grab the field tag
		tag, ok := reflect.StructTag(f.Tag).Lookup(gopium.NAME)
		// build group tag
		gtag := stg.tag
		// if we want to build discrete groups
		if stg.discrete {
			// use default group tag
			// and append index of field to it
			group := fmt.Sprintf("%s-%d", gopium.NAME, i+1)
			gtag = fmt.Sprintf("group:%s;%s", group, stg.tag)
		}
		// in case gopium tag already exists
		// and force is set - replace tag
		// in case gopium tag already exists
		// and force isn't set - do nothing
		// in case tag is not empty and
		// gopium tag doesn't exist - append tag
		// in case tag is empty - set tag
		fulltag := fmt.Sprintf(`%s:"%s"`, gopium.NAME, gtag)
		switch {
		case ok && stg.force:
			f.Tag = strings.Replace(f.Tag, tag, gtag, 1)
		case ok:
			break
		case f.Tag != "":
			f.Tag += " " + fulltag
		default:
			f.Tag = fulltag
		}
	}
	return r, ctx.Err()
}
