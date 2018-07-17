package collector

import (
	"context"
	"testing"

	"github.com/fabric8-services/fabric8-notification/configuration"

	"github.com/goadesign/goa/uuid"
	"github.com/stretchr/testify/assert"
)

func TestConfigureVarsSetWebURL(t *testing.T) {
	config, err := configuration.NewData()
	if err != nil {
		assert.NoError(t, err)
	}
	u, v, err := ConfiguredVars(config, EmptyResolver)(context.Background(), "id", uuid.NewV4())
	assert.NoError(t, err)
	assert.Len(t, u, 0)
	assert.Len(t, v, 1)
	assert.NotEmpty(t, v["webURL"])
}

func TestConfigureVarsNilResolverVars(t *testing.T) {
	config, err := configuration.NewData()
	if err != nil {
		assert.NoError(t, err)
	}
	u, v, err := ConfiguredVars(config, NilVarResolver)(context.Background(), "id", uuid.NewV4())
	assert.NoError(t, err)
	assert.Len(t, u, 0)
	assert.Len(t, v, 1)
	assert.NotEmpty(t, v["webURL"])
}

func EmptyResolver(context.Context, string, uuid.UUID) (users []Receiver, templateValues map[string]interface{}, err error) {
	return []Receiver{}, map[string]interface{}{}, nil
}

func NilVarResolver(context.Context, string, uuid.UUID) (users []Receiver, templateValues map[string]interface{}, err error) {
	return []Receiver{}, nil, nil
}
