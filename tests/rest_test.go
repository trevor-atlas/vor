package tests

import (
	"github.com/trevor-atlas/vor/rest"
	"io"
)

type mockHTTP struct{}

func (h *mockHTTP) Body(body io.Reader) rest.RequestBuilder {
	return h
}

func (h *mockHTTP) WithHeader(key, value string) rest.RequestBuilder {
	return h
}

func (h *mockHTTP) GET() ([]byte, error) {
	return []byte{0, 0, 0, 0, 0}, nil
}

func (h *mockHTTP) POST() ([]byte, error) {
	return []byte{0, 0, 0, 0, 0}, nil
}

func (h *mockHTTP) URL(url string) rest.RequestBuilder {
	return h
}
