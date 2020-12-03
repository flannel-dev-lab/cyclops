package router

import (
	"fmt"
	"github.com/flannel-dev-lab/cyclops/v2"
	"io/ioutil"
	"math/rand"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"strings"
	"testing"
	"time"
)

func TestRouterParam(t *testing.T) {
	r := New(true, nil, nil)
	r.Get("/users/:id", func(w http.ResponseWriter, r *http.Request) {})
	req, _ := http.NewRequest("GET", "/users/1", nil)
	w := httptest.NewRecorder()

	h, err := r.find(req)
	if err != nil {
		t.Errorf("%s: unable to get handler: %s", t.Name(), err.Error())
	}

	h(w, req)

	if cyclops.Param(req, "id") != "1" {
		t.Errorf(fmt.Sprintf("%s: params do not match", t.Name()))
	}
}

func TestRouterInitializationWithCustomNotFoundHandler(t *testing.T) {
	notFoundHandler := func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	}

	r := New(true, notFoundHandler, nil)
	r.Get("/use", func(w http.ResponseWriter, r *http.Request) {})

	req, _ := http.NewRequest("GET", "/users/1", nil)
	w := httptest.NewRecorder()

	h, err := r.find(req)
	if err != nil {
		t.Errorf("%s: unable to get handler: %s", t.Name(), err.Error())
	}

	h(w, req)

	if w.Code != http.StatusNotFound {
		t.Errorf("%s: unable to get not found status code", t.Name())
	}
}

func TestRouterInitializationWithCustomMethodNotAllowedHandler(t *testing.T) {
	notAllowedHandler := func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusMethodNotAllowed)
	}

	r := New(true, nil, notAllowedHandler)
	r.Post("/use", func(w http.ResponseWriter, r *http.Request) {})

	req, _ := http.NewRequest("GET", "/use", nil)
	w := httptest.NewRecorder()

	h, err := r.find(req)
	if err != nil {
		t.Errorf("%s: unable to get handler: %s", t.Name(), err.Error())
	}

	h(w, req)

	if w.Code != http.StatusMethodNotAllowed {
		t.Errorf("%s: unable to get not found status code", t.Name())
	}
}

func TestRouterInitializationWithCyclopsNotFoundHandler(t *testing.T) {
	r := New(true, nil, nil)
	r.Get("/use", func(w http.ResponseWriter, r *http.Request) {})

	req, _ := http.NewRequest("GET", "/users/1", nil)
	w := httptest.NewRecorder()

	h, err := r.find(req)
	if err != nil {
		t.Errorf("%s: unable to get handler: %s", t.Name(), err.Error())
	}

	h(w, req)

	if w.Code != http.StatusNotFound {
		t.Errorf("%s: unable to get not found status code", t.Name())
	}
}

func TestRouterInitializationWithCyclopsMethodNotAllowedHandler(t *testing.T) {
	r := New(true, nil, nil)
	r.Post("/use", func(w http.ResponseWriter, r *http.Request) {})

	req, _ := http.NewRequest("GET", "/use", nil)
	w := httptest.NewRecorder()

	h, err := r.find(req)
	if err != nil {
		t.Errorf("%s: unable to get handler: %s", t.Name(), err.Error())
	}

	h(w, req)

	if w.Code != http.StatusMethodNotAllowed {
		t.Errorf("%s: unable to get not found status code", t.Name())
	}
}

func TestRouter_IncorrectInputs(t *testing.T) {
	cases := []struct {
		Method  string
		Path    string
		Handler http.HandlerFunc
	}{
		{"", "/", nil},
		{"GET", "users/", nil},
		{"GET", "/path", nil},
	}

	for _, testCase := range cases {
		codeThatPanics(t, testCase.Method, testCase.Path, testCase.Handler)
	}
}

func codeThatPanics(t *testing.T, method, path string, handler http.HandlerFunc) {
	defer func() { recover() }()

	r := New(true, nil, nil)
	r.add(method, path, handler)

	t.Errorf("should have panicked")
}

func TestRouter_ServeHTTP(t *testing.T) {
	r := New(true, nil, nil)
	r.Post("/use", func(w http.ResponseWriter, r *http.Request) {})

	req, _ := http.NewRequest("POST", "/use", nil)
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Errorf("%s: expected 200 got %d", t.Name(), w.Code)
	}
}

func TestRouter_FileServer(t *testing.T) {
	err := os.Mkdir("static", 0777)
	if err != nil {
		t.Fatalf("%s: %s", t.Name(), err.Error())
	}

	err = ioutil.WriteFile("static/index.html", nil, 0777)
	if err != nil {
		t.Fatalf("%s: %s", t.Name(), err.Error())
	}

	cases := []struct {
		DirectoryPath string
		StatusCode    int
	}{
		{"static", http.StatusOK},
		{"bad_path", http.StatusNotFound},
	}

	for _, testCase := range cases {
		r := New(false, nil, nil)
		r.RegisterStatic(testCase.DirectoryPath, "/static/")

		req, _ := http.NewRequest("GET", "/static/", nil)
		w := httptest.NewRecorder()

		r.ServeHTTP(w, req)

		if w.Code != testCase.StatusCode {
			t.Errorf("%s: expected %d got %d", t.Name(), testCase.StatusCode, w.Code)
		}
	}

	err = os.RemoveAll("static")
	if err != nil {
		t.Fatalf("%s: %s", t.Name(), err.Error())
	}
}

