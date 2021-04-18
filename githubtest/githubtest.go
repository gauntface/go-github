package githubtest

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/google/go-github/v35/github"
)

func New(t *testing.T, opts ...Option) (*github.Client, func()) {
	t.Helper()

	// Sort through inputs from test to dictate behavior of the
	o := &options{}
	for _, oo := range opts {
		err := oo(o)
		if err != nil {
			t.Fatalf("Failed to process githubtest.New() options: %v", err)
		}
	}

	ts := httptest.NewServer(http.HandlerFunc(fakeResponses(o)))
	gh := github.NewClient(ts.Client())
	tu := ts.URL
	if !strings.HasSuffix(tu, "/") {
		tu = tu + "/"
	}
	u, err := url.Parse(tu)
	if err != nil {
		t.Fatalf("Failed to parse test server URL: %v", err)
	}
	gh.BaseURL = u
	return gh, ts.Close
}

func fakeResponses(opts *options) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/user/orgs":
			b, err := json.Marshal(opts.Organizations[""])
			if err != nil {
				fmt.Printf("Failed to marshal response for %v\n", r.URL)
				return
			}
			io.WriteString(w, string(b))
		default:
			fmt.Printf("Unhandled network request: %v\n", r.URL)
		}
	}
}

type options struct {
	Organizations map[string][]*github.Organization
}

type Option func(o *options) error

func WithOrganizations(orgs map[string][]*github.Organization) Option {
	return func(o *options) error {
		o.Organizations = orgs
		return nil
	}
}
