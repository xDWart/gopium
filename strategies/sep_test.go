package strategies

import (
	"context"
	"reflect"
	"testing"

	"github.com/1pkg/gopium/collections"
	"github.com/1pkg/gopium/gopium"
	"github.com/1pkg/gopium/tests/mocks"
)

func TestSep(t *testing.T) {
	// prepare
	cctx, cancel := context.WithCancel(context.Background())
	cancel()
	table := map[string]struct {
		sep sep
		c   gopium.Curator
		ctx context.Context
		o   gopium.Struct
		r   gopium.Struct
		err error
	}{
		"empty struct should be applied to empty struct": {
			sep: sepl1b,
			c:   mocks.Maven{SCache: []int64{32}},
			ctx: context.Background(),
		},
		"non empty struct should be applied to expected aligned struct": {
			sep: sepl2b,
			c:   mocks.Maven{SCache: []int64{16, 16, 16}},
			ctx: context.Background(),
			o: gopium.Struct{
				Name: "test",
				Fields: []gopium.Field{
					{
						Name: "test",
						Size: 8,
					},
				},
			},
			r: gopium.Struct{
				Name: "test",
				Fields: []gopium.Field{
					{
						Name: "test",
						Size: 8,
					},
					collections.PadField(16),
				},
			},
		},
		"non empty struct should be applied to expected aligned struct custom bytes": {
			sep: sepbb.Bytes(12),
			c:   mocks.Maven{SCache: []int64{16, 16, 16}},
			ctx: context.Background(),
			o: gopium.Struct{
				Name: "test",
				Fields: []gopium.Field{
					{
						Name: "test",
						Size: 8,
					},
				},
			},
			r: gopium.Struct{
				Name: "test",
				Fields: []gopium.Field{
					{
						Name: "test",
						Size: 8,
					},
					collections.PadField(12),
				},
			},
		},
		"non empty struct should be applied to expected aligned struct custom bytes and cache line": {
			sep: sepl2b.Bytes(12),
			c:   mocks.Maven{SCache: []int64{16, 16, 16}},
			ctx: context.Background(),
			o: gopium.Struct{
				Name: "test",
				Fields: []gopium.Field{
					{
						Name: "test",
						Size: 8,
					},
				},
			},
			r: gopium.Struct{
				Name: "test",
				Fields: []gopium.Field{
					{
						Name: "test",
						Size: 8,
					},
					collections.PadField(16),
				},
			},
		},
		"non empty struct should be applied to expected aligned struct on canceled context": {
			sep: sepl3b,
			c:   mocks.Maven{SCache: []int64{16, 16, 16}},
			ctx: cctx,
			o: gopium.Struct{
				Name: "test",
				Fields: []gopium.Field{
					{
						Name: "test",
						Size: 8,
					},
				},
			},
			r: gopium.Struct{
				Name: "test",
				Fields: []gopium.Field{
					{
						Name: "test",
						Size: 8,
					},
					collections.PadField(16),
				},
			},
			err: context.Canceled,
		},
		"mixed struct should be applied to expected aligned struct with sys top": {
			sep: sepsyst,
			c:   mocks.Maven{SAlign: 24, SCache: []int64{16, 32, 64}},
			ctx: context.Background(),
			o: gopium.Struct{
				Name: "test",
				Fields: []gopium.Field{
					{
						Name: "test1",
						Size: 32,
					},
					{
						Name: "test2",
						Size: 8,
					},
					{
						Name: "test3",
						Size: 8,
					},
					{
						Name: "test4",
						Size: 3,
					},
				},
			},
			r: gopium.Struct{
				Name: "test",
				Fields: []gopium.Field{
					collections.PadField(24),
					{
						Name: "test1",
						Size: 32,
					},
					{
						Name: "test2",
						Size: 8,
					},
					{
						Name: "test3",
						Size: 8,
					},
					{
						Name: "test4",
						Size: 3,
					},
				},
			},
		},
		"mixed struct should be applied to expected aligned struct with sys top cutom bytes": {
			sep: sepsyst.Bytes(100),
			c:   mocks.Maven{SAlign: 24, SCache: []int64{16, 32, 64}},
			ctx: context.Background(),
			o: gopium.Struct{
				Name: "test",
				Fields: []gopium.Field{
					{
						Name: "test1",
						Size: 32,
					},
					{
						Name: "test2",
						Size: 8,
					},
					{
						Name: "test3",
						Size: 8,
					},
					{
						Name: "test4",
						Size: 3,
					},
				},
			},
			r: gopium.Struct{
				Name: "test",
				Fields: []gopium.Field{
					collections.PadField(24),
					{
						Name: "test1",
						Size: 32,
					},
					{
						Name: "test2",
						Size: 8,
					},
					{
						Name: "test3",
						Size: 8,
					},
					{
						Name: "test4",
						Size: 3,
					},
				},
			},
		},
		"mixed struct should be applied to expected aligned struct with sys bottom": {
			sep: sepsysb,
			c:   mocks.Maven{SAlign: 24, SCache: []int64{16, 32, 64}},
			ctx: context.Background(),
			o: gopium.Struct{
				Name: "test",
				Fields: []gopium.Field{
					{
						Name: "test1",
						Size: 32,
					},
					{
						Name: "test2",
						Size: 8,
					},
					{
						Name: "test3",
						Size: 8,
					},
					{
						Name: "test4",
						Size: 3,
					},
				},
			},
			r: gopium.Struct{
				Name: "test",
				Fields: []gopium.Field{
					{
						Name: "test1",
						Size: 32,
					},
					{
						Name: "test2",
						Size: 8,
					},
					{
						Name: "test3",
						Size: 8,
					},
					{
						Name: "test4",
						Size: 3,
					},
					collections.PadField(24),
				},
			},
		},
		"mixed struct should be applied to expected aligned struct with cache top": {
			sep: sepl3t,
			c:   mocks.Maven{SAlign: 24, SCache: []int64{16, 32, 64}},
			ctx: context.Background(),
			o: gopium.Struct{
				Name: "test",
				Fields: []gopium.Field{
					{
						Name: "test1",
						Size: 32,
					},
					{
						Name: "test2",
						Size: 8,
					},
					{
						Name: "test3",
						Size: 8,
					},
					{
						Name: "test4",
						Size: 3,
					},
					{
						Name: "test5",
						Size: 1,
					},
				},
			},
			r: gopium.Struct{
				Name: "test",
				Fields: []gopium.Field{
					collections.PadField(64),
					{
						Name: "test1",
						Size: 32,
					},
					{
						Name: "test2",
						Size: 8,
					},
					{
						Name: "test3",
						Size: 8,
					},
					{
						Name: "test4",
						Size: 3,
					},
					{
						Name: "test5",
						Size: 1,
					},
				},
			},
		},
	}
	for name, tcase := range table {
		t.Run(name, func(t *testing.T) {
			// prepare
			sep := tcase.sep.Curator(tcase.c)
			// exec
			r, err := sep.Apply(tcase.ctx, tcase.o)
			// check
			if !reflect.DeepEqual(r, tcase.r) {
				t.Errorf("actual %v doesn't equal to expected %v", r, tcase.r)
			}
			if !reflect.DeepEqual(err, tcase.err) {
				t.Errorf("actual %v doesn't equal to expected %v", err, tcase.err)
			}
		})
	}
}
