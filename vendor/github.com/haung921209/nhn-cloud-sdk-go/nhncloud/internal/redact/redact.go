// Package redact masks credentials before they reach SDK debug output or
// response capture files. Every debug path in the SDK must go through it.
package redact

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"sort"
	"strings"
)

// Mask replaces every redacted value.
const Mask = "***"

var sensitiveHeaderNames = map[string]bool{
	"x-auth-token":               true,
	"x-tc-authentication-secret": true,
	"x-tc-authentication-id":     true,
	"x-tc-app-key":               true,
	"authorization":              true,
	"x-nhn-authorization":        true,
	"proxy-authorization":        true,
	"cookie":                     true,
	"set-cookie":                 true,
}

var sensitiveHeaderParts = []string{"secret", "token", "password", "key", "auth", "cookie", "credential"}

var sensitiveKeyParts = []string{
	"password", "passwd", "secret", "token", "privatekey", "private_key",
	"accesskey", "access_key", "secretkey", "credential", "kubeconfig", "payload",
	"appkey", "app_key", "apikey", "api_key",
	"symmetrickey", "plaintext", "adminpass", "authorization", "cookie",
}

// IsSensitiveHeader reports whether a header value must not be printed.
func IsSensitiveHeader(name string) bool {
	n := strings.ToLower(name)
	return sensitiveHeaderNames[n] || containsAny(n, sensitiveHeaderParts)
}

// IsSensitiveKey reports whether a JSON object key (or query parameter)
// holds a value that must not be printed.
func IsSensitiveKey(key string) bool {
	k := strings.ToLower(key)
	// "key" alone is the key material in key manager responses; "auth" alone
	// is matched exactly so fields like authorityType stay readable.
	return k == "key" || k == "auth" || containsAny(k, sensitiveKeyParts)
}

func containsAny(s string, parts []string) bool {
	for _, p := range parts {
		if strings.Contains(s, p) {
			return true
		}
	}
	return false
}

// Headers renders h as sorted "  Name: value" lines with sensitive values masked.
func Headers(h http.Header) string {
	names := make([]string, 0, len(h))
	for name := range h {
		names = append(names, name)
	}
	sort.Strings(names)
	var b strings.Builder
	for _, name := range names {
		for _, v := range h[name] {
			if IsSensitiveHeader(name) {
				v = Mask
			}
			fmt.Fprintf(&b, "  %s: %s\n", name, v)
		}
	}
	return b.String()
}

// URL renders u without userinfo and with sensitive query values masked.
func URL(u *url.URL) string {
	if u == nil {
		return ""
	}
	c := *u
	c.User = nil
	c.Path, c.RawPath = Path(c.Path), ""
	if c.RawQuery != "" {
		q := c.Query()
		for k := range q {
			if IsSensitiveKey(k) {
				q[k] = []string{Mask}
			}
		}
		c.RawQuery = q.Encode()
	}
	return c.String()
}

// URLString is URL for a raw URL string. Unparseable input is not printed.
func URLString(raw string) string {
	u, err := url.Parse(raw)
	if err != nil {
		return "<unparseable URL omitted>"
	}
	return URL(u)
}

// JSON returns body as compact JSON with sensitive values masked. ok is false
// when body is not a JSON document; callers must then not print the body.
func JSON(body []byte) (out []byte, ok bool) {
	dec := json.NewDecoder(bytes.NewReader(body))
	dec.UseNumber()
	var v interface{}
	if err := dec.Decode(&v); err != nil || dec.More() {
		return nil, false
	}
	var buf bytes.Buffer
	enc := json.NewEncoder(&buf)
	enc.SetEscapeHTML(false)
	if err := enc.Encode(maskValue(v)); err != nil {
		return nil, false
	}
	return bytes.TrimRight(buf.Bytes(), "\n"), true
}

func maskValue(v interface{}) interface{} {
	switch t := v.(type) {
	case map[string]interface{}:
		for k, child := range t {
			if IsSensitiveKey(k) {
				t[k] = Mask
			} else {
				t[k] = maskValue(child)
			}
		}
	case []interface{}:
		for i, child := range t {
			t[i] = maskValue(child)
		}
	}
	return v
}

// Body renders a raw body for debug output: masked, indented JSON, or only
// the length when the body is not JSON.
func Body(body []byte) string {
	if len(body) == 0 {
		return ""
	}
	out, ok := JSON(body)
	if !ok {
		return fmt.Sprintf("<%d bytes, non-JSON body omitted>", len(body))
	}
	var buf bytes.Buffer
	if json.Indent(&buf, out, "", "  ") != nil {
		return string(out)
	}
	return buf.String()
}

// Value renders a request body value (struct, map, string, []byte,
// url.Values, io.Reader) for debug output using the same rules as Body.
func Value(v interface{}) string {
	switch t := v.(type) {
	case nil:
		return ""
	case []byte:
		return Body(t)
	case string:
		return Body([]byte(t))
	case io.Reader:
		return "<streamed body omitted>"
	}
	data, err := json.Marshal(v)
	if err != nil {
		return "<unencodable body omitted>"
	}
	return Body(data)
}

// Debugf writes one debug line to stderr. Debug output never goes to stdout,
// which belongs to the calling program.
func Debugf(format string, args ...interface{}) {
	fmt.Fprintf(os.Stderr, format, args...)
}

// Path masks the segment after "appkey" or "appkeys"; some APIs carry the
// appkey in the path (/appkey/<appkey>/...).
func Path(p string) string {
	segs := strings.Split(p, "/")
	for i := 0; i+1 < len(segs); i++ {
		if s := strings.ToLower(segs[i]); (s == "appkey" || s == "appkeys") && segs[i+1] != "" {
			segs[i+1] = Mask
		}
	}
	return strings.Join(segs, "/")
}
