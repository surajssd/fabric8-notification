package collector

import (
	"context"

	authapi "github.com/fabric8-services/fabric8-notification/auth/api"
	"github.com/fabric8-services/fabric8-notification/wit"
	"github.com/fabric8-services/fabric8-notification/wit/api"
	"github.com/goadesign/goa/uuid"
)

func NewCVEResolver(authClient *authapi.Client, witClient *api.Client) ReceiverResolver {
	return func(
		ctx context.Context, url string, revisionID uuid.UUID,
	) ([]Receiver, map[string]interface{}, error) {
		codebases, err := wit.GetCodebases(ctx, witClient, url)
		if err != nil {
			return nil, nil, err
		}
		spaceIDs := collectCodebasesSpaces(codebases)
		spaces, err := wit.GetSpaces(ctx, witClient, spaceIDs)
		if err != nil {
			return nil, nil, err
		}

		users := collectSpacesUsers(spaces)
		resolved, err := resolveAllUsers(ctx, authClient, SliceUniq(users), []*authapi.UserData{}, false)
		if err != nil {
			return nil, nil, err
		}
		return resolved, map[string]interface{}{}, nil
	}
}

func collectCodebasesSpaces(codebases *api.CodebaseList) []uuid.UUID {
	var spaces []uuid.UUID
	for _, cb := range codebases.Data {
		if cb.Attributes.CveScan != nil && *cb.Attributes.CveScan {
			space, err := uuid.FromString(*cb.Relationships.Space.Data.ID)
			if err == nil {
				spaces = append(spaces, space)
			}
		}
	}
	return spaces
}

func collectSpacesUsers(spaces []*api.SpaceSingle) []uuid.UUID {
	var users []uuid.UUID
	for _, space := range spaces {
		users = append(users, *space.Data.Relationships.OwnedBy.Data.ID)
	}
	return users
}
