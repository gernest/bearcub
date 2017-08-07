package bearcub_test

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httputil"
	"net/url"
	"testing"

	"github.com/gernest/bearcub"
)

/*

TODO: Make all tests pass below. Good Luck!

*/

type variableReplaceTest struct {
	// Description of test
	Description string
	// variables to use for replacing.
	Variables string
	// request to replace varibales in.
	Req *http.Request
	// Optional []byte or func() io.ReadCloser to populate Req.Body
	Body interface{}

	// Expected result
	Expected string
}

var reqReplaceVarsTests = []variableReplaceTest{
	{
		Description: "Should handle undefined variables",
		Req: &http.Request{
			Method: "GET",
			URL: &url.URL{
				Scheme: "http",
				Host:   "{host}",
				Path:   "/",
			},
			Proto:      "HTTP/1.1",
			ProtoMajor: 1,
			ProtoMinor: 1,
			Header: http.Header{
				"Accept":           {"text/html,application/xhtml+xml,application/xml;q=0.9,*/*;q=0.8"},
				"Accept-Charset":   {"ISO-8859-1,utf-8;q=0.7,*;q=0.7"},
				"Accept-Encoding":  {"gzip,deflate"},
				"Accept-Language":  {"en-us,en;q=0.5"},
				"Keep-Alive":       {"300"},
				"Proxy-Connection": {"keep-alive"},
				"User-Agent":       {"Fake"},
			},
			Body:  nil,
			Close: false,
			Host:  "www.techcrunch.com",
			Form:  map[string][]string{},
		},
		Expected: "GET / HTTP/1.1\r\n" +
			"Host: {host}\r\n" +
			"User-Agent: Fake\r\n" +
			"Accept: text/html,application/xhtml+xml,application/xml;q=0.9,*/*;q=0.8\r\n" +
			"Accept-Charset: ISO-8859-1,utf-8;q=0.7,*;q=0.7\r\n" +
			"Accept-Encoding: gzip,deflate\r\n" +
			"Accept-Language: en-us,en;q=0.5\r\n" +
			"Keep-Alive: 300\r\n" +
			"Proxy-Connection: keep-alive\r\n\r\n",
	},
	{
		Description: "Simple: If variable doesn't exist don't modify request.",
		Req: &http.Request{
			Method: "GET",
			URL: &url.URL{
				Scheme: "http",
				Host:   "{host}",
				Path:   "/",
			},
			Proto:      "HTTP/1.1",
			ProtoMajor: 1,
			ProtoMinor: 1,
			Header: http.Header{
				"Accept":           {"text/html,application/xhtml+xml,application/xml;q=0.9,*/*;q=0.8"},
				"Accept-Charset":   {"ISO-8859-1,utf-8;q=0.7,*;q=0.7"},
				"Accept-Encoding":  {"gzip,deflate"},
				"Accept-Language":  {"en-us,en;q=0.5"},
				"Keep-Alive":       {"300"},
				"Proxy-Connection": {"keep-alive"},
				"User-Agent":       {"Fake"},
			},
			Body:  nil,
			Close: false,
			Host:  "www.techcrunch.com",
			Form:  map[string][]string{},
		},
		Variables: `{
      "host2": "www.techcrunch.com"
    }`,
		Expected: "GET / HTTP/1.1\r\n" +
			"Host: {host}\r\n" +
			"User-Agent: Fake\r\n" +
			"Accept: text/html,application/xhtml+xml,application/xml;q=0.9,*/*;q=0.8\r\n" +
			"Accept-Charset: ISO-8859-1,utf-8;q=0.7,*;q=0.7\r\n" +
			"Accept-Encoding: gzip,deflate\r\n" +
			"Accept-Language: en-us,en;q=0.5\r\n" +
			"Keep-Alive: 300\r\n" +
			"Proxy-Connection: keep-alive\r\n\r\n",
	},
	{
		Description: "Simple: Should replace any variables that exist, and for missing variables don't modify the reqeust.",
		Req: &http.Request{
			Method: "{method}",
			URL: &url.URL{
				Scheme: "http",
				Host:   "{host}",
				Path:   "/",
			},
			Proto:      "HTTP/1.1",
			ProtoMajor: 1,
			ProtoMinor: 1,
			Header: http.Header{
				"Accept":           {"text/html,application/xhtml+xml,application/xml;q=0.9,*/*;q=0.8"},
				"Accept-Charset":   {"ISO-8859-1,utf-8;q=0.7,*;q=0.7"},
				"Accept-Encoding":  {"gzip,deflate"},
				"Accept-Language":  {"en-us,en;q=0.5"},
				"Keep-Alive":       {"300"},
				"Proxy-Connection": {"keep-alive"},
				"User-Agent":       {"Fake"},
			},
			Body:  nil,
			Close: false,
			Host:  "www.techcrunch.com",
			Form:  map[string][]string{},
		},
		Variables: `{
      "host": "www.techcrunch.com"
    }`,
		Expected: "{method} / HTTP/1.1\r\n" +
			"Host: www.techcrunch.com\r\n" +
			"User-Agent: Fake\r\n" +
			"Accept: text/html,application/xhtml+xml,application/xml;q=0.9,*/*;q=0.8\r\n" +
			"Accept-Charset: ISO-8859-1,utf-8;q=0.7,*;q=0.7\r\n" +
			"Accept-Encoding: gzip,deflate\r\n" +
			"Accept-Language: en-us,en;q=0.5\r\n" +
			"Keep-Alive: 300\r\n" +
			"Proxy-Connection: keep-alive\r\n\r\n",
	},
	{
		Description: "Simple: Replace request host with variable",
		Req: &http.Request{
			Method: "GET",
			URL: &url.URL{
				Scheme: "http",
				Host:   "{host}",
				Path:   "/",
			},
			Proto:      "HTTP/1.1",
			ProtoMajor: 1,
			ProtoMinor: 1,
			Header: http.Header{
				"Accept":           {"text/html,application/xhtml+xml,application/xml;q=0.9,*/*;q=0.8"},
				"Accept-Charset":   {"ISO-8859-1,utf-8;q=0.7,*;q=0.7"},
				"Accept-Encoding":  {"gzip,deflate"},
				"Accept-Language":  {"en-us,en;q=0.5"},
				"Keep-Alive":       {"300"},
				"Proxy-Connection": {"keep-alive"},
				"User-Agent":       {"Fake"},
			},
			Body:  nil,
			Close: false,
			Host:  "www.techcrunch.com",
			Form:  map[string][]string{},
		},
		Variables: `{
      "host": "techcrunch.com"
    }`,
		Expected: "GET / HTTP/1.1\r\n" +
			"Host: techcrunch.com\r\n" +
			"User-Agent: Fake\r\n" +
			"Accept: text/html,application/xhtml+xml,application/xml;q=0.9,*/*;q=0.8\r\n" +
			"Accept-Charset: ISO-8859-1,utf-8;q=0.7,*;q=0.7\r\n" +
			"Accept-Encoding: gzip,deflate\r\n" +
			"Accept-Language: en-us,en;q=0.5\r\n" +
			"Keep-Alive: 300\r\n" +
			"Proxy-Connection: keep-alive\r\n\r\n",
	},
	{
		Description: "Simple: Replace request query string with variable",
		Req: &http.Request{
			Method: "GET",
			URL: &url.URL{
				Scheme:   "http",
				Host:     "todos.stoplight.io",
				Path:     "/todos",
				RawQuery: "apikey={apikey}",
			},
			Proto:      "HTTP/1.1",
			ProtoMajor: 1,
			ProtoMinor: 1,
			Header: http.Header{
				"Accept":           {"text/html,application/xhtml+xml,application/xml;q=0.9,*/*;q=0.8"},
				"Accept-Charset":   {"ISO-8859-1,utf-8;q=0.7,*;q=0.7"},
				"Accept-Encoding":  {"gzip,deflate"},
				"Accept-Language":  {"en-us,en;q=0.5"},
				"Keep-Alive":       {"300"},
				"Proxy-Connection": {"keep-alive"},
				"User-Agent":       {"Fake"},
			},
			Body:  nil,
			Close: false,
			Host:  "todos.stoplight.io",
			Form:  map[string][]string{},
		},
		Variables: `{
      "apikey": "123"
    }`,
		Expected: "GET /todos?apikey=123 HTTP/1.1\r\n" +
			"Host: todos.stoplight.io\r\n" +
			"User-Agent: Fake\r\n" +
			"Accept: text/html,application/xhtml+xml,application/xml;q=0.9,*/*;q=0.8\r\n" +
			"Accept-Charset: ISO-8859-1,utf-8;q=0.7,*;q=0.7\r\n" +
			"Accept-Encoding: gzip,deflate\r\n" +
			"Accept-Language: en-us,en;q=0.5\r\n" +
			"Keep-Alive: 300\r\n" +
			"Proxy-Connection: keep-alive\r\n\r\n",
	},
	{
		Description: "Simple: Replace request query string with variable that is a number",
		Req: &http.Request{
			Method: "GET",
			URL: &url.URL{
				Scheme:   "http",
				Host:     "todos.stoplight.io",
				Path:     "/todos",
				RawQuery: "apikey={apikey}",
			},
			Proto:      "HTTP/1.1",
			ProtoMajor: 1,
			ProtoMinor: 1,
			Header: http.Header{
				"Accept":           {"text/html,application/xhtml+xml,application/xml;q=0.9,*/*;q=0.8"},
				"Accept-Charset":   {"ISO-8859-1,utf-8;q=0.7,*;q=0.7"},
				"Accept-Encoding":  {"gzip,deflate"},
				"Accept-Language":  {"en-us,en;q=0.5"},
				"Keep-Alive":       {"300"},
				"Proxy-Connection": {"keep-alive"},
				"User-Agent":       {"Fake"},
			},
			Body:  nil,
			Close: false,
			Host:  "todos.stoplight.io",
			Form:  map[string][]string{},
		},
		Variables: `{
      "apikey": 123
    }`,
		Expected: "GET /todos?apikey=123 HTTP/1.1\r\n" +
			"Host: todos.stoplight.io\r\n" +
			"User-Agent: Fake\r\n" +
			"Accept: text/html,application/xhtml+xml,application/xml;q=0.9,*/*;q=0.8\r\n" +
			"Accept-Charset: ISO-8859-1,utf-8;q=0.7,*;q=0.7\r\n" +
			"Accept-Encoding: gzip,deflate\r\n" +
			"Accept-Language: en-us,en;q=0.5\r\n" +
			"Keep-Alive: 300\r\n" +
			"Proxy-Connection: keep-alive\r\n\r\n",
	},
	{
		Description: "Nested: Replace request host and path with variable",
		Req: &http.Request{
			Method: "GET",
			URL: &url.URL{
				Scheme: "http",
				Host:   "www.{url.host}",
				Path:   "{url.path}",
			},
			ProtoMajor:       1,
			ProtoMinor:       1,
			Header:           http.Header{},
			TransferEncoding: []string{"chunked"},
		},
		Body: []byte("abcdef"),
		Variables: `{
      "url": {
        "host": "google.com",
        "path": "/search"
      }
    }`,
		Expected: "GET /search HTTP/1.1\r\n" +
			"Host: www.google.com\r\n" +
			"User-Agent: Go-http-client/1.1\r\n" +
			"Transfer-Encoding: chunked\r\n\r\n" +
			chunk("abcdef") + chunk(""),
	},
	{
		Description: "Simple: Replace request (plain) body with variable.",
		Req: &http.Request{
			Method: "POST",
			URL: &url.URL{
				Scheme: "http",
				Host:   "www.google.com",
				Path:   "/search",
			},
			ProtoMajor:       1,
			ProtoMinor:       1,
			Header:           http.Header{},
			Close:            true,
			TransferEncoding: []string{"chunked"},
		},
		Body: []byte("{body}"),
		Variables: `{
      "body": "abcdef",
    }`,
		Expected: "POST /search HTTP/1.1\r\n" +
			"Host: www.google.com\r\n" +
			"User-Agent: Go-http-client/1.1\r\n" +
			"Connection: close\r\n" +
			"Transfer-Encoding: chunked\r\n\r\n" +
			chunk("abcdef") + chunk(""),
	},

	// TODO: Add a request with json body, and replace variables on it.
	// TODO: Add a request with url encoded body and replace variables on it.
	// TODO: Add a request with xml body, and replace a variable in it.
	// TODO: Add a request and replace a header value with a variable.
	// TODO: Add a request with this header: "Authorization: Bearer {apikey}" and apikey = "123". Result header should be "Authorization: Bearer 123".
	// TODO: Add a request with this header: "Authorization: Bearer {apikey}" and apikey = 123. Result header should be "Authorization: Bearer 123"
	// TODO: Add a request with this json body: {"foo": "{bool}"}, and use these variables {"bool": true}. The result body should be {"foo": true}.
	// TODO: Add a request with this json body: {"foo": "{str}"}, and use these variables {"str": "bar"}. The result body should be {"foo": "bar"}.
	// TODO: Add a request with this json body: {"foo": "{array}", and use these variables {"array": [1,2,3]}. The result body should be {"foo": [1,2,3]}.
	// TODO: Add a request with this json body: {"foo": "{array[1]}", and use these variables {"array": [1,2,3]}. The result body should be {"foo": 2}.
	// TODO: Add a request with this json body: {"foo": "{array[1]}", and use these variables {"array": ["1","2","3"]}. The result body should be {"foo": "2"}.
	// TODO: Add a request with this header: "Bear: foo,{honey}", and use these variables {"honey": "yum"}. Result header should be "Bear: foo,yum".
	// TODO (Thoughts on if this is too hard?): Add a request with this body: {foo: {"bar": 123, "cat": {$.test}}}, and use these variables {"test": 123}. Body should be {foo: {"bar": 123, "cat": 123}}.

	/*

	  Bonus: If the above wasn't hard enough. Like the instructions above, write your test first, then get them to pass. You do not have to benchmark if you don't want though, but if you do, please follow instructions above.

	*/
	// TODO: Make replace variable function generic, so that any object can be passed into it. Ideally, this is what you would start out with first, then create a replace variable function for handling requests.
	// TODO: Uncomment test case below. Body is {{key}}, and body should equal abcdef.
	// TODO TIP: You will need to run variable replacement twice, phase one will replace key with body {body}, and the second phase will replace {body} with abcdef.
	//{
	//  Description: Bonus: Replace request body with nested variable syntax.",
	//  Req: &http.Request{
	//    Method: "POST",
	//    URL: &url.URL{
	//      Scheme: "http",
	//      Host:   "www.google.com",
	//      Path:   "/search",
	//    },
	//    ProtoMajor:       1,
	//    ProtoMinor:       1,
	//    Header:           http.Header{},
	//    Close:            true,
	//    TransferEncoding: []string{"chunked"},
	//  },
	//  Body: []byte("{{key}"),
	//  Variables: `{
	//    "key": "body",
	//    "body": "abcdef",
	//  }`,
	//  Expected: "POST /search HTTP/1.1\r\n" +
	//    "Host: www.google.com\r\n" +
	//    "User-Agent: Go-http-client/1.1\r\n" +
	//    "Connection: close\r\n" +
	//    "Transfer-Encoding: chunked\r\n\r\n" +
	//    chunk("abcdef") + chunk(""),
	//},
	// TODO: Add a request with form-data body(no file) and replace variables on it.
	// TODO: Add a request with form-data body(Should have a file, file doesn't need to a be a variable but you need to make sure the file is still there after replacing variables) and replace variables on it.
	// TODO: Handle a request url has extra slashes and a varible. Replace the variable and clean up the request url.
	// TODO: Handle a request with a stringified JSON body "{\"foo\": \"{bar}\"}", and use these variables {"bar": "bear"}. Expected body is {"foo": "bear"}.
	// TODO: Handle a request with a stringified JSON body "{\"foo\": \"{bar}\"}", and use these variables {"bar": 123}. Expected body is {"foo": 123}.
	// TODO: Handle a request with a (invalid json) stringified JSON body "{\"foo\": \"{bar}\"}]", and use these variables {"bar": "bear"}. Expected body is "{\"foo\": \"bear\"}]".
	// TODO: Handle a request with content type url-encoded, but the input body is a json object. Input: headers -> "Content-Type: application/x-www-form-urlencoded" -> {"foo": "{bear}"}, variables -> {"bear": "bar"}. Expected body is foo=bar.
	// TODO: Return array of variable keys that couldn't be found.

	// TODO Surprise: Add something cool that makes sense for this. For example, spell check for variable keys and path. Help me explain this better. You could set up a proxy, and replace variables are requests come into it. More than happy to help you come up with an idea or help you flesh it out.
}

func TestRequestReplaceVariables(t *testing.T) {
	for i := range reqReplaceVarsTests {
		tt := &reqReplaceVarsTests[i]

		setBody := func() {
			if tt.Body == nil {
				return
			}
			switch b := tt.Body.(type) {
			case []byte:
				tt.Req.Body = ioutil.NopCloser(bytes.NewReader(b))
			case func() io.ReadCloser:
				tt.Req.Body = b()
			}
		}
		setBody()

		if tt.Req.Header == nil {
			tt.Req.Header = make(http.Header)
		}

		// TODO: Call replace variable function here.
		bearcub.Replace(tt.Req, tt.Variables)
		dump, err := httputil.DumpRequestOut(tt.Req, true)
		if err != nil {
			t.Errorf("Test %s, error building reqeust %s.", tt.Description, err.Error())
		}

		if tt.Expected != "" {
			sraw := string(dump)
			if sraw != tt.Expected {
				t.Errorf("Test %s, expecting:\n%s\nGot:\n%s\n", tt.Description, tt.Expected, sraw)
				continue
			}
		}
	}
}

func chunk(s string) string {
	return fmt.Sprintf("%x\r\n%s\r\n", len(s), s)
}
