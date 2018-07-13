package auth

import (
	"context"
	"fmt"
	"net/http"
	"net/url"

	"github.com/fabric8-services/fabric8-notification/auth/api"

	"github.com/fabric8-services/fabric8-wit/goasupport"
	"github.com/goadesign/goa/uuid"
)

/*
	Took out the auth api client code relevant to User API in order to add a custom
	client.JWTSigner.Sign(req) before the http request is made.

	This wasn't present in the auto-generated client for GET /api/users/ID
	because we don't set "a.Security("jwt")" in auth's design/account.go
	for the `Show Users` action.
*/

func GetUser(ctx context.Context, client *api.Client, uID uuid.UUID) (*api.User, error) {
	resp, err := showUsers(goasupport.ForwardContextRequestID(ctx), client, api.ShowUsersPath(uID.String()), nil, nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("non %v status code for %v, returned %v", http.StatusOK, "GET user", resp.StatusCode)
	}

	if err != nil {
		return nil, err
	}
	return client.DecodeUser(resp)
}

// retrieve user for the given ID.
func showUsers(ctx context.Context, client *api.Client, path string, ifModifiedSince *string, ifNoneMatch *string) (*http.Response, error) {
	req, err := newShowUsersRequest(ctx, client, path, ifModifiedSince, ifNoneMatch)
	if err != nil {
		return nil, err
	}
	return client.Do(ctx, req)
}

// newShowUsersRequest create the request corresponding to the show action endpoint of the users resource.
func newShowUsersRequest(ctx context.Context, client *api.Client, path string, ifModifiedSince *string, ifNoneMatch *string) (*http.Request, error) {
	scheme := client.Scheme
	if scheme == "" {
		scheme = "http"
	}
	u := url.URL{Host: client.Host, Scheme: scheme, Path: path}
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

	// This wasn't present in the auto-generated client for GET /api/users/ID
	// because we don't set "a.Security("jwt")" in auth's design/account.go
	// for the `Show Users` action.
	if client.JWTSigner != nil {
		client.JWTSigner.Sign(req)
	}
	return req, nil
}
