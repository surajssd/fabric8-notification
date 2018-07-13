package auth

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"strconv"

	"github.com/fabric8-services/fabric8-notification/auth/api"

	"github.com/fabric8-services/fabric8-wit/goasupport"
	"github.com/goadesign/goa/uuid"
)

type CollaboratorCollector interface {
	GetSpaceCollaborators(ctx context.Context, client *api.Client, spaceID uuid.UUID) (*api.UserList, error)
}

type AuthCollector struct {
}

func (c *AuthCollector) GetSpaceCollaborators(ctx context.Context, client *api.Client, spaceID uuid.UUID) (*api.UserList, error) {
	pageLimit := 100
	pageOffset := "0"
	resp, err := listCollaborators(goasupport.ForwardContextRequestID(ctx), client, api.ListCollaboratorsPath(spaceID), &pageLimit, &pageOffset, nil, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to make request to get list of collaborators")
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("non %v status code for %v, returned %v", http.StatusOK, "GET collaborators", resp.StatusCode)
	}
	return client.DecodeUserList(resp)
}

// list collaborators for the given space ID.
func listCollaborators(ctx context.Context, client *api.Client, path string, pageLimit *int, pageOffset *string, ifModifiedSince *string, ifNoneMatch *string) (*http.Response, error) {
	req, err := newListCollaboratorsRequest(ctx, client, path, pageLimit, pageOffset, ifModifiedSince, ifNoneMatch)
	if err != nil {
		return nil, err
	}
	return client.Do(ctx, req)
}

// newListCollaboratorsRequest create the request corresponding to the list action endpoint of the collaborators resource.
func newListCollaboratorsRequest(ctx context.Context, client *api.Client, path string, pageLimit *int, pageOffset *string, ifModifiedSince *string, ifNoneMatch *string) (*http.Request, error) {
	scheme := client.Scheme
	if scheme == "" {
		scheme = "http"
	}
	u := url.URL{Host: client.Host, Scheme: scheme, Path: path}
	values := u.Query()
	if pageLimit != nil {
		tmp3 := strconv.Itoa(*pageLimit)
		values.Set("page[limit]", tmp3)
	}
	if pageOffset != nil {
		values.Set("page[offset]", *pageOffset)
	}
	u.RawQuery = values.Encode()
	req, err := http.NewRequest("GET", u.String(), nil)
	if err != nil {
		return nil, err
	}
	header := req.Header
	if ifModifiedSince != nil {

		header.Set("If-Modified-Since", *ifModifiedSince)
	}
	if ifNoneMatch != nil {

		header.Set("If-None-Match", *ifNoneMatch)
	}

	// This wasn't present in the auto-generated client for GET /api/spaces/ID/collaborators/ID
	// because we don't set "a.Security("jwt")" in auth's design/account.go
	// for the `List Collaborators` action.
	if client.JWTSigner != nil {
		client.JWTSigner.Sign(req)
	}

	return req, nil
}
