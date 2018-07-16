package controller

import (
	"context"
	"fmt"
	"testing"

	jwt "github.com/dgrijalva/jwt-go"
	goajwt "github.com/goadesign/goa/middleware/security/jwt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/fabric8-services/fabric8-notification/app"
	"github.com/fabric8-services/fabric8-notification/app/test"
	"github.com/fabric8-services/fabric8-notification/auth"
	"github.com/fabric8-services/fabric8-notification/collector"
	"github.com/fabric8-services/fabric8-notification/configuration"
	"github.com/fabric8-services/fabric8-notification/email"
	"github.com/fabric8-services/fabric8-notification/template"
	"github.com/fabric8-services/fabric8-notification/types"
	"github.com/fabric8-services/fabric8-notification/validator"
	"github.com/goadesign/goa"
)

func TestNotifySendUnknownType(t *testing.T) {

	resolvers := collector.NewRegistry()
	typeRegistry := &template.AssetRegistry{}
	notifier := &email.CallbackNotifier{Callback: func(ctx context.Context, notification email.Notification) { fmt.Println(notification) }}
	ctrl := NewNotifyController(goa.New("send-test"), resolvers, typeRegistry, notifier)

	payload := &app.SendNotifyPayload{
		Data: &app.Notification{
			Attributes: &app.NotificationAttributes{
				ID:   "13131",
				Type: "unknown.create",
			},
			Type: "notifications",
		},
	}

	test.SendNotifyBadRequest(t, nil, nil, ctrl, payload)
}

func TestNotifySendWithCustomParam(t *testing.T) {

	config, _ := configuration.GetData()

	resolvers := collector.NewRegistry()
	typeRegistry := &template.AssetRegistry{}
	authClient, _ := auth.NewCachedClient(config.GetWITURL())

	resolvers.Register(types.UserEmailUpdate, collector.ConfiguredVars(config, collector.NewUserResolver(authClient)), validator.ValidateUser)
	notifier := &email.CallbackNotifier{
		Callback: func(ctx context.Context, notification email.Notification) {
			require.NotNil(t, notification.CustomAttributes)
			assert.Equal(t, notification.CustomAttributes["verifyURL"], "https://someurl.openshift.io")
		}}

	ctrl := NewNotifyController(goa.New("send-test"), resolvers, typeRegistry, notifier)

	payload := &app.SendNotifyPayload{
		Data: &app.Notification{
			Attributes: &app.NotificationAttributes{
				ID:   "13132",
				Type: string(types.UserEmailUpdate),
				Custom: map[string]interface{}{
					"verifyURL": "https://someurl.openshift.io",
				},
			},
			Type: "notifications",
		},
	}

	claims := jwt.MapClaims{}
	claims["service_accountname"] = "fabric8-auth"
	ctx := goajwt.WithJWT(context.Background(), jwt.NewWithClaims(jwt.SigningMethodRS512, claims))

	test.SendNotifyAccepted(t, ctx, nil, ctrl, payload)
}

func TestNotifySendWithoutustomParamBadRequest(t *testing.T) {

	config, _ := configuration.GetData()

	resolvers := collector.NewRegistry()
	typeRegistry := &template.AssetRegistry{}
	authClient, _ := auth.NewCachedClient(config.GetWITURL())

	resolvers.Register(types.UserEmailUpdate, collector.ConfiguredVars(config, collector.NewUserResolver(authClient)), validator.ValidateUser)
	notifier := &email.CallbackNotifier{
		Callback: func(ctx context.Context, notification email.Notification) {
			require.NotNil(t, notification.CustomAttributes)
			assert.Equal(t, notification.CustomAttributes["verifyURL"], "https://someurl.openshift.io")
		}}

	ctrl := NewNotifyController(goa.New("send-test"), resolvers, typeRegistry, notifier)

	payload := &app.SendNotifyPayload{
		Data: &app.Notification{
			Attributes: &app.NotificationAttributes{
				ID:   "13132",
				Type: string(types.UserEmailUpdate),
			},
			Type: "notifications",
		},
	}

	claims := jwt.MapClaims{}
	claims["service_accountname"] = "fabric8-auth"
	ctx := goajwt.WithJWT(context.Background(), jwt.NewWithClaims(jwt.SigningMethodRS512, claims))

	test.SendNotifyBadRequest(t, ctx, nil, ctrl, payload)
}

