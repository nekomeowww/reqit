package reqit

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/url"
	"reflect"
	"strings"

	"github.com/samber/lo"
)

var (
	ErrRequestCalledWithNilRequest = errors.New("directly called Result() without Request being created")
)

type Request struct {
	request *http.Request

	client   *Client
	response *http.Response

	header   http.Header
	method   string
	uri      string
	param    url.Values
	body     io.Reader
	result   *bytes.Buffer
	rawError error
}

func newRequest(c *Client, method, uri string, param url.Values) *Request {
	req := &Request{client: c, header: http.Header{}}
	if req.client.rawError != nil {
		req.rawError = req.client.rawError
		return req
	}

	if c.Options.baseURL != nil {
		req.uri = c.Options.baseURL.String() + uri
	} else {
		if !strings.HasPrefix(uri, "http://") && !strings.HasPrefix(uri, "https://") {
			req.uri = "http://" + uri
		}
	}

	req.param = param

	return req
}

func (r *Request) do() error {
	var newRequest *http.Request
	var err error
	if r.method == http.MethodGet {
		newRequest, err = http.NewRequest(r.method, r.uri, strings.NewReader(r.param.Encode()))
		if err != nil {
			r.rawError = err
			return err
		}
	} else {
		newRequest, err = http.NewRequest(r.method, r.uri, r.body)
		if err != nil {
			r.rawError = err
			return err
		}
	}

	newRequest.Header = http.Header{}
	mergeHeadersFrom(newRequest.Header, r.client.Options.Headers)
	mergeHeadersFrom(newRequest.Header, r.header)

	r.request = newRequest

	resp, err := r.client.Do(r.request)
	if err != nil {
		r.rawError = err
		return err
	}

	r.response = resp
	defer r.response.Body.Close()

	r.result = new(bytes.Buffer)
	_, err = io.Copy(r.result, resp.Body)
	if err != nil {
		r.rawError = err
		return err
	}

	return nil
}

func (r *Request) unmarshalToJSON(dest any) error {
	err := json.Unmarshal(r.result.Bytes(), dest)
	if err != nil {
		r.rawError = err
		return err
	}

	return nil
}

func (r *Request) parseDestKind(dest any) reflect.Kind {
	refVal := reflect.ValueOf(dest)
	if refVal.Kind() == reflect.Ptr {
		return refVal.Elem().Kind()
	}

	return refVal.Kind()
}

func (r *Request) Send() *Request {
	if r.rawError != nil {
		return r
	}

	err := r.do()
	if err != nil {
		r.rawError = err
		return r
	}

	return r
}

func (r *Request) Result(dest any) (int, error) {
	if r.rawError != nil {
		return 0, r.rawError
	}

	switch r.parseDestKind(dest) {
	case reflect.String:
		switch v := dest.(type) {
		case *string:
			*v = r.result.String()
		}
	case reflect.Map, reflect.Array, reflect.Struct:
		err := r.unmarshalToJSON(dest)
		if err != nil {
			return 0, err
		}
	}

	return r.response.StatusCode, nil
}

func (r *Request) Error() error {
	return r.rawError
}

func unpointer(v any) any {
	refVal := reflect.ValueOf(v)
	if refVal.Kind() == reflect.Ptr {
		return refVal.Elem().Interface()
	}

	return v
}

func (r *Request) WithBody(body any) *Request {
	switch v := body.(type) {
	case io.Reader:
		r.body = v
	case string:
		r.body = strings.NewReader(v)
		r.header.Set("Content-Type", "text/plain")
	case *string:
		if v != nil {
			r.body = strings.NewReader(*v)
		}

		r.body = strings.NewReader("")
		r.header.Set("Content-Type", "text/plain")
	default:
		body := unpointer(body)
		refVal := reflect.ValueOf(body)
		switch refVal.Kind() {
		case reflect.Struct, reflect.Array, reflect.Map:
			bodyData, err := json.Marshal(body)
			if err != nil {
				r.rawError = err
				return r
			}

			r.body = bytes.NewReader(bodyData)
			r.header.Set("Content-Type", "application/json")
		}
	}

	return r
}

func (r *Request) Headers(headers http.Header) *Request {
	mergeHeadersFrom(r.request.Header, headers)
	return r
}

func (r *Request) HeadersFromMap(headers map[string]string) *Request {
	for k, v := range headers {
		r.request.Header.Set(k, v)
	}

	return r
}

func mergeHeadersFrom(headers http.Header, fromHeaders http.Header) {
	if headers == nil || fromHeaders == nil {
		// fallback
		return
	}

	for key, values := range fromHeaders {
		// append then unique, drop the duplicates
		headers[key] = lo.Uniq(append(headers[key], values...))
	}
}
