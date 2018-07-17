package collector

import (
	"context"

	"github.com/fabric8-services/fabric8-notification/configuration"

	"github.com/goadesign/goa/uuid"
)

func ConfiguredVars(config *configuration.Data, resolver ReceiverResolver) ReceiverResolver {
	return func(
		ctx context.Context, id string, revisionID uuid.UUID,
	) ([]Receiver, map[string]interface{}, error) {
		r, v, err := resolver(ctx, id, revisionID)
		if err != nil {
			return r, v, err
		}
		if v == nil {
			v = map[string]interface{}{}
		}
		v["webURL"] = config.GetWebURL()
		return r, v, err
	}
}
