package wit

import (
	"context"
	"fmt"
	"net/http"
	"net/url"

	"github.com/fabric8-services/fabric8-notification/wit/api"

	"github.com/fabric8-services/fabric8-wit/goasupport"
	goaclient "github.com/goadesign/goa/client"
	"github.com/goadesign/goa/uuid"
	"github.com/gregjones/httpcache"
)

func NewCachedClient(hostURL string) (*api.Client, error) {

	u, err := url.Parse(hostURL)
	if err != nil {
		return nil, err
	}

	tp := httpcache.NewMemoryCacheTransport()
	client := http.Client{Transport: tp}

	c := api.New(goaclient.HTTPClientDoer(&client))
	c.Host = u.Host
	c.Scheme = u.Scheme
	return c, nil
}

func GetWorkItem(ctx context.Context, client *api.Client, wiID uuid.UUID) (*api.WorkItemSingle, error) {
	resp, err := client.ShowWorkitem(goasupport.ForwardContextRequestID(ctx), api.ShowWorkitemPath(wiID), nil, nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("non %v status code for %v, returned %v", http.StatusOK, "GET workitem", resp.StatusCode)
	}

	if err != nil {
		return nil, err
	}
	return client.DecodeWorkItemSingle(resp)
}

func GetArea(ctx context.Context, client *api.Client, areaID uuid.UUID) (*api.AreaSingle, error) {
	resp, err := client.ShowArea(goasupport.ForwardContextRequestID(ctx), api.ShowAreaPath(areaID.String()), nil, nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("non %v status code for %v, returned %v", http.StatusOK, "GET space", resp.StatusCode)
	}

	if err != nil {
		return nil, err
	}
	return client.DecodeAreaSingle(resp)
}

func GetWorkItemType(ctx context.Context, client *api.Client, witID uuid.UUID) (*api.WorkItemTypeSingle, error) {
	resp, err := client.ShowWorkitemtype(goasupport.ForwardContextRequestID(ctx), api.ShowWorkitemtypePath(witID), nil, nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("non %v status code for %v, returned %v", http.StatusOK, "GET workitemtype", resp.StatusCode)
	}

	if err != nil {
		return nil, err
	}
	return client.DecodeWorkItemTypeSingle(resp)
}

func GetComment(ctx context.Context, client *api.Client, cID uuid.UUID) (*api.CommentSingle, error) {
	resp, err := client.ShowComments(goasupport.ForwardContextRequestID(ctx), api.ShowCommentsPath(cID), nil, nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("non %v status code for %v, returned %v", http.StatusOK, "GET comment", resp.StatusCode)
	}

	if err != nil {
		return nil, err
	}
	return client.DecodeCommentSingle(resp)
}

func GetComments(ctx context.Context, client *api.Client, wiID uuid.UUID) (*api.CommentList, error) {
	pageLimit := 100
	pageOffset := "0"

	resp, err := client.ListWorkItemComments(goasupport.ForwardContextRequestID(ctx), api.ListWorkItemCommentsPath(wiID), &pageLimit, &pageOffset, nil, nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("non %v status code for %v, returned %v", http.StatusOK, "GET comments", resp.StatusCode)
	}

	if err != nil {
		return nil, err
	}
	return client.DecodeCommentList(resp)
}

func GetSpace(ctx context.Context, client *api.Client, spaceID uuid.UUID) (*api.SpaceSingle, error) {
	resp, err := client.ShowSpace(goasupport.ForwardContextRequestID(ctx), api.ShowSpacePath(spaceID), nil, nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("non %v status code for %v, returned %v", http.StatusOK, "GET space", resp.StatusCode)
	}

	if err != nil {
		return nil, err
	}
	return client.DecodeSpaceSingle(resp)
}

func GetSpaces(ctx context.Context, client *api.Client, spaceIDs []uuid.UUID) ([]*api.SpaceSingle, error) {
	var spaces []*api.SpaceSingle
	for _, spaceID := range spaceIDs {
		space, err := GetSpace(ctx, client, spaceID)
		if err != nil {
			return nil, err
		}
		spaces = append(spaces, space)
	}
	return spaces, nil
}

func GetCodebases(ctx context.Context, client *api.Client, url string) (*api.CodebaseList, error) {
	resp, err := client.CodebasesSearch(context.Background(), api.CodebasesSearchPath(), url, nil, nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("non %v status code for %v, returned %v", http.StatusOK, "GET space", resp.StatusCode)
	}
	if err != nil {
		return nil, err
	}
	return client.DecodeCodebaseList(resp)
}
