package vdf

import (
	"encoding/json"
	"fmt"
	"io"
	"unicode"
)

// NewDecoder returns a Decoder reading from r.
func NewDecoder(r io.Reader) *Decoder {
	return &Decoder{r: r}
}

// Decoder implements the same interface as encoding/json
// for decoding vdf (Valve data format). It does so in
// an _extremely_ ugly way: "massaging" vdf into JSON
// and passing it to an encoding/json.Decoder. Oops!
type Decoder struct {
	r io.Reader
}

// Decode reads the next VDF-encoded value from its input
// and stores it in the value pointed to by v using
// encoding/json tags.
func (d *Decoder) Decode(v any) error {
	pr, pw := io.Pipe()
	// buf := new(bytes.Buffer)
	// mw := io.MultiWriter(buf, pw)
	// defer func() {
	// 	fmt.Println(buf.String(), d.depth)
	// }()

	// Massage vdf into json to use the encoding/json package.
	go func() {
		defer pw.Close()
		var (
			depth            int
			eof              bool
			lineDoubleQuotes int
		)

		for {
			var p [32768]byte

			n, err := d.r.Read(p[:])
			if err != nil {
				pw.CloseWithError(err)
				return
			}

			for i, c := range p[:n] {
				var r = []byte{c}

				// TODO: This breaks if the buffer doesn't have
				// enough bytes to be able to look ahead to
				// determine if a comma is needed. To work around this
				// for now, we just make the buffer enormous, as the
				// smaller it is, the more "breakpoints" there are in
				// the read and the more likely one of those breakpoints
				// is to land on a spot that needs to be looked past
				// to determine if a comma is needed.
				commaIfNeeded := func() {
					// Look ahead to see if we're at the end of an object
					// to determine if we need a trailing comma.
					for _, q := range p[i+1:] {
						if q == '}' {
							break
						} else if unicode.IsSpace(rune(q)) {
						} else {
							r = append(r, ',')
							break
						}
					}
				}

				switch c {
				case '{':
					depth = depth + 1
				case '}':
					if depth <= 1 {
						// steamcmd's terminal doens't EOF
						// when it's done printing the vdf,
						// so we artificially EOF when we
						// know that the reader is out of
						// relevant data.
						eof = true
					} else {
						commaIfNeeded()
					}

					depth = depth - 1
				case '"':
					lineDoubleQuotes = lineDoubleQuotes + 1
					switch lineDoubleQuotes {
					case 2:
						r = append(r, ':')
					case 4:
						commaIfNeeded()
					case 1:
					case 3:
					default:
						pw.CloseWithError(fmt.Errorf(`vdf: string contained "`))
						return
					}
				case '\n':
					lineDoubleQuotes = 0
				}

				if depth > 0 {
					if _, err := pw.Write(r); err != nil {
						pw.CloseWithError(err)
						return
					}
				}

				if eof {
					return
				}
			}
		}
	}()

	return json.NewDecoder(pr).Decode(v)
}
