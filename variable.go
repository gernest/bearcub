// Package bearcub provides simple API for a simple mustache like string
// templating
package bearcub

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"strings"
	"text/template"
	"unicode"

	"github.com/pytlesk4/m"
)

const (
	open    = '{'
	closing = '}'
)

func dummy(a string) string {
	return "{" + a + "}"
}

// JSONReplacer uses a map of json decoded values i.e map[string]interface{} to
// offer retriving string representation of set values
type JSONReplacer struct {
	O map[string]interface{}
}

// NewJSONReplacer returns a new instance of JSONReplacer. src must be a valid
// json string object.
func NewJSONReplacer(src []byte) (*JSONReplacer, error) {
	j := &JSONReplacer{O: make(map[string]interface{})}
	err := json.Unmarshal(src, &j.O)
	if err != nil {
		return nil, err
	}
	return j, nil
}

// Replace returns a string representation of the value with key a.
func (j *JSONReplacer) Replace(a string) string {
	a = strings.TrimSpace(a)
	if v, ok := m.GetOK(j.O, a); ok {
		return fmt.Sprint(v)
	}
	return "{" + a + "}"
}

// ReplaceString replaces any occurrence of keys inside { } from the src with
// values returned by the replacer function r
func ReplaceString(out io.Writer, src []byte, r func(string) string) error {
	funcs := make(template.FuncMap)
	funcs["replace"] = r
	var o bytes.Buffer
	err := Rewrite(&o, string(src))
	if err != nil {
		return err
	}
	tpl, err := template.New("replacing").Funcs(funcs).Delims("<<", ">>").Parse(
		o.String(),
	)
	if err != nil {
		return err
	}
	return tpl.Execute(out, nil)
}

// Replace replace any string templates in the request objects with variables.
func Replace(req *http.Request, variables string) error {
	var r func(string) string
	if variables == "" {
		r = dummy
	} else {
		jr, err := NewJSONReplacer([]byte(variables))
		if err != nil {
			return err
		}
		r = jr.Replace
	}
	var o bytes.Buffer
	rep := func(a string) string {
		if hasCurly(a) {
			ReplaceString(&o, []byte(a), r)
			a = o.String()
			o.Reset()
		}
		return a
	}
	req.URL.Host = rep(req.URL.Host)
	req.Host = ""
	req.URL.Path = rep(req.URL.Path)
	req.URL.RawQuery = rep(req.URL.RawQuery)
	if req.Body != nil {
		b, err := ioutil.ReadAll(req.Body)
		if err != nil {
			return err
		}
		req.Body = ioutil.NopCloser(strings.NewReader(rep(string(b))))
	}
	return nil
}

func hasCurly(src string) bool {
	return strings.Contains(src, "{")
}

func Rewrite(out io.Writer, src string) error {
	r := strings.NewReader(src)
	o := &bytes.Buffer{}
	for {
		ch, _, err := r.ReadRune()
		if err != nil {
			if err == io.EOF {
				out.Write(o.Bytes())
				return nil
			}
			return err
		}
		switch ch {
		case open:
			if err = consumeSpace(r); err != nil {
				return err
			}
			next, _, err := r.ReadRune()
			if err != nil {
				return err
			}
			if unicode.IsLetter(next) {
				k, err := readUntil(r, closing)
				if err != nil {
					return err
				}
				k = string(next) + k
				o.WriteString(attachReplaceFunc(k))
			} else {
				o.WriteRune(open)
				o.WriteRune(next)
			}
		default:
			o.WriteRune(ch)
		}
	}
}

func attachReplaceFunc(src string) string {
	return fmt.Sprintf("<<replace \"%s\">>", src)
}

func consumeSpace(r *strings.Reader) error {
	for {
		ch, _, err := r.ReadRune()
		if err != nil {
			return err
		}
		if !unicode.IsSpace(ch) {
			r.UnreadRune()
			return err
		}
	}
}

func readUntil(r *strings.Reader, ch rune) (string, error) {
	var o bytes.Buffer
	for {
		n, _, err := r.ReadRune()
		if err != nil {
			return "", err
		}
		if n == ch {
			return o.String(), nil
		}
		o.WriteRune(n)
	}
}
