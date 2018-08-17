package rest

import (
	"io"
	"io/ioutil"
	"net/http"
)

type RequestBuilder interface {
	BODY(body io.Reader) RequestBuilder
	WithHeader(key, value string) RequestBuilder
	URL(url string) RequestBuilder
	GET() ([]byte, error)
	POST() ([]byte, error)
}

type HTTP struct {
	Client  *http.Client
	Request *http.Request
	Body    io.Reader
	Url     string
	Headers map[string]string
}

func (h *HTTP) URL(url string) RequestBuilder {
	h.Url = url
	return h
}

func (h *HTTP) BODY(body io.Reader) RequestBuilder {
	h.Body = body
	return h
}

func NewHTTPClient(client *http.Client) RequestBuilder {
	h := new(HTTP)
	h.Client = client
	h.Headers = make(map[string]string)
	h.Request, _ = http.NewRequest("", "", nil)
	return h
}

func (h *HTTP) WithHeader(key, value string) RequestBuilder {
	h.Headers[key] = value
	return h
}

func (h *HTTP) POST() ([]byte, error) {
	h.Request, _ = http.NewRequest(http.MethodPost, h.Url, h.Body)

	if len(h.Headers) != 0 {
		for k, v := range h.Headers {
			h.Request.Header.Add(k, v)
			delete(h.Headers, k)
		}
	}

	response, resErr := h.Client.Do(h.Request)
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
	h.Request, _ = http.NewRequest(http.MethodGet, h.Url, nil)

	if len(h.Headers) != 0 {
		for k, v := range h.Headers {
			h.Request.Header.Add(k, v)
			delete(h.Headers, k)
		}
	}
	resp, reqErr := h.Client.Do(h.Request)

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
