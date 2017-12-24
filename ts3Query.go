//ts3Query implements the
package ts3Query

import (
	"bytes"
	"fmt"
	"io"
	"regexp"
	"strings"
	"time"
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
	rw    io.ReadWriter
	b     []byte
	delay time.Duration
}

// New returns an instance of ts3Query
//
// By default teamspeak will only allow 10 commands every 3 seconds from external IP's.
// so make sure to set the delay to over 300ms to avoid being banned for flooding.
func New(rw io.ReadWriter, delay time.Duration) Ts3Query {
	return Ts3Query{rw: rw, b: make([]byte, 1024*1024), delay: delay}
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

// use what this function solves as a programming test in the future if needed :D
func unEscapeString(input string) (result string) {

	var prev string
	var skip bool
	var last string

	for i, str := range input {
		// if we just started skip doing stuff and pre load
		if i == 0 {
			prev = string(str)
			continue
		}
		// if we decided to skip then shuffle the values down and skip this iteration
		if skip {
			prev = string(str)
			last = prev
			skip = false
			continue
		}
		// loop to check for a pattern match
		for k, v := range unescapeTable {
			// combine the previous rune with the current and see if it matches the pattern
			if (prev + string(str)) == k {
				// set the previous to the unescaped version
				prev = string(v)
				// set skip so we ignore this str and it wont be saved
				skip = true
				break
			}
		}

		result += prev
		prev = string(str)
		last = string(str)
	}
	result += last
	return
}

func getError(ts3Msg string) error {
	s := strings.Split(ts3Msg, "error ")
	l := len(s)
	rawerr := strings.TrimLeft(s[l-1], " ")
	items := strings.Split(rawerr, " ")
	m := make(map[string]string)
	for _, v := range items {
		x := strings.Split(v, "=")
		if len(x) >= 2 {
			m[x[0]] = x[1]
		}
	}
	if m["id"] != "0" {
		return fmt.Errorf("ID:%s msg:%s", m["id"], unEscapeString(m["msg"]))
	}
	return nil
}

func (t *Ts3Query) readResponse() (res string, err error) {

	reg := regexp.MustCompile("error id=[0-9]* msg=")
	for {

		n, err := t.rw.Read(t.b)
		if err != nil {
			break
		}
		if n != 0 {
			res += string(t.b[0:n])
			if reg.Match([]byte(res)) {
				break
			}
		}
	}
	return
}

func (t *Ts3Query) sendMessage(msg string) error {
	<-time.After(t.delay)
	_, err := t.rw.Write([]byte(msg + "\n"))
	return err
}
