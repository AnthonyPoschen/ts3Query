//ts3Query implements the
package ts3Query

import (
	"bytes"
	"fmt"
	"io"
	"strings"
)

var (
	escapeTable = map[byte][]byte{
		'\\': []byte("\\\\"),
		'/':  []byte("\\/"),
		' ':  []byte("\\s"),
		'|':  []byte("\\p"),
		'\a': []byte("\\a"),
		'\b': []byte("\\b"),
		'\f': []byte("\\f"),
		'\n': []byte("\\n"),
		'\r': []byte("\\r"),
		'\t': []byte("\\t"),
		'\v': []byte("\\v"),
	}
	unescapeTable = map[string]byte{
		string(escapeTable['\\']): '\\',
		string(escapeTable['/']):  '/',
		string(escapeTable[' ']):  ' ',
		string(escapeTable['|']):  '|',
		string(escapeTable['\a']): '\a',
		string(escapeTable['\b']): '\b',
		string(escapeTable['\f']): '\f',
		string(escapeTable['\n']): '\n',
		string(escapeTable['\r']): '\r',
	}
)

// Ts3Query is the main object which contains all the features
type Ts3Query struct {
	rw io.ReadWriter
	b  []byte
}

// New returns an instance of ts3Query
func New(rw io.ReadWriter) Ts3Query {
	return Ts3Query{rw: rw, b: make([]byte, 1024*1024)}
}

// escapesString escapes the input to be sent to the ts3 client.
func escapeString(input string) string {

	var result string
	r := bytes.NewReader([]byte(input))
	for {
		b, err := r.ReadByte()
		if err != nil {
			break
		}
		// check if the current byte is in the lookup table
		if val, ok := escapeTable[b]; ok {
			result += string(val)
			continue
		}
		result += string(b)
	}
	return result
}

func getError(ts3Msg string) error {
	s := strings.Split(ts3Msg, "error ")
	l := len(s)
	rawerr := strings.TrimLeft(s[l-1], " ")
	items := strings.Split(rawerr, " ")
	m := make(map[string]string)
	for _, v := range items {
		x := strings.Split(v, "=")
		m[x[0]] = x[1]
	}
	if m["id"] != "0" {
		return fmt.Errorf("%s", m["msg"])
	}
	return nil
}

func (t *Ts3Query) readResponse() (res string, err error) {

	for {
		n, err := t.rw.Read(t.b)
		if err != nil {
			break
		}
		if n != 0 {
			res = string(t.b[0:n])
			break
		}
	}
	return
}
