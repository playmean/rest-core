package restcore

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
)

type HttpRequester struct {
	root string

	client  *http.Client
	headers map[string]string
}

type HttpRequesterOptions struct {
	Query map[string]string
}

func NewHttpRequester(root string) *HttpRequester {
	req := &HttpRequester{
		root: root,

		client:  &http.Client{},
		headers: make(map[string]string),
	}

	return req
}

func (p *HttpRequester) SetHeaders(headers map[string]string) {
	for k, v := range headers {
		p.headers[k] = v
	}
}

func (p *HttpRequester) Get(route string, reqOptions ...*HttpRequesterOptions) (*http.Response, []byte, error) {
	return p.makeRequest("GET", route, nil, reqOptions)
}

func (p *HttpRequester) Post(route string, body []byte, reqOptions ...*HttpRequesterOptions) (*http.Response, []byte, error) {
	return p.makeRequest("POST", route, body, reqOptions)
}

func (p *HttpRequester) Put(route string, body []byte, reqOptions ...*HttpRequesterOptions) (*http.Response, []byte, error) {
	return p.makeRequest("PUT", route, body, reqOptions)
}

func (p *HttpRequester) Delete(route string, body []byte, reqOptions ...*HttpRequesterOptions) (*http.Response, []byte, error) {
	return p.makeRequest("DELETE", route, body, reqOptions)
}

func (p *HttpRequester) makeRequest(method string, route string, body []byte, reqOptions []*HttpRequesterOptions) (*http.Response, []byte, error) {
	options := p.getReqOptions(reqOptions)
	path := p.getReqPath(route, options)

	var r *http.Request
	var err error

	if body != nil {
		r, err = http.NewRequest(method, path, bytes.NewReader(body))
	} else {
		r, err = http.NewRequest(method, path, nil)
	}
	if err != nil {
		return nil, nil, NewApiError(&ApiErrorOptions{
			Code:     "REQUEST",
			Message:  fmt.Sprintf("error while creating %s request", method),
			Original: err,
		})
	}

	return p.doRequest(r)
}

func (p *HttpRequester) getReqOptions(reqOptions []*HttpRequesterOptions) *HttpRequesterOptions {
	if len(reqOptions) > 0 {
		return reqOptions[0]
	}

	return new(HttpRequesterOptions)
}

func (p *HttpRequester) getReqPath(route string, options *HttpRequesterOptions) string {
	path := p.root + route

	if options.Query != nil {
		params := url.Values{}

		for key, value := range options.Query {
			params.Set(key, value)
		}

		path += "?" + params.Encode()
	}

	return path
}

func (p *HttpRequester) doRequest(r *http.Request) (*http.Response, []byte, error) {
	for k, v := range p.headers {
		r.Header.Add(k, v)
	}

	resp, err := p.client.Do(r)
	if err != nil {
		return nil, nil, NewApiError(&ApiErrorOptions{
			Code:     "REQUEST",
			Message:  "error while sending request",
			Original: err,
		})
	}

	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return resp, nil, NewApiError(&ApiErrorOptions{
			Code:     "REQUEST",
			Message:  "error while reading response body",
			Original: err,
		})
	}

	if resp.StatusCode != 200 {
		return resp, respBody, NewApiError(&ApiErrorOptions{
			Code:    "REQUEST",
			Message: string(respBody),
		})
	}

	return resp, respBody, nil
}
