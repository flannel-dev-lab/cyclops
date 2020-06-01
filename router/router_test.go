package router

import (
	"fmt"
	"github.com/flannel-dev-lab/cyclops"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

type route struct {
	method  string
	path    string
	handler http.HandlerFunc
}

var (
	api = []route{
		// OAuth Authorizations
		{"GET", "/authorizations", nil},
		{"GET", "/authorizations/:id", nil},
		{"POST", "/authorizations", nil},
		{"PUT", "/authorizations/clients/:client_id", nil},
		{"PATCH", "/authorizations/:id", nil},
		{"DELETE", "/authorizations/:id", nil},
		{"GET", "/applications/:client_id/tokens/:access_token", nil},
		{"DELETE", "/applications/:client_id/tokens", nil},
		{"DELETE", "/applications/:client_id/tokens/:access_token", nil},

		// Activity
		{"GET", "/events", nil},
		{"GET", "/repos/:owner/:repo/events", nil},
		{"GET", "/networks/:owner/:repo/events", nil},
		{"GET", "/orgs/:org/events", nil},
		{"GET", "/users/:user/received_events", nil},
		{"GET", "/users/:user/received_events/public", nil},
		{"GET", "/users/:user/events", nil},
		{"GET", "/users/:user/events/public", nil},
		{"GET", "/users/:user/events/orgs/:org", nil},
		{"GET", "/feeds", nil},
		{"GET", "/notifications", nil},
		{"GET", "/repos/:owner/:repo/notifications", nil},
		{"PUT", "/notifications", nil},
		{"PUT", "/repos/:owner/:repo/notifications", nil},
		{"GET", "/notifications/threads/:id", nil},
		{"PATCH", "/notifications/threads/:id", nil},
		{"GET", "/notifications/threads/:id/subscription", nil},
		{"PUT", "/notifications/threads/:id/subscription", nil},
		{"DELETE", "/notifications/threads/:id/subscription", nil},
		{"GET", "/repos/:owner/:repo/stargazers", nil},
		{"GET", "/users/:user/starred", nil},
		{"GET", "/user/starred", nil},
		{"GET", "/user/starred/:owner/:repo", nil},
		{"PUT", "/user/starred/:owner/:repo", nil},
		{"DELETE", "/user/starred/:owner/:repo", nil},
		{"GET", "/repos/:owner/:repo/subscribers", nil},
		{"GET", "/users/:user/subscriptions", nil},
		{"GET", "/user/subscriptions", nil},
		{"GET", "/repos/:owner/:repo/subscription", nil},
		{"PUT", "/repos/:owner/:repo/subscription", nil},
		{"DELETE", "/repos/:owner/:repo/subscription", nil},
		{"GET", "/user/subscriptions/:owner/:repo", nil},
		{"PUT", "/user/subscriptions/:owner/:repo", nil},
		{"DELETE", "/user/subscriptions/:owner/:repo", nil},

		// Gists
		{"GET", "/users/:user/gists", nil},
		{"GET", "/gists", nil},
		{"GET", "/gists/public", nil},
		{"GET", "/gists/starred", nil},
		{"GET", "/gists/:id", nil},
		{"POST", "/gists", nil},
		{"PATCH", "/gists/:id", nil},
		{"PUT", "/gists/:id/star", nil},
		{"DELETE", "/gists/:id/star", nil},
		{"GET", "/gists/:id/star", nil},
		{"POST", "/gists/:id/forks", nil},
		{"DELETE", "/gists/:id", nil},

		// Git Data
		{"GET", "/repos/:owner/:repo/git/blobs/:sha", nil},
		{"POST", "/repos/:owner/:repo/git/blobs", nil},
		{"GET", "/repos/:owner/:repo/git/commits/:sha", nil},
		{"POST", "/repos/:owner/:repo/git/commits", nil},
		{"GET", "/repos/:owner/:repo/git/refs/*ref", nil},
		{"GET", "/repos/:owner/:repo/git/refs", nil},
		{"POST", "/repos/:owner/:repo/git/refs", nil},
		{"PATCH", "/repos/:owner/:repo/git/refs/*ref", nil},
		{"DELETE", "/repos/:owner/:repo/git/refs/*ref", nil},
		{"GET", "/repos/:owner/:repo/git/tags/:sha", nil},
		{"POST", "/repos/:owner/:repo/git/tags", nil},
		{"GET", "/repos/:owner/:repo/git/trees/:sha", nil},
		{"POST", "/repos/:owner/:repo/git/trees", nil},

		// Issues
		{"GET", "/issues", nil},
		{"GET", "/user/issues", nil},
		{"GET", "/orgs/:org/issues", nil},
		{"GET", "/repos/:owner/:repo/issues", nil},
		{"GET", "/repos/:owner/:repo/issues/:number", nil},
		{"POST", "/repos/:owner/:repo/issues", nil},
		{"PATCH", "/repos/:owner/:repo/issues/:number", nil},
		{"GET", "/repos/:owner/:repo/assignees", nil},
		{"GET", "/repos/:owner/:repo/assignees/:assignee", nil},
		{"GET", "/repos/:owner/:repo/issues/:number/comments", nil},
		{"GET", "/repos/:owner/:repo/issues/comments", nil},
		{"GET", "/repos/:owner/:repo/issues/comments/:id", nil},
		{"POST", "/repos/:owner/:repo/issues/:number/comments", nil},
		{"PATCH", "/repos/:owner/:repo/issues/comments/:id", nil},
		{"DELETE", "/repos/:owner/:repo/issues/comments/:id", nil},
		{"GET", "/repos/:owner/:repo/issues/:number/events", nil},
		{"GET", "/repos/:owner/:repo/issues/events", nil},
		{"GET", "/repos/:owner/:repo/issues/events/:id", nil},
		{"GET", "/repos/:owner/:repo/labels", nil},
		{"GET", "/repos/:owner/:repo/labels/:name", nil},
		{"POST", "/repos/:owner/:repo/labels", nil},
		{"PATCH", "/repos/:owner/:repo/labels/:name", nil},
		{"DELETE", "/repos/:owner/:repo/labels/:name", nil},
		{"GET", "/repos/:owner/:repo/issues/:number/labels", nil},
		{"POST", "/repos/:owner/:repo/issues/:number/labels", nil},
		{"DELETE", "/repos/:owner/:repo/issues/:number/labels/:name", nil},
		{"PUT", "/repos/:owner/:repo/issues/:number/labels", nil},
		{"DELETE", "/repos/:owner/:repo/issues/:number/labels", nil},
		{"GET", "/repos/:owner/:repo/milestones/:number/labels", nil},
		{"GET", "/repos/:owner/:repo/milestones", nil},
		{"GET", "/repos/:owner/:repo/milestones/:number", nil},
		{"POST", "/repos/:owner/:repo/milestones", nil},
		{"PATCH", "/repos/:owner/:repo/milestones/:number", nil},
		{"DELETE", "/repos/:owner/:repo/milestones/:number", nil},

		// Miscellaneous
		{"GET", "/emojis", nil},
		{"GET", "/gitignore/templates", nil},
		{"GET", "/gitignore/templates/:name", nil},
		{"POST", "/markdown", nil},
		{"POST", "/markdown/raw", nil},
		{"GET", "/meta", nil},
		{"GET", "/rate_limit", nil},

		// Organizations
		{"GET", "/users/:user/orgs", nil},
		{"GET", "/user/orgs", nil},
		{"GET", "/orgs/:org", nil},
		{"PATCH", "/orgs/:org", nil},
		{"GET", "/orgs/:org/members", nil},
		{"GET", "/orgs/:org/members/:user", nil},
		{"DELETE", "/orgs/:org/members/:user", nil},
		{"GET", "/orgs/:org/public_members", nil},
		{"GET", "/orgs/:org/public_members/:user", nil},
		{"PUT", "/orgs/:org/public_members/:user", nil},
		{"DELETE", "/orgs/:org/public_members/:user", nil},
		{"GET", "/orgs/:org/teams", nil},
		{"GET", "/teams/:id", nil},
		{"POST", "/orgs/:org/teams", nil},
		{"PATCH", "/teams/:id", nil},
		{"DELETE", "/teams/:id", nil},
		{"GET", "/teams/:id/members", nil},
		{"GET", "/teams/:id/members/:user", nil},
		{"PUT", "/teams/:id/members/:user", nil},
		{"DELETE", "/teams/:id/members/:user", nil},
		{"GET", "/teams/:id/repos", nil},
		{"GET", "/teams/:id/repos/:owner/:repo", nil},
		{"PUT", "/teams/:id/repos/:owner/:repo", nil},
		{"DELETE", "/teams/:id/repos/:owner/:repo", nil},
		{"GET", "/user/teams", nil},

		// Pull Requests
		{"GET", "/repos/:owner/:repo/pulls", nil},
		{"GET", "/repos/:owner/:repo/pulls/:number", nil},
		{"POST", "/repos/:owner/:repo/pulls", nil},
		{"PATCH", "/repos/:owner/:repo/pulls/:number", nil},
		{"GET", "/repos/:owner/:repo/pulls/:number/commits", nil},
		{"GET", "/repos/:owner/:repo/pulls/:number/files", nil},
		{"GET", "/repos/:owner/:repo/pulls/:number/merge", nil},
		{"PUT", "/repos/:owner/:repo/pulls/:number/merge", nil},
		{"GET", "/repos/:owner/:repo/pulls/:number/comments", nil},
		{"GET", "/repos/:owner/:repo/pulls/comments", nil},
		{"GET", "/repos/:owner/:repo/pulls/comments/:number", nil},
		{"PUT", "/repos/:owner/:repo/pulls/:number/comments", nil},
		{"PATCH", "/repos/:owner/:repo/pulls/comments/:number", nil},
		{"DELETE", "/repos/:owner/:repo/pulls/comments/:number", nil},

		// Repositories
		{"GET", "/user/repos", nil},
		{"GET", "/users/:user/repos", nil},
		{"GET", "/orgs/:org/repos", nil},
		{"GET", "/repositories", nil},
		{"POST", "/user/repos", nil},
		{"POST", "/orgs/:org/repos", nil},
		{"GET", "/repos/:owner/:repo", nil},
		{"PATCH", "/repos/:owner/:repo", nil},
		{"GET", "/repos/:owner/:repo/contributors", nil},
		{"GET", "/repos/:owner/:repo/languages", nil},
		{"GET", "/repos/:owner/:repo/teams", nil},
		{"GET", "/repos/:owner/:repo/tags", nil},
		{"GET", "/repos/:owner/:repo/branches", nil},
		{"GET", "/repos/:owner/:repo/branches/:branch", nil},
		{"DELETE", "/repos/:owner/:repo", nil},
		{"GET", "/repos/:owner/:repo/collaborators", nil},
		{"GET", "/repos/:owner/:repo/collaborators/:user", nil},
		{"PUT", "/repos/:owner/:repo/collaborators/:user", nil},
		{"DELETE", "/repos/:owner/:repo/collaborators/:user", nil},
		{"GET", "/repos/:owner/:repo/comments", nil},
		{"GET", "/repos/:owner/:repo/commits/:sha/comments", nil},
		{"POST", "/repos/:owner/:repo/commits/:sha/comments", nil},
		{"GET", "/repos/:owner/:repo/comments/:id", nil},
		{"PATCH", "/repos/:owner/:repo/comments/:id", nil},
		{"DELETE", "/repos/:owner/:repo/comments/:id", nil},
		{"GET", "/repos/:owner/:repo/commits", nil},
		{"GET", "/repos/:owner/:repo/commits/:sha", nil},
		{"GET", "/repos/:owner/:repo/readme", nil},
		{"GET", "/repos/:owner/:repo/contents/*path", nil},
		{"PUT", "/repos/:owner/:repo/contents/*path", nil},
		{"DELETE", "/repos/:owner/:repo/contents/*path", nil},
		{"GET", "/repos/:owner/:repo/:archive_format/:ref", nil},
		{"GET", "/repos/:owner/:repo/keys", nil},
		{"GET", "/repos/:owner/:repo/keys/:id", nil},
		{"POST", "/repos/:owner/:repo/keys", nil},
		{"PATCH", "/repos/:owner/:repo/keys/:id", nil},
		{"DELETE", "/repos/:owner/:repo/keys/:id", nil},
		{"GET", "/repos/:owner/:repo/downloads", nil},
		{"GET", "/repos/:owner/:repo/downloads/:id", nil},
		{"DELETE", "/repos/:owner/:repo/downloads/:id", nil},
		{"GET", "/repos/:owner/:repo/forks", nil},
		{"POST", "/repos/:owner/:repo/forks", nil},
		{"GET", "/repos/:owner/:repo/hooks", nil},
		{"GET", "/repos/:owner/:repo/hooks/:id", nil},
		{"POST", "/repos/:owner/:repo/hooks", nil},
		{"PATCH", "/repos/:owner/:repo/hooks/:id", nil},
		{"POST", "/repos/:owner/:repo/hooks/:id/tests", nil},
		{"DELETE", "/repos/:owner/:repo/hooks/:id", nil},
		{"POST", "/repos/:owner/:repo/merges", nil},
		{"GET", "/repos/:owner/:repo/releases", nil},
		{"GET", "/repos/:owner/:repo/releases/:id", nil},
		{"POST", "/repos/:owner/:repo/releases", nil},
		{"PATCH", "/repos/:owner/:repo/releases/:id", nil},
		{"DELETE", "/repos/:owner/:repo/releases/:id", nil},
		{"GET", "/repos/:owner/:repo/releases/:id/assets", nil},
		{"GET", "/repos/:owner/:repo/stats/contributors", nil},
		{"GET", "/repos/:owner/:repo/stats/commit_activity", nil},
		{"GET", "/repos/:owner/:repo/stats/code_frequency", nil},
		{"GET", "/repos/:owner/:repo/stats/participation", nil},
		{"GET", "/repos/:owner/:repo/stats/punch_card", nil},
		{"GET", "/repos/:owner/:repo/statuses/:ref", nil},
		{"POST", "/repos/:owner/:repo/statuses/:ref", nil},

		// Search
		{"GET", "/search/repositories", nil},
		{"GET", "/search/code", nil},
		{"GET", "/search/issues", nil},
		{"GET", "/search/users", nil},
		{"GET", "/legacy/issues/search/:owner/:repository/:state/:keyword", nil},
		{"GET", "/legacy/repos/search/:keyword", nil},
		{"GET", "/legacy/user/search/:keyword", nil},
		{"GET", "/legacy/user/email/:email", nil},

		// Users
		{"GET", "/users/:user", nil},
		{"GET", "/user", nil},
		{"PATCH", "/user", nil},
		{"GET", "/users", nil},
		{"GET", "/user/emails", nil},
		{"POST", "/user/emails", nil},
		{"DELETE", "/user/emails", nil},
		{"GET", "/users/:user/followers", nil},
		{"GET", "/user/followers", nil},
		{"GET", "/users/:user/following", nil},
		{"GET", "/user/following", nil},
		{"GET", "/user/following/:user", nil},
		{"GET", "/users/:user/following/:target_user", nil},
		{"PUT", "/user/following/:user", nil},
		{"DELETE", "/user/following/:user", nil},
		{"GET", "/users/:user/keys", nil},
		{"GET", "/user/keys", nil},
		{"GET", "/user/keys/:id", nil},
		{"POST", "/user/keys", nil},
		{"PATCH", "/user/keys/:id", nil},
		{"DELETE", "/user/keys/:id", nil},
	}
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

	cases := []struct{
		DirectoryPath string
		StatusCode int
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

	cases := []struct{
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
	cases := []struct{
		Method string
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
