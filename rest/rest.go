package rest

import (
	"io"
	"io/ioutil"
	"net/http"
)

type RequestBuilder interface {
	Body(body io.Reader) RequestBuilder
	WithHeader(key, value string) RequestBuilder
	URL(url string) RequestBuilder
	GET() ([]byte, error)
	POST() ([]byte, error)
}

type HTTP struct {
	client  *http.Client
	request *http.Request
	body    io.Reader
	url     string
	headers map[string]string
}

func (h *HTTP) URL(url string) RequestBuilder {
	h.url = url
	return h
}

func (h *HTTP) Body(body io.Reader) RequestBuilder {
	h.body = body
	return h
}

func NewHTTPClient(client *http.Client) RequestBuilder {
	h := new(HTTP)
	h.client = client
	h.headers = make(map[string]string)
	h.request, _ = http.NewRequest("", "", nil)
	return h
}

func (h *HTTP) WithHeader(key, value string) RequestBuilder {
	h.headers[key] = value
	return h
}

func (h *HTTP) POST() ([]byte, error) {
	h.request, _ = http.NewRequest(http.MethodPost, h.url, h.body)

	if len(h.headers) != 0 {
		for k, v := range h.headers {
			h.request.Header.Add(k, v)
			delete(h.headers, k)
		}
	}

	response, resErr := h.client.Do(h.request)
	if resErr != nil {
		return nil, resErr
	}

	contents, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	return contents, nil
}

func (h *HTTP) GET() ([]byte, error) {
	h.request, _ = http.NewRequest(http.MethodGet, h.url, nil)

	if len(h.headers) != 0 {
		for k, v := range h.headers {
			h.request.Header.Add(k, v)
			delete(h.headers, k)
		}
	}
	resp, reqErr := h.client.Do(h.request)

	if reqErr != nil {
		return nil, reqErr
	}

	defer resp.Body.Close()

	contents, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return contents, nil
}
