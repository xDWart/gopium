package fmtio

import (
	"bytes"
	"encoding/csv"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io"
	"io/ioutil"
	"strconv"
	"strings"

	"github.com/1pkg/gopium/gopium"
)

// Jsonb defines bytes implementation
// which uses json marshal with indent
// to serialize flat collection to byte slice
func Jsonb(sts []gopium.Struct) ([]byte, error) {
	// just use json marshal with indent
	return json.MarshalIndent(sts, "", "\t")
}

// Xmlb defines bytes implementation
// which uses xml marshal with indent
// to serialize flat collection to byte slice
func Xmlb(sts []gopium.Struct) ([]byte, error) {
	// just use xml marshal with indent
	return xml.MarshalIndent(sts, "", "\t")
}

// Csvb defines bytes implementation
// which serializes flat collection
// to formatted csv byte slice
func Csvb(rw io.ReadWriter) gopium.Bytes {
	return func(sts []gopium.Struct) ([]byte, error) {
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
		for _, st := range sts {
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

// Mdtb defines bytes implementation
// which serializes flat collection
// to formatted markdown table byte slice
func Mdtb(sts []gopium.Struct) ([]byte, error) {
	// prepare buffer and collections
	var buf bytes.Buffer
	if len(sts) > 0 {
		// write header
		// no error should be
		// checked as it uses
		// buffered writer
		_, _ = buf.WriteString("| Struct Name | Struct Doc | Struct Comment | Field Name | Field Type | Field Size | Field Align | Field Tag | Field Exported | Field Embedded | Field Doc | Field Comment |\n")
		_, _ = buf.WriteString("| :---: | :---: | :---: | :---: | :---: | :---: | :---: | :---: | :---: | :---: | :---: | :---: |\n")
		for _, st := range sts {
			// go through all fields
			// and write then one by one
			for _, f := range st.Fields {
				// no error should be
				// checked as it uses
				// buffered writer
				_, _ = buf.WriteString(
					fmt.Sprintf("| %s | %s | %s | %s | %s | %d | %d | %s | %s | %s | %s | %s |\n",
						st.Name,
						strings.Join(st.Doc, " "),
						strings.Join(st.Comment, " "),
						f.Name,
						f.Type,
						f.Size,
						f.Align,
						f.Tag,
						strconv.FormatBool(f.Exported),
						strconv.FormatBool(f.Embedded),
						strings.Join(f.Doc, " "),
						strings.Join(f.Comment, " "),
					),
				)
			}
		}
	}
	return buf.Bytes(), nil
}
