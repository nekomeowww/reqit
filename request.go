package reqit

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"strings"

	"github.com/samber/lo"
)

var (
	ErrRequestCalledWithNilRequest = errors.New("directly called Result() without Request being created")
)

type Request struct {
	*http.Request

	Client   *Client
	Response *http.Response

	result   *bytes.Buffer
	rawError error
}

func newRequest(c *Client, method, uri string, param io.Reader) *Request {
	req := &Request{Client: c}
	if req.Client.rawError != nil {
		req.rawError = req.Client.rawError
		return req
	}

	if c.Options.baseURL != nil {
		uri = c.Options.baseURL.String() + uri
	} else {
		if !strings.HasPrefix(uri, "http://") && !strings.HasPrefix(uri, "https://") {
			uri = "http://" + uri
		}
	}

	r, err := http.NewRequest(method, uri, param)
	if err != nil {
		req.rawError = err
		return req
	}

	r.Header = http.Header{}
	mergeHeadersFrom(r.Header, c.Options.Headers)

	req.Request = r
	return req
}

func (r *Request) do() error {
	resp, err := r.Client.Do(r.Request)
	if err != nil {
		r.rawError = err
		return err
	}

	r.Response = resp
	defer r.Body.Close()

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

func (r *Request) Result(dest any) *Request {
	if r.Request == nil {
		if r.rawError != nil {
			return r
		}

		r.rawError = ErrRequestCalledWithNilRequest
		return r
	}

	err := r.do()
	if err != nil {
		return r
	}

	err = r.unmarshalToJSON(dest)
	if err != nil {
		return r
	}

	return r
}

func (r *Request) Error() error {
	return r.rawError
}

func (r *Request) Headers(headers http.Header) *Request {
	mergeHeadersFrom(r.Request.Header, headers)
	return r
}

func (r *Request) HeadersFromMap(headers map[string]string) *Request {
	for k, v := range headers {
		r.Header.Set(k, v)
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
