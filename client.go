package reqit

import (
	"net/http"
	"net/url"
	"strings"
)

type Client struct {
	*http.Client

	Options *ClientOptions

	rawError error
}

func NewClient(options ...ClientOptions) *Client {
	client := &Client{
		Client:  &http.Client{},
		Options: new(ClientOptions),
	}
	if len(options) > 0 {
		configureClientByClientOptions(client, options[0])
	}

	return client
}

type ClientOptions struct {
	BaseURL string
	Headers http.Header

	baseURL *url.URL
}

func configureClientByClientOptions(c *Client, options ClientOptions) {
	if options.Headers != nil {
		c.Options.Headers = options.Headers
	}
	if options.BaseURL != "" {
		if !strings.HasPrefix(options.BaseURL, "http://") && !strings.HasPrefix(options.BaseURL, "https://") {
			options.BaseURL = "http://" + options.BaseURL
		}

		url, err := url.Parse(options.BaseURL)
		if err != nil {
			c.rawError = err
			return
		} else {
			c.Options.baseURL = url
			c.Options.BaseURL = c.Options.baseURL.String()
		}
	}
}

func (c *Client) Get(uri string, query map[string]string) *Request {
	return newRequest(c, http.MethodGet, uri, mapToQuery(query))
}

func (c *Client) Post(uri string, query map[string]string) *Request {
	return newRequest(c, http.MethodPost, uri+mapToQuery(query).Encode(), nil)
}
