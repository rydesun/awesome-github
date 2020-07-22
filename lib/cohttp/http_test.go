package cohttp

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"
)

var testServer *httptest.Server

func init() {
	handler := func(rw http.ResponseWriter, req *http.Request) {
		if req.URL.Path == "/text" {
			rw.Write([]byte("simple text"))
		}
		if req.URL.Path == "/text/wrong_path" {
			rw.WriteHeader(400)
			rw.Write([]byte("permisson denied"))
		}
		if req.URL.Path == "/json" {
			rw.Write([]byte(`{"foo": "bar"}`))
		}
		if req.URL.Path == "/json/wrong_path" {
			rw.WriteHeader(500)
		}
		if req.URL.Path == "/json/invalid" {
			rw.Write([]byte(`{`))
		}
	}
	testServer = httptest.NewServer(http.HandlerFunc(handler))
}

func TestClient_Text(t *testing.T) {
	require := require.New(t)
	testCases := []struct {
		in     string
		out    string
		hasErr bool
		netErr bool
	}{
		{
			in:  "/text",
			out: "simple text",
		},
		{
			in:     "/text/wrong_path",
			out:    "permisson denied",
			hasErr: true,
		},
		{
			hasErr: true,
			netErr: true,
		},
	}
	client := NewClient(*testServer.Client(), 16, 20, nil)
	for _, tc := range testCases {
		t.Run("server"+tc.in, func(t *testing.T) {
			var url string
			if tc.netErr {
				url = "invalid url"
			} else {
				url = testServer.URL + tc.in
			}
			req, err := http.NewRequest(
				http.MethodGet, url, nil)
			require.Nil(err)
			result, err := client.Text(req)
			if !tc.hasErr {
				require.Nil(err)
			} else {
				require.NotNil(err)
			}
			require.Equal(tc.out, result)
		})
	}
}

func TestClient_Json(t *testing.T) {
	require := require.New(t)
	type Result struct {
		Foo string
	}
	testCases := []struct {
		in     string
		out    Result
		hasErr bool
		netErr bool
	}{
		{
			in: "/json",
			out: Result{
				Foo: "bar",
			},
		},
		{
			in:     "/json/wrong_path",
			hasErr: true,
		},
		{
			in:     "/json/invalid",
			hasErr: true,
		},
		{
			hasErr: true,
			netErr: true,
		},
	}

	client := NewClient(*testServer.Client(), 16, 20, nil)
	for _, tc := range testCases {
		t.Run("server"+tc.in, func(t *testing.T) {
			var url string
			if tc.netErr {
				url = "invalid url"
			} else {
				url = testServer.URL + tc.in
			}
			req, err := http.NewRequest(
				http.MethodGet, url, nil)
			require.Nil(err)
			var result Result
			err = client.Json(req, &result)
			if !tc.hasErr {
				require.Nil(err)
			} else {
				require.NotNil(err)
			}
			require.Equal(tc.out, result)
		})
	}
}
