// Copyright 2020 The go-github AUTHORS. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package githubtest

import (
	"context"
	"errors"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-github/v35/github"
)

func TestGithubTest_ExampleUsage(t *testing.T) {
	want := []string{"org-1", "org-2"}
	var wantError error = nil

	ghFake, close := New(t, WithOrganizations(map[string][]*github.Organization{
		"": []*github.Organization{
			{Name: github.String("org-1")},
			{Name: github.String("org-2")},
		},
	}))
	defer close()

	c := myClient{
		gh: ghFake,
	}

	got, err := c.AllOrgNames(context.Background())
	if !errors.Is(err, wantError) {
		t.Fatalf("Unexpected error; got %v, want %v", err, wantError)
	}

	if diff := cmp.Diff(got, want); diff != "" {
		t.Fatalf("Unexpected result; diff %v", diff)
	}
}

type myClient struct {
	gh *github.Client
}

func (c *myClient) AllOrgNames(ctx context.Context) ([]string, error) {
	orgs, _, err := c.gh.Organizations.List(ctx, "", nil)
	if err != nil {
		return nil, err
	}

	on := []string{}
	for _, o := range orgs {
		on = append(on, o.GetName())
	}
	return on, nil
}
