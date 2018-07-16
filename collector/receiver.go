package collector

import (
	"context"

	"github.com/fabric8-services/fabric8-notification/types"
)

type ReceiverResolver func(context.Context, string) (users []Receiver, templateValues map[string]interface{}, err error)
type ParamValidator func(context.Context, map[string]interface{}) error

type Receiver struct {
	FullName string
	EMail    string
}

type Registry interface {
	Register(types.NotificationType, ReceiverResolver, ParamValidator)
	Get(types.NotificationType) (ReceiverResolver, bool)
	Validator(types.NotificationType) (ParamValidator, bool)
	Notifiers(types.NotificationType) []string
}

type localRegistry struct {
	reg        map[types.NotificationType]ReceiverResolver
	validators map[types.NotificationType]ParamValidator
}

func NewRegistry() Registry {
	return &localRegistry{
		reg:        map[types.NotificationType]ReceiverResolver{},
		validators: map[types.NotificationType]ParamValidator{},
	}
}

func (r *localRegistry) Get(nType types.NotificationType) (ReceiverResolver, bool) {
	res, b := r.reg[nType]
	return res, b
}

func (r *localRegistry) Validator(nType types.NotificationType) (ParamValidator, bool) {
	v, found := r.validators[nType]
	return v, found
}

func (r *localRegistry) Notifiers(nType types.NotificationType) []string {
	return nType.Notifiers()
}

func (r *localRegistry) Register(nType types.NotificationType, resolver ReceiverResolver, validator ParamValidator) {
	if resolver != nil {
		r.reg[nType] = resolver
	}
	if validator != nil {
		r.validators[nType] = validator
	}
}