func TestRouter_FileServerOpen(t *testing.T) {
	err := os.Mkdir("static", 0777)
	if err != nil {
		t.Fatalf("%s: %s", t.Name(), err.Error())
	}

	err = ioutil.WriteFile("static/index.html", nil, 0777)
	if err != nil {
		t.Fatalf("%s: %s", t.Name(), err.Error())
	}

	cases := []struct {
		DirectoryPath string
		ErrorExpected bool
	}{
		{"static", false},
		{"bad_path", true},
	}

	for _, testCase := range cases {
		fs := FileSystem{http.Dir(testCase.DirectoryPath)}

		_, err := fs.Open(testCase.DirectoryPath)
		fmt.Println(err)
	}

	err = os.RemoveAll("static")
	if err != nil {
		t.Fatalf("%s: %s", t.Name(), err.Error())
	}
}

func TestRouter_TwoParam(t *testing.T) {
	r := New(false, nil, nil)
	r.Get("/users/:uid/files/:fid", func(w http.ResponseWriter, r *http.Request) {})

	req, _ := http.NewRequest("GET", "/users/1/files/1", nil)
	w := httptest.NewRecorder()

	h, err := r.find(req)
	if err != nil {
		t.Errorf("%s: unable to get handler: %s", t.Name(), err.Error())
	}

	h(w, req)

	if cyclops.Param(req, "uid") != "1" {
		t.Errorf(fmt.Sprintf("%s: params do not match", t.Name()))
	}

	if cyclops.Param(req, "fid") != "1" {
		t.Errorf(fmt.Sprintf("%s: params do not match", t.Name()))
	}
}

func TestRouterMicroParam(t *testing.T) {
	r := New(false, nil, nil)
	r.Get("/:a/:b/:c", func(w http.ResponseWriter, r *http.Request) {})

	req, _ := http.NewRequest("GET", "/1/2/3", nil)
	w := httptest.NewRecorder()

	h, err := r.find(req)
	if err != nil {
		t.Errorf("%s: unable to get handler: %s", t.Name(), err.Error())
	}

	h(w, req)

	if cyclops.Param(req, "a") != "1" {
		t.Errorf(fmt.Sprintf("%s: params do not match", t.Name()))
	}

	if cyclops.Param(req, "b") != "2" {
		t.Errorf(fmt.Sprintf("%s: params do not match", t.Name()))
	}

	if cyclops.Param(req, "c") != "3" {
		t.Errorf(fmt.Sprintf("%s: params do not match", t.Name()))
	}
}

func TestRouter_HttpMethods(t *testing.T) {
	cases := []struct {
		Method  string
		Handler http.HandlerFunc
	}{
		{http.MethodGet, func(writer http.ResponseWriter, request *http.Request) {}},
		{http.MethodPost, func(writer http.ResponseWriter, request *http.Request) {}},
		{http.MethodConnect, func(writer http.ResponseWriter, request *http.Request) {}},
		{http.MethodDelete, func(writer http.ResponseWriter, request *http.Request) {}},
		{http.MethodPatch, func(writer http.ResponseWriter, request *http.Request) {}},
		{http.MethodPut, func(writer http.ResponseWriter, request *http.Request) {}},
		{http.MethodTrace, func(writer http.ResponseWriter, request *http.Request) {}},
		{http.MethodHead, func(writer http.ResponseWriter, request *http.Request) {}},
		{http.MethodOptions, func(writer http.ResponseWriter, request *http.Request) {}},
	}

	for _, testCase := range cases {
		r := New(false, nil, nil)

		switch testCase.Method {
		case http.MethodGet:
			r.Get("/", testCase.Handler)
		case http.MethodPost:
			r.Post("/", testCase.Handler)
		case http.MethodConnect:
			r.Connect("/", testCase.Handler)
		case http.MethodDelete:
			r.Delete("/", testCase.Handler)
		case http.MethodPatch:
			r.Patch("/", testCase.Handler)
		case http.MethodPut:
			r.Put("/", testCase.Handler)
		case http.MethodTrace:
			r.Trace("/", testCase.Handler)
		case http.MethodHead:
			r.Head("/", testCase.Handler)
		case http.MethodOptions:
			r.Options("/", testCase.Handler)
		}

		req, _ := http.NewRequest(testCase.Method, "/", nil)

		h, err := r.find(req)
		if err != nil {
			t.Errorf("%s: unable to get handler: %s", t.Name(), err.Error())
		}

		if h == nil {
			t.Errorf("%s: nil handler", t.Name())
		}
	}
}

