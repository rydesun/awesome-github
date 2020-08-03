package fakeg

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
)

func HtmlServer(testdataHoler DataHolder) (testServer *httptest.Server, err error) {
	awesomeHtmlReadme, err := testdataHoler.GetHtmlAwesomeReadme()
	if err != nil {
		return
	}
	handler := func(rw http.ResponseWriter, req *http.Request) {
		if req.URL.Path == "/tester/awesome-test/blob/master/README.md" {
			rw.Write([]byte(awesomeHtmlReadme))
		}
		if req.URL.Path == "/invalid/invalid/blob/master/README.md" {
			rw.WriteHeader(400)
		}
	}
	testServer = httptest.NewServer(http.HandlerFunc(handler))
	return
}

func ApiServer(testdataHolder DataHolder) (testServer *httptest.Server, err error) {
	jsonUser, err := testdataHolder.GetUser()
	if err != nil {
		return
	}
	jsonRepoXpath, err := testdataHolder.GetJsonRepo("antchfx", "xpath")
	if err != nil {
		return
	}
	jsonRepoGlmatrix, err := testdataHolder.GetJsonRepo("technohippy", "go-glmatrix")
	if err != nil {
		return
	}
	jsonRepoGlfw, err := testdataHolder.GetJsonRepo("goxjs", "glfw")
	if err != nil {
		return
	}
	handler := func(rw http.ResponseWriter, req *http.Request) {
		if !strings.HasPrefix(req.URL.Path, "/graphql") {
			rw.WriteHeader(400)
			return
		}
		// There is no need to implement a GraphQL server
		raw, err := ioutil.ReadAll(req.Body)
		defer req.Body.Close()
		if err != nil {
			rw.WriteHeader(400)
			return
		}
		if strings.Contains(string(raw), "login") {
			bear := req.Header.Get("Authorization")
			if strings.HasSuffix(bear, "123456") {
				rw.Write([]byte(jsonUser))
			} else {
				rw.WriteHeader(401)
			}
			return
		}
		if strings.Contains(string(raw), "antchfx") {
			rw.Write([]byte(jsonRepoXpath))
			return
		}
		if strings.Contains(string(raw), "technohippy") {
			rw.Write([]byte(jsonRepoGlmatrix))
			return
		}
		if strings.Contains(string(raw), "goxjs") {
			rw.Write([]byte(jsonRepoGlfw))
			return
		}
		rw.WriteHeader(400)
	}
	testServer = httptest.NewServer(http.HandlerFunc(handler))
	return
}
