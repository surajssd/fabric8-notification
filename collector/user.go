package collector

import (
	"context"
	"fmt"

	"github.com/fabric8-services/fabric8-notification/auth/api"
	"github.com/goadesign/goa/uuid"
)

func NewUserResolver(c *api.Client) ReceiverResolver {
	return func(
		ctx context.Context, id string, revisionID uuid.UUID,
	) ([]Receiver, map[string]interface{}, error) {
		userID, err := uuid.FromString(id)
		if err != nil {
			return []Receiver{}, nil, fmt.Errorf("unable to lookup user based on id %v", id)
		}
		return User(ctx, c, userID)
	}
}

func User(ctx context.Context, c *api.Client, userID uuid.UUID) ([]Receiver, map[string]interface{}, error) {

	var values = map[string]interface{}{}
	var users []uuid.UUID

	users = append(users, userID)
	resolved, err := resolveAllUsers(ctx, c, SliceUniq(users), nil, true)
	if err != nil {
		return nil, nil, err
	}
	return resolved, values, nil
}
