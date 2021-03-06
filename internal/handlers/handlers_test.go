package handlers

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
)

type postData struct {
	key   string
	value string
}

var theTests = []struct {
	name       string
	url        string
	method     string
	params     []postData
	statusCode int
}{
	{"home", "/", "GET", []postData{}, http.StatusOK},
	{"about", "/", "GET", []postData{}, http.StatusOK},
	{"general quarters", "/generals-quarters", "GET", []postData{}, http.StatusOK},
	{"majors suite", "/majors-suite", "GET", []postData{}, http.StatusOK},
	{"search availability", "/search-availability", "GET", []postData{}, http.StatusOK},
	{"contact", "/contact", "GET", []postData{}, http.StatusOK},
	{"make reservation", "/make-reservation", "GET", []postData{}, http.StatusOK},
	{"search availability", "/search-availability", "POST", []postData{
		{key: "start", value: "2020-01-01"},
		{key: "end", value: "2020-01-11"},
	}, http.StatusOK},
	{"search availability JSON", "/search-availability-json", "POST", []postData{
		{key: "start", value: "2020-01-01"},
		{key: "end", value: "2020-01-11"},
	}, http.StatusOK},
	{"make reservation post", "/make-reservation", "POST", []postData{
		{key: "first_name", value: "John"},
		{key: "last_name", value: "Doe"},
		{key: "email", value: "some.email@google.com"},
		{key: "phone", value: "777 555 999"},
	}, http.StatusOK},
}

func TestHandlers(t *testing.T) {
	routes := GetRoutes()

	testServer := httptest.NewServer(routes)
	defer testServer.Close()

	for _, test := range theTests {
		if test.method == "GET" {
			res, err := testServer.Client().Get(testServer.URL + test.url)
			if err != nil {
				t.Log(err)
				t.Fatal(err)
			}

			if res.StatusCode != test.statusCode {
				t.Errorf("for %s expected %d but got %d", test.name, test.statusCode, res.StatusCode)
			}
		} else {
			values := url.Values{}

			for _, i := range test.params {
				values.Add(i.key, i.value)
			}

			res, err := testServer.Client().PostForm(testServer.URL+test.url, values)
			if err != nil {
				t.Log(err)
				t.Fatal(err)
			}

			if res.StatusCode != test.statusCode {
				t.Errorf("for %s expected %d but got %d", test.name, test.statusCode, res.StatusCode)
			}
		}
	}

}