func TestRouter_UpdateRoute(t *testing.T) {
	r := New(false, nil, nil)

	r.Get("/users/hello", func(writer http.ResponseWriter, request *http.Request) {
		fmt.Fprintf(writer, "hello")
	})

	r.Get("/users/hello", func(writer http.ResponseWriter, request *http.Request) {
		fmt.Fprintf(writer, "hello world")
	})

	req, _ := http.NewRequest("GET", "/users/hello", nil)
	w := httptest.NewRecorder()

	h, err := r.find(req)
	if err != nil {
		t.Errorf("%s: unable to get handler: %s", t.Name(), err.Error())
	}

	h(w, req)

	if w.Body.String() != "hello world" {
		t.Errorf("%s: expected 'hello world' got %s", t.Name(), w.Body.String())
	}
}


// Benchmarks
type TestRoutes struct {
	method  string
	path    string
	handler http.HandlerFunc
}

var (
	api = []TestRoutes{
		// OAuth Authorizations
		{"GET", "/authorizations", func(writer http.ResponseWriter, request *http.Request) {}},
		{"GET", "/authorizations/:id", func(writer http.ResponseWriter, request *http.Request) {}},
		{"POST", "/authorizations", func(writer http.ResponseWriter, request *http.Request) {}},
		{"PUT", "/authorizations/clients/:client_id", func(writer http.ResponseWriter, request *http.Request) {}},
		{"PATCH", "/authorizations/:id", func(writer http.ResponseWriter, request *http.Request) {}},
		{"DELETE", "/authorizations/:id", func(writer http.ResponseWriter, request *http.Request) {}},
		{"GET", "/applications/:client_id/tokens/:access_token", func(writer http.ResponseWriter, request *http.Request) {}},
		{"DELETE", "/applications/:client_id/tokens", func(writer http.ResponseWriter, request *http.Request) {}},
		{"DELETE", "/applications/:client_id/tokens/:access_token", func(writer http.ResponseWriter, request *http.Request) {}},

		// Activity
		{"GET", "/events", func(writer http.ResponseWriter, request *http.Request) {}},
		{"GET", "/repos/:owner/:repo/events", func(writer http.ResponseWriter, request *http.Request) {}},
		{"GET", "/networks/:owner/:repo/events", func(writer http.ResponseWriter, request *http.Request) {}},
		{"GET", "/orgs/:org/events", func(writer http.ResponseWriter, request *http.Request) {}},
		{"GET", "/users/:user/received_events", func(writer http.ResponseWriter, request *http.Request) {}},
		{"GET", "/users/:user/received_events/public", func(writer http.ResponseWriter, request *http.Request) {}},
		{"GET", "/users/:user/events", func(writer http.ResponseWriter, request *http.Request) {}},
		{"GET", "/users/:user/events/public", func(writer http.ResponseWriter, request *http.Request) {}},
		{"GET", "/users/:user/events/orgs/:org", func(writer http.ResponseWriter, request *http.Request) {}},
		{"GET", "/feeds", func(writer http.ResponseWriter, request *http.Request) {}},
		{"GET", "/notifications", func(writer http.ResponseWriter, request *http.Request) {}},
		{"GET", "/repos/:owner/:repo/notifications", func(writer http.ResponseWriter, request *http.Request) {}},
		{"PUT", "/notifications", func(writer http.ResponseWriter, request *http.Request) {}},
		{"PUT", "/repos/:owner/:repo/notifications", func(writer http.ResponseWriter, request *http.Request) {}},
		{"GET", "/notifications/threads/:id", func(writer http.ResponseWriter, request *http.Request) {}},
		{"PATCH", "/notifications/threads/:id", func(writer http.ResponseWriter, request *http.Request) {}},
		{"GET", "/notifications/threads/:id/subscription", func(writer http.ResponseWriter, request *http.Request) {}},
		{"PUT", "/notifications/threads/:id/subscription", func(writer http.ResponseWriter, request *http.Request) {}},
		{"DELETE", "/notifications/threads/:id/subscription", func(writer http.ResponseWriter, request *http.Request) {}},
		{"GET", "/repos/:owner/:repo/stargazers", func(writer http.ResponseWriter, request *http.Request) {}},
		{"GET", "/users/:user/starred", func(writer http.ResponseWriter, request *http.Request) {}},
		{"GET", "/user/starred", func(writer http.ResponseWriter, request *http.Request) {}},
		{"GET", "/user/starred/:owner/:repo", func(writer http.ResponseWriter, request *http.Request) {}},
		{"PUT", "/user/starred/:owner/:repo", func(writer http.ResponseWriter, request *http.Request) {}},
		{"DELETE", "/user/starred/:owner/:repo", func(writer http.ResponseWriter, request *http.Request) {}},
		{"GET", "/repos/:owner/:repo/subscribers", func(writer http.ResponseWriter, request *http.Request) {}},
		{"GET", "/users/:user/subscriptions", func(writer http.ResponseWriter, request *http.Request) {}},
		{"GET", "/user/subscriptions", func(writer http.ResponseWriter, request *http.Request) {}},
		{"GET", "/repos/:owner/:repo/subscription", func(writer http.ResponseWriter, request *http.Request) {}},
		{"PUT", "/repos/:owner/:repo/subscription", func(writer http.ResponseWriter, request *http.Request) {}},
		{"DELETE", "/repos/:owner/:repo/subscription", func(writer http.ResponseWriter, request *http.Request) {}},
		{"GET", "/user/subscriptions/:owner/:repo", func(writer http.ResponseWriter, request *http.Request) {}},
		{"PUT", "/user/subscriptions/:owner/:repo", func(writer http.ResponseWriter, request *http.Request) {}},
		{"DELETE", "/user/subscriptions/:owner/:repo", func(writer http.ResponseWriter, request *http.Request) {}},

		// Gists
		{"GET", "/users/:user/gists", func(writer http.ResponseWriter, request *http.Request) {}},
		{"GET", "/gists", func(writer http.ResponseWriter, request *http.Request) {}},
		{"GET", "/gists/public", func(writer http.ResponseWriter, request *http.Request) {}},
		{"GET", "/gists/starred", func(writer http.ResponseWriter, request *http.Request) {}},
		{"GET", "/gists/:id", func(writer http.ResponseWriter, request *http.Request) {}},
		{"POST", "/gists", func(writer http.ResponseWriter, request *http.Request) {}},
		{"PATCH", "/gists/:id", func(writer http.ResponseWriter, request *http.Request) {}},
		{"PUT", "/gists/:id/star", func(writer http.ResponseWriter, request *http.Request) {}},
		{"DELETE", "/gists/:id/star", func(writer http.ResponseWriter, request *http.Request) {}},
		{"GET", "/gists/:id/star", func(writer http.ResponseWriter, request *http.Request) {}},
		{"POST", "/gists/:id/forks", func(writer http.ResponseWriter, request *http.Request) {}},
		{"DELETE", "/gists/:id", func(writer http.ResponseWriter, request *http.Request) {}},

		// Git Data
		{"GET", "/repos/:owner/:repo/git/blobs/:sha", func(writer http.ResponseWriter, request *http.Request) {}},
		{"POST", "/repos/:owner/:repo/git/blobs", func(writer http.ResponseWriter, request *http.Request) {}},
		{"GET", "/repos/:owner/:repo/git/commits/:sha", func(writer http.ResponseWriter, request *http.Request) {}},
		{"POST", "/repos/:owner/:repo/git/commits", func(writer http.ResponseWriter, request *http.Request) {}},
		{"GET", "/repos/:owner/:repo/git/refs/*ref", func(writer http.ResponseWriter, request *http.Request) {}},
		{"GET", "/repos/:owner/:repo/git/refs", func(writer http.ResponseWriter, request *http.Request) {}},
		{"POST", "/repos/:owner/:repo/git/refs", func(writer http.ResponseWriter, request *http.Request) {}},
		{"PATCH", "/repos/:owner/:repo/git/refs/*ref", func(writer http.ResponseWriter, request *http.Request) {}},
		{"DELETE", "/repos/:owner/:repo/git/refs/*ref", func(writer http.ResponseWriter, request *http.Request) {}},
		{"GET", "/repos/:owner/:repo/git/tags/:sha", func(writer http.ResponseWriter, request *http.Request) {}},
		{"POST", "/repos/:owner/:repo/git/tags", func(writer http.ResponseWriter, request *http.Request) {}},
		{"GET", "/repos/:owner/:repo/git/trees/:sha", func(writer http.ResponseWriter, request *http.Request) {}},
		{"POST", "/repos/:owner/:repo/git/trees", func(writer http.ResponseWriter, request *http.Request) {}},

		// Issues
		{"GET", "/issues", func(writer http.ResponseWriter, request *http.Request) {}},
		{"GET", "/user/issues", func(writer http.ResponseWriter, request *http.Request) {}},
		{"GET", "/orgs/:org/issues", func(writer http.ResponseWriter, request *http.Request) {}},
		{"GET", "/repos/:owner/:repo/issues", func(writer http.ResponseWriter, request *http.Request) {}},
		{"GET", "/repos/:owner/:repo/issues/:number", func(writer http.ResponseWriter, request *http.Request) {}},
		{"POST", "/repos/:owner/:repo/issues", func(writer http.ResponseWriter, request *http.Request) {}},
		{"PATCH", "/repos/:owner/:repo/issues/:number", func(writer http.ResponseWriter, request *http.Request) {}},
		{"GET", "/repos/:owner/:repo/assignees", func(writer http.ResponseWriter, request *http.Request) {}},
		{"GET", "/repos/:owner/:repo/assignees/:assignee", func(writer http.ResponseWriter, request *http.Request) {}},
		{"GET", "/repos/:owner/:repo/issues/:number/comments", func(writer http.ResponseWriter, request *http.Request) {}},
		{"GET", "/repos/:owner/:repo/issues/comments", func(writer http.ResponseWriter, request *http.Request) {}},
		{"GET", "/repos/:owner/:repo/issues/comments/:id", func(writer http.ResponseWriter, request *http.Request) {}},
		{"POST", "/repos/:owner/:repo/issues/:number/comments", func(writer http.ResponseWriter, request *http.Request) {}},
		{"PATCH", "/repos/:owner/:repo/issues/comments/:id", func(writer http.ResponseWriter, request *http.Request) {}},
		{"DELETE", "/repos/:owner/:repo/issues/comments/:id", func(writer http.ResponseWriter, request *http.Request) {}},
		{"GET", "/repos/:owner/:repo/issues/:number/events", func(writer http.ResponseWriter, request *http.Request) {}},
		{"GET", "/repos/:owner/:repo/issues/events", func(writer http.ResponseWriter, request *http.Request) {}},
		{"GET", "/repos/:owner/:repo/issues/events/:id", func(writer http.ResponseWriter, request *http.Request) {}},
		{"GET", "/repos/:owner/:repo/labels", func(writer http.ResponseWriter, request *http.Request) {}},
		{"GET", "/repos/:owner/:repo/labels/:name", func(writer http.ResponseWriter, request *http.Request) {}},
		{"POST", "/repos/:owner/:repo/labels", func(writer http.ResponseWriter, request *http.Request) {}},
		{"PATCH", "/repos/:owner/:repo/labels/:name", func(writer http.ResponseWriter, request *http.Request) {}},
		{"DELETE", "/repos/:owner/:repo/labels/:name", func(writer http.ResponseWriter, request *http.Request) {}},
		{"GET", "/repos/:owner/:repo/issues/:number/labels", func(writer http.ResponseWriter, request *http.Request) {}},
		{"POST", "/repos/:owner/:repo/issues/:number/labels", func(writer http.ResponseWriter, request *http.Request) {}},
		{"DELETE", "/repos/:owner/:repo/issues/:number/labels/:name", func(writer http.ResponseWriter, request *http.Request) {}},
		{"PUT", "/repos/:owner/:repo/issues/:number/labels", func(writer http.ResponseWriter, request *http.Request) {}},
		{"DELETE", "/repos/:owner/:repo/issues/:number/labels", func(writer http.ResponseWriter, request *http.Request) {}},
		{"GET", "/repos/:owner/:repo/milestones/:number/labels", func(writer http.ResponseWriter, request *http.Request) {}},
		{"GET", "/repos/:owner/:repo/milestones", func(writer http.ResponseWriter, request *http.Request) {}},
		{"GET", "/repos/:owner/:repo/milestones/:number", func(writer http.ResponseWriter, request *http.Request) {}},
		{"POST", "/repos/:owner/:repo/milestones", func(writer http.ResponseWriter, request *http.Request) {}},
		{"PATCH", "/repos/:owner/:repo/milestones/:number", func(writer http.ResponseWriter, request *http.Request) {}},
		{"DELETE", "/repos/:owner/:repo/milestones/:number", func(writer http.ResponseWriter, request *http.Request) {}},

		// Miscellaneous
		{"GET", "/emojis", func(writer http.ResponseWriter, request *http.Request) {}},
		{"GET", "/gitignore/templates", func(writer http.ResponseWriter, request *http.Request) {}},
		{"GET", "/gitignore/templates/:name", func(writer http.ResponseWriter, request *http.Request) {}},
		{"POST", "/markdown", func(writer http.ResponseWriter, request *http.Request) {}},
		{"POST", "/markdown/raw", func(writer http.ResponseWriter, request *http.Request) {}},
		{"GET", "/meta", func(writer http.ResponseWriter, request *http.Request) {}},
		{"GET", "/rate_limit", func(writer http.ResponseWriter, request *http.Request) {}},

		// Organizations
		{"GET", "/users/:user/orgs", func(writer http.ResponseWriter, request *http.Request) {}},
		{"GET", "/user/orgs", func(writer http.ResponseWriter, request *http.Request) {}},
		{"GET", "/orgs/:org", func(writer http.ResponseWriter, request *http.Request) {}},
		{"PATCH", "/orgs/:org", func(writer http.ResponseWriter, request *http.Request) {}},
		{"GET", "/orgs/:org/members", func(writer http.ResponseWriter, request *http.Request) {}},
		{"GET", "/orgs/:org/members/:user", func(writer http.ResponseWriter, request *http.Request) {}},
		{"DELETE", "/orgs/:org/members/:user", func(writer http.ResponseWriter, request *http.Request) {}},
		{"GET", "/orgs/:org/public_members", func(writer http.ResponseWriter, request *http.Request) {}},
		{"GET", "/orgs/:org/public_members/:user", func(writer http.ResponseWriter, request *http.Request) {}},
		{"PUT", "/orgs/:org/public_members/:user", func(writer http.ResponseWriter, request *http.Request) {}},
		{"DELETE", "/orgs/:org/public_members/:user", func(writer http.ResponseWriter, request *http.Request) {}},
		{"GET", "/orgs/:org/teams", func(writer http.ResponseWriter, request *http.Request) {}},
		{"GET", "/teams/:id", func(writer http.ResponseWriter, request *http.Request) {}},
		{"POST", "/orgs/:org/teams", func(writer http.ResponseWriter, request *http.Request) {}},
		{"PATCH", "/teams/:id", func(writer http.ResponseWriter, request *http.Request) {}},
		{"DELETE", "/teams/:id", func(writer http.ResponseWriter, request *http.Request) {}},
		{"GET", "/teams/:id/members", func(writer http.ResponseWriter, request *http.Request) {}},
		{"GET", "/teams/:id/members/:user", func(writer http.ResponseWriter, request *http.Request) {}},
		{"PUT", "/teams/:id/members/:user", func(writer http.ResponseWriter, request *http.Request) {}},
		{"DELETE", "/teams/:id/members/:user", func(writer http.ResponseWriter, request *http.Request) {}},
		{"GET", "/teams/:id/repos", func(writer http.ResponseWriter, request *http.Request) {}},
		{"GET", "/teams/:id/repos/:owner/:repo", func(writer http.ResponseWriter, request *http.Request) {}},
		{"PUT", "/teams/:id/repos/:owner/:repo", func(writer http.ResponseWriter, request *http.Request) {}},
		{"DELETE", "/teams/:id/repos/:owner/:repo", func(writer http.ResponseWriter, request *http.Request) {}},
		{"GET", "/user/teams", func(writer http.ResponseWriter, request *http.Request) {}},

		// Pull Requests
		{"GET", "/repos/:owner/:repo/pulls", func(writer http.ResponseWriter, request *http.Request) {}},
		{"GET", "/repos/:owner/:repo/pulls/:number", func(writer http.ResponseWriter, request *http.Request) {}},
		{"POST", "/repos/:owner/:repo/pulls", func(writer http.ResponseWriter, request *http.Request) {}},
		{"PATCH", "/repos/:owner/:repo/pulls/:number", func(writer http.ResponseWriter, request *http.Request) {}},
		{"GET", "/repos/:owner/:repo/pulls/:number/commits", func(writer http.ResponseWriter, request *http.Request) {}},
		{"GET", "/repos/:owner/:repo/pulls/:number/files", func(writer http.ResponseWriter, request *http.Request) {}},
		{"GET", "/repos/:owner/:repo/pulls/:number/merge", func(writer http.ResponseWriter, request *http.Request) {}},
		{"PUT", "/repos/:owner/:repo/pulls/:number/merge", func(writer http.ResponseWriter, request *http.Request) {}},
		{"GET", "/repos/:owner/:repo/pulls/:number/comments", func(writer http.ResponseWriter, request *http.Request) {}},
		{"GET", "/repos/:owner/:repo/pulls/comments", func(writer http.ResponseWriter, request *http.Request) {}},
		{"GET", "/repos/:owner/:repo/pulls/comments/:number", func(writer http.ResponseWriter, request *http.Request) {}},
		{"PUT", "/repos/:owner/:repo/pulls/:number/comments", func(writer http.ResponseWriter, request *http.Request) {}},
		{"PATCH", "/repos/:owner/:repo/pulls/comments/:number", func(writer http.ResponseWriter, request *http.Request) {}},
		{"DELETE", "/repos/:owner/:repo/pulls/comments/:number", func(writer http.ResponseWriter, request *http.Request) {}},

		// Repositories
		{"GET", "/user/repos", func(writer http.ResponseWriter, request *http.Request) {}},
		{"GET", "/users/:user/repos", func(writer http.ResponseWriter, request *http.Request) {}},
		{"GET", "/orgs/:org/repos", func(writer http.ResponseWriter, request *http.Request) {}},
		{"GET", "/repositories", func(writer http.ResponseWriter, request *http.Request) {}},
		{"POST", "/user/repos", func(writer http.ResponseWriter, request *http.Request) {}},
		{"POST", "/orgs/:org/repos", func(writer http.ResponseWriter, request *http.Request) {}},
		{"GET", "/repos/:owner/:repo", func(writer http.ResponseWriter, request *http.Request) {}},
		{"PATCH", "/repos/:owner/:repo", func(writer http.ResponseWriter, request *http.Request) {}},
		{"GET", "/repos/:owner/:repo/contributors", func(writer http.ResponseWriter, request *http.Request) {}},
		{"GET", "/repos/:owner/:repo/languages", func(writer http.ResponseWriter, request *http.Request) {}},
		{"GET", "/repos/:owner/:repo/teams", func(writer http.ResponseWriter, request *http.Request) {}},
		{"GET", "/repos/:owner/:repo/tags", func(writer http.ResponseWriter, request *http.Request) {}},
		{"GET", "/repos/:owner/:repo/branches", func(writer http.ResponseWriter, request *http.Request) {}},
		{"GET", "/repos/:owner/:repo/branches/:branch", func(writer http.ResponseWriter, request *http.Request) {}},
		{"DELETE", "/repos/:owner/:repo", func(writer http.ResponseWriter, request *http.Request) {}},
		{"GET", "/repos/:owner/:repo/collaborators", func(writer http.ResponseWriter, request *http.Request) {}},
		{"GET", "/repos/:owner/:repo/collaborators/:user", func(writer http.ResponseWriter, request *http.Request) {}},
		{"PUT", "/repos/:owner/:repo/collaborators/:user", func(writer http.ResponseWriter, request *http.Request) {}},
		{"DELETE", "/repos/:owner/:repo/collaborators/:user", func(writer http.ResponseWriter, request *http.Request) {}},
		{"GET", "/repos/:owner/:repo/comments", func(writer http.ResponseWriter, request *http.Request) {}},
		{"GET", "/repos/:owner/:repo/commits/:sha/comments", func(writer http.ResponseWriter, request *http.Request) {}},
		{"POST", "/repos/:owner/:repo/commits/:sha/comments", func(writer http.ResponseWriter, request *http.Request) {}},
		{"GET", "/repos/:owner/:repo/comments/:id", func(writer http.ResponseWriter, request *http.Request) {}},
		{"PATCH", "/repos/:owner/:repo/comments/:id", func(writer http.ResponseWriter, request *http.Request) {}},
		{"DELETE", "/repos/:owner/:repo/comments/:id", func(writer http.ResponseWriter, request *http.Request) {}},
		{"GET", "/repos/:owner/:repo/commits", func(writer http.ResponseWriter, request *http.Request) {}},
		{"GET", "/repos/:owner/:repo/commits/:sha", func(writer http.ResponseWriter, request *http.Request) {}},
		{"GET", "/repos/:owner/:repo/readme", func(writer http.ResponseWriter, request *http.Request) {}},
		{"GET", "/repos/:owner/:repo/contents/*path", func(writer http.ResponseWriter, request *http.Request) {}},
		{"PUT", "/repos/:owner/:repo/contents/*path", func(writer http.ResponseWriter, request *http.Request) {}},
		{"DELETE", "/repos/:owner/:repo/contents/*path", func(writer http.ResponseWriter, request *http.Request) {}},
		{"GET", "/repos/:owner/:repo/:archive_format/:ref", func(writer http.ResponseWriter, request *http.Request) {}},
		{"GET", "/repos/:owner/:repo/keys", func(writer http.ResponseWriter, request *http.Request) {}},
		{"GET", "/repos/:owner/:repo/keys/:id", func(writer http.ResponseWriter, request *http.Request) {}},
		{"POST", "/repos/:owner/:repo/keys", func(writer http.ResponseWriter, request *http.Request) {}},
		{"PATCH", "/repos/:owner/:repo/keys/:id", func(writer http.ResponseWriter, request *http.Request) {}},
		{"DELETE", "/repos/:owner/:repo/keys/:id", func(writer http.ResponseWriter, request *http.Request) {}},
		{"GET", "/repos/:owner/:repo/downloads", func(writer http.ResponseWriter, request *http.Request) {}},
		{"GET", "/repos/:owner/:repo/downloads/:id", func(writer http.ResponseWriter, request *http.Request) {}},
		{"DELETE", "/repos/:owner/:repo/downloads/:id", func(writer http.ResponseWriter, request *http.Request) {}},
		{"GET", "/repos/:owner/:repo/forks", func(writer http.ResponseWriter, request *http.Request) {}},
		{"POST", "/repos/:owner/:repo/forks", func(writer http.ResponseWriter, request *http.Request) {}},
		{"GET", "/repos/:owner/:repo/hooks", func(writer http.ResponseWriter, request *http.Request) {}},
		{"GET", "/repos/:owner/:repo/hooks/:id", func(writer http.ResponseWriter, request *http.Request) {}},
		{"POST", "/repos/:owner/:repo/hooks", func(writer http.ResponseWriter, request *http.Request) {}},
		{"PATCH", "/repos/:owner/:repo/hooks/:id", func(writer http.ResponseWriter, request *http.Request) {}},
		{"POST", "/repos/:owner/:repo/hooks/:id/tests", func(writer http.ResponseWriter, request *http.Request) {}},
		{"DELETE", "/repos/:owner/:repo/hooks/:id", func(writer http.ResponseWriter, request *http.Request) {}},
		{"POST", "/repos/:owner/:repo/merges", func(writer http.ResponseWriter, request *http.Request) {}},
		{"GET", "/repos/:owner/:repo/releases", func(writer http.ResponseWriter, request *http.Request) {}},
		{"GET", "/repos/:owner/:repo/releases/:id", func(writer http.ResponseWriter, request *http.Request) {}},
		{"POST", "/repos/:owner/:repo/releases", func(writer http.ResponseWriter, request *http.Request) {}},
		{"PATCH", "/repos/:owner/:repo/releases/:id", func(writer http.ResponseWriter, request *http.Request) {}},
		{"DELETE", "/repos/:owner/:repo/releases/:id", func(writer http.ResponseWriter, request *http.Request) {}},
		{"GET", "/repos/:owner/:repo/releases/:id/assets", func(writer http.ResponseWriter, request *http.Request) {}},
		{"GET", "/repos/:owner/:repo/stats/contributors", func(writer http.ResponseWriter, request *http.Request) {}},
		{"GET", "/repos/:owner/:repo/stats/commit_activity", func(writer http.ResponseWriter, request *http.Request) {}},
		{"GET", "/repos/:owner/:repo/stats/code_frequency", func(writer http.ResponseWriter, request *http.Request) {}},
		{"GET", "/repos/:owner/:repo/stats/participation", func(writer http.ResponseWriter, request *http.Request) {}},
		{"GET", "/repos/:owner/:repo/stats/punch_card", func(writer http.ResponseWriter, request *http.Request) {}},
		{"GET", "/repos/:owner/:repo/statuses/:ref", func(writer http.ResponseWriter, request *http.Request) {}},
		{"POST", "/repos/:owner/:repo/statuses/:ref", func(writer http.ResponseWriter, request *http.Request) {}},

		// Search
		{"GET", "/search/repositories", func(writer http.ResponseWriter, request *http.Request) {}},
		{"GET", "/search/code", func(writer http.ResponseWriter, request *http.Request) {}},
		{"GET", "/search/issues", func(writer http.ResponseWriter, request *http.Request) {}},
		{"GET", "/search/users", func(writer http.ResponseWriter, request *http.Request) {}},
		{"GET", "/legacy/issues/search/:owner/:repository/:state/:keyword", func(writer http.ResponseWriter, request *http.Request) {}},
		{"GET", "/legacy/repos/search/:keyword", func(writer http.ResponseWriter, request *http.Request) {}},
		{"GET", "/legacy/user/search/:keyword", func(writer http.ResponseWriter, request *http.Request) {}},
		{"GET", "/legacy/user/email/:email", func(writer http.ResponseWriter, request *http.Request) {}},

		// Users
		{"GET", "/users/:user", func(writer http.ResponseWriter, request *http.Request) {}},
		{"GET", "/user", func(writer http.ResponseWriter, request *http.Request) {}},
		{"PATCH", "/user", func(writer http.ResponseWriter, request *http.Request) {}},
		{"GET", "/users", func(writer http.ResponseWriter, request *http.Request) {}},
		{"GET", "/user/emails", func(writer http.ResponseWriter, request *http.Request) {}},
		{"POST", "/user/emails", func(writer http.ResponseWriter, request *http.Request) {}},
		{"DELETE", "/user/emails", func(writer http.ResponseWriter, request *http.Request) {}},
		{"GET", "/users/:user/followers", func(writer http.ResponseWriter, request *http.Request) {}},
		{"GET", "/user/followers", func(writer http.ResponseWriter, request *http.Request) {}},
		{"GET", "/users/:user/following", func(writer http.ResponseWriter, request *http.Request) {}},
		{"GET", "/user/following", func(writer http.ResponseWriter, request *http.Request) {}},
		{"GET", "/user/following/:user", func(writer http.ResponseWriter, request *http.Request) {}},
		{"GET", "/users/:user/following/:target_user", func(writer http.ResponseWriter, request *http.Request) {}},
		{"PUT", "/user/following/:user", func(writer http.ResponseWriter, request *http.Request) {}},
		{"DELETE", "/user/following/:user", func(writer http.ResponseWriter, request *http.Request) {}},
		{"GET", "/users/:user/keys", func(writer http.ResponseWriter, request *http.Request) {}},
		{"GET", "/user/keys", func(writer http.ResponseWriter, request *http.Request) {}},
		{"GET", "/user/keys/:id", func(writer http.ResponseWriter, request *http.Request) {}},
		{"POST", "/user/keys", func(writer http.ResponseWriter, request *http.Request) {}},
		{"PATCH", "/user/keys/:id", func(writer http.ResponseWriter, request *http.Request) {}},
		{"DELETE", "/user/keys/:id", func(writer http.ResponseWriter, request *http.Request) {}},
	}
)

