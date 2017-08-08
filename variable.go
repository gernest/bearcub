package bearcub

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/pytlesk4/m"
)

// TODO: implement replace variables functions.
// Variables are valid json, they can be either an array or object.

const (
	open    = '{'
	closing = '}'
)

func dummy(a string) string {
	return a
}

type JSONReplacer struct {
	O map[string]interface{}
}

func NewJSONReplacer(src []byte) (*JSONReplacer, error) {
	j := &JSONReplacer{O: make(map[string]interface{})}
	err := json.Unmarshal(src, &j.O)
	if err != nil {
		return nil, err
	}
	return j, nil
}

func (j *JSONReplacer) Replace(a string) string {
	a = strings.TrimSpace(a)
	if v, ok := m.GetOK(j.O, a); ok {
		return fmt.Sprint(v)
	}
	return a
}

func ReplaceString(out io.Writer, src []byte, r func(string) string) error {
	rd := bytes.NewReader(src)
	var isOpen bool
	o := &bytes.Buffer{}
	buf := &bytes.Buffer{}
	for {
		ch, _, err := rd.ReadRune()
		if err != nil {
			if err != io.EOF {
				return err
			}
			break
		}
		switch ch {
		case open:
			isOpen = true
		case closing:
			isOpen = false
			s := buf.String()
			e := r(s)
			if s == e {
				o.WriteRune(open)
				o.WriteString(s)
				o.WriteRune(closing)
			} else {
				o.WriteString(e)
			}
			buf.Reset()
			continue
		default:
			if isOpen {
				buf.WriteRune(ch)
			} else {
				o.WriteRune(ch)
			}
		}
	}
	out.Write(o.Bytes())
	return nil
}

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
