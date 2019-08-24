package proxy

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/blend/go-sdk/assert"
	"github.com/blend/go-sdk/logger"
)

func TestUpstreamWithoutHopHeaders(t *testing.T) {
	assert := assert.New(t)

	u := NewUpstream(MustParseURL("http://localhost:5000"))
	assert.NotNil(u.ReverseProxy)
}

func TestUpstreamServeHTTP(t *testing.T) {
	assert := assert.New(t)

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, "OK!\n")
	}))
	defer srv.Close()

	u := NewUpstream(MustParseURL(srv.URL))
	assert.NotNil(u.ReverseProxy)
	u.Log = logger.None()

	proxy := New()
	proxy.WithUpstream(u)

	server := &http.Server{}
	server.Handler = proxy
	server.Addr = fmt.Sprintf("localhost:%s", "5000")
	go server.ListenAndServe()
	defer server.Close()

	req, err := http.NewRequest("GET", "http://localhost:5000", nil)
	assert.Nil(err)

	res, err := http.DefaultClient.Do(req)
	assert.Nil(err)

	assert.Empty(res.Header.Get("X-Forwarded-For"))
	assert.Empty(res.Header.Get("X-Forwarded-Port"))
}