func Benchmark_AddRoute(b *testing.B) {
	b.StopTimer()
	b.ReportAllocs()

	b.StartTimer()

	for _, testCase := range api {
		r := New(false, nil, nil)

		switch testCase.method {
		case http.MethodGet:
			r.Get(testCase.path, testCase.handler)
		case http.MethodPost:
			r.Post(testCase.path, testCase.handler)
		case http.MethodConnect:
			r.Connect(testCase.path, testCase.handler)
		case http.MethodDelete:
			r.Delete(testCase.path, testCase.handler)
		case http.MethodPatch:
			r.Patch(testCase.path, testCase.handler)
		case http.MethodPut:
			r.Put(testCase.path, testCase.handler)
		case http.MethodTrace:
			r.Trace(testCase.path, testCase.handler)
		case http.MethodHead:
			r.Head(testCase.path, testCase.handler)
		case http.MethodOptions:
			r.Options(testCase.path, testCase.handler)
		}
	}
}

func Benchmark_GetRoute(b *testing.B) {
	b.StopTimer()
	b.ReportAllocs()

	r := New(false, nil, nil)

	for _, testCase := range api {
		switch testCase.method {
		case http.MethodGet:
			r.Get(testCase.path, testCase.handler)
		case http.MethodPost:
			r.Post(testCase.path, testCase.handler)
		case http.MethodConnect:
			r.Connect(testCase.path, testCase.handler)
		case http.MethodDelete:
			r.Delete(testCase.path, testCase.handler)
		case http.MethodPatch:
			r.Patch(testCase.path, testCase.handler)
		case http.MethodPut:
			r.Put(testCase.path, testCase.handler)
		case http.MethodTrace:
			r.Trace(testCase.path, testCase.handler)
		case http.MethodHead:
			r.Head(testCase.path, testCase.handler)
		case http.MethodOptions:
			r.Options(testCase.path, testCase.handler)
		}
	}

	params := url.Values{}

	rand.Seed(time.Now().UnixNano())
	min := 0
	max := len(api) - 1

	b.StartTimer()
	for i := 0; i < b.N; i++ {

		randomIdx := rand.Intn(max - min + 1) + min

		r.tree.traverse(strings.Split(api[randomIdx].path, "/"), params)
	}
}