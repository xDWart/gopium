package fmtio

import (
	"encoding/csv"
	"encoding/json"
	"encoding/xml"
	"io"
	"io/ioutil"
	"strconv"
	"strings"

	"1pkg/gopium/collections"
)

// Bytes defines abstraction for formatting
// gopium flat collection to byte slice
type Bytes func(collections.Flat) ([]byte, error)

// Jsonb defines bytes implementation
// which uses json marshal with indent
// to serialize flat collection to byte slice
func Jsonb(f collections.Flat) ([]byte, error) {
	// just use json marshal with indent
	return json.MarshalIndent(f.Sorted(), "", "\t")
}

// Xmlb defines bytes implementation
// which uses xml marshal with indent
// to serialize flat collection to byte slice
func Xmlb(f collections.Flat) ([]byte, error) {
	// just use xml marshal with indent
	return xml.MarshalIndent(f.Sorted(), "", "\t")
}

// Csvb defines bytes implementation
// which serializes flat collection
// to formatted csv byte slice
func Csvb(rw io.ReadWriter) Bytes {
	return func(f collections.Flat) ([]byte, error) {
		// prepare csv writer
		w := csv.NewWriter(rw)
		// write header
		// no error should be
		// checked as it uses
		// buffered writer
		_ = w.Write([]string{
			"Struct Name",
			"Struct Doc",
			"Struct Comment",
			"Field Name",
			"Field Type",
			"Field Size",
			"Field Align",
			"Field Tag",
			"Field Exported",
			"Field Embedded",
			"Field Doc",
			"Field Comment",
		})
		for _, st := range f.Sorted() {
			// go through all fields
			// and write then one by one
			for _, f := range st.Fields {
				// no error should be
				// checked as it uses
				// buffered writer
				_ = w.Write([]string{
					st.Name,
					strings.Join(st.Doc, " "),
					strings.Join(st.Comment, " "),
					f.Name,
					f.Type,
					strconv.Itoa(int(f.Size)),
					strconv.Itoa(int(f.Align)),
					f.Tag,
					strconv.FormatBool(f.Exported),
					strconv.FormatBool(f.Embedded),
					strings.Join(f.Doc, " "),
					strings.Join(f.Comment, " "),
				})
			}
			// flush to buf
			w.Flush()
			// check flush error
			if err := w.Error(); err != nil {
				return nil, err
			}
		}
		// and return buf result
		return ioutil.ReadAll(rw)
	}
}