func TestNotifySendWithCustomParamBadRequest(t *testing.T) {

	config, _ := configuration.GetData()

	resolvers := collector.NewRegistry()
	typeRegistry := &template.AssetRegistry{}
	authClient, _ := auth.NewCachedClient(config.GetWITURL())

	resolvers.Register(types.UserEmailUpdate, collector.ConfiguredVars(config, collector.NewUserResolver(authClient)), validator.ValidateUser)
	notifier := &email.CallbackNotifier{
		Callback: func(ctx context.Context, notification email.Notification) {
			require.NotNil(t, notification.CustomAttributes)
			assert.Equal(t, notification.CustomAttributes["verifyURL"], "https://someurl.openshift.io")
		}}

	ctrl := NewNotifyController(goa.New("send-test"), resolvers, typeRegistry, notifier)

	payload := &app.SendNotifyPayload{
		Data: &app.Notification{
			Attributes: &app.NotificationAttributes{
				ID:   "13132",
				Type: string(types.UserEmailUpdate),
				Custom: map[string]interface{}{
					"somthing_else": "https://someurl.openshift.io",
				},
			},
			Type: "notifications",
		},
	}

	claims := jwt.MapClaims{}
	claims["service_accountname"] = "fabric8-auth"
	ctx := goajwt.WithJWT(context.Background(), jwt.NewWithClaims(jwt.SigningMethodRS512, claims))

	test.SendNotifyBadRequest(t, ctx, nil, ctrl, payload)
}

func TestNotifySendWithoutCustomParamSuccess(t *testing.T) {

	config, _ := configuration.GetData()

	resolvers := collector.NewRegistry()
	typeRegistry := &template.AssetRegistry{}
	authClient, _ := auth.NewCachedClient(config.GetWITURL())

	resolvers.Register(types.WorkitemUpdate, collector.ConfiguredVars(config, collector.NewUserResolver(authClient)), nil)
	notifier := &email.CallbackNotifier{
		Callback: func(ctx context.Context, notification email.Notification) {
		}}

	ctrl := NewNotifyController(goa.New("send-test"), resolvers, typeRegistry, notifier)

	payload := &app.SendNotifyPayload{
		Data: &app.Notification{
			Attributes: &app.NotificationAttributes{
				ID:   "13132",
				Type: string(types.WorkitemUpdate),
			},
			Type: "notifications",
		},
	}

	test.SendNotifyAccepted(t, nil, nil, ctrl, payload)
}

func TestNotifySendUnauthorized(t *testing.T) {
	notifier := &email.CallbackNotifier{Callback: func(ctx context.Context, notification email.Notification) { /* blank */ }}
	ctrl := NewNotifyController(goa.New("send-test"), collector.NewRegistry(), &template.AssetRegistry{}, notifier)

	payload := &app.SendNotifyPayload{
		Data: &app.Notification{
			Attributes: &app.NotificationAttributes{
				ID:   "git@github.com:testrepo/testproject1.git",
				Type: string(types.AnalyticsNotifyCVE),
			},
			Type: "notifications",
		},
	}

	claims := jwt.MapClaims{}
	claims["service_accountname"] = "test-service"
	ctx := goajwt.WithJWT(context.Background(), jwt.NewWithClaims(jwt.SigningMethodRS512, claims))

	test.SendNotifyUnauthorized(t, ctx, nil, ctrl, payload)
}

func TestValidateNotifier(t *testing.T) {

	t.Run("service_allowed", func(t *testing.T) {
		claims := jwt.MapClaims{}
		claims["service_accountname"] = "fabric8-gemini-server"
		ctx := goajwt.WithJWT(context.Background(), jwt.NewWithClaims(jwt.SigningMethodRS512, claims))
		assert.True(t, validateNotifier(ctx, types.AnalyticsNotifyCVE.Notifiers()))
	})

	t.Run("service_not_allowed", func(t *testing.T) {
		claims := jwt.MapClaims{}
		claims["service_accountname"] = "test-service"
		ctx := goajwt.WithJWT(context.Background(), jwt.NewWithClaims(jwt.SigningMethodRS512, claims))
		assert.False(t, validateNotifier(ctx, types.AnalyticsNotifyCVE.Notifiers()))
	})

	t.Run("service_missing", func(t *testing.T) {
		claims := jwt.MapClaims{}
		ctx := goajwt.WithJWT(context.Background(), jwt.NewWithClaims(jwt.SigningMethodRS512, claims))
		assert.False(t, validateNotifier(ctx, types.AnalyticsNotifyCVE.Notifiers()))
	})

	t.Run("no_service_restriction", func(t *testing.T) {
		claims := jwt.MapClaims{}
		ctx := goajwt.WithJWT(context.Background(), jwt.NewWithClaims(jwt.SigningMethodRS512, claims))
		assert.True(t, validateNotifier(ctx, types.WorkitemCreate.Notifiers()))
	})

}
