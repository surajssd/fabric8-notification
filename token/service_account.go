package token

import (
	"context"
	"errors"

	"github.com/fabric8-services/fabric8-notification/auth/api"

	errs "github.com/fabric8-services/fabric8-wit/errors"
	"github.com/fabric8-services/fabric8-wit/goasupport"
	"github.com/fabric8-services/fabric8-wit/log"
)

type Fabric8ServiceAccountTokenClient struct {
	client        *api.Client
	accountID     string
	accountSecret string
}

func NewFabric8ServiceAccountTokenClient(client *api.Client, accountID string, accountSecret string) *Fabric8ServiceAccountTokenClient {
	return &Fabric8ServiceAccountTokenClient{
		client:        client,
		accountID:     accountID,
		accountSecret: accountSecret,
	}
}

type Fabric8ServiceAccountTokenService interface {
	Get(ctx context.Context) (string, error)
}

func (c *Fabric8ServiceAccountTokenClient) Get(ctx context.Context) (string, error) {
	tokenString, err := getServiceAccountToken(ctx, c.client, c.accountID, c.accountSecret)
	if err != nil {
		log.Error(ctx, map[string]interface{}{
			"err": err,
		}, "unable to get service account token")
		return "", err
	}

	if tokenString == nil {
		return "", errs.NewInternalError(ctx, errors.New("couldn't generate service account token"))
	}
	return tokenString.AccessToken, nil
}

func getServiceAccountToken(ctx context.Context, client *api.Client, serviceAccountID string, serviceAccountSecret string) (*api.ExternalToken, error) {
	payload := api.TokenExchange{
		ClientID:     serviceAccountID,
		ClientSecret: &serviceAccountSecret,
		GrantType:    "client_credentials",
	}
	resp, err := client.ExchangeToken(goasupport.ForwardContextRequestID(ctx), api.ExchangeTokenPath(), &payload, "application/x-www-form-urlencoded")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if err != nil {
		log.Error(ctx, map[string]interface{}{
			"err": err,
		}, "failed to get service account token from auth")
		return nil, err
	}

	return client.DecodeExternalToken(resp)

}
