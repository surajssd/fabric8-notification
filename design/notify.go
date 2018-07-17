package design

import (
	d "github.com/goadesign/goa/design"
	a "github.com/goadesign/goa/design/apidsl"
)

var notification = a.Type("Notification", func() {
	a.Description(`JSONAPI for the Notification object. See also http://jsonapi.org/format/#document-resource-object`)
	a.Attribute("type", d.String, func() {
		a.Enum("notifications")
	})
	a.Attribute("id", d.UUID, "ID of notification", func() {
		a.Example("40bbdd3d-8b5d-4fd6-ac90-7236b669af04")
	})
	a.Attribute("attributes", notificationAttributes)
	//a.Attribute("links", genericLinks)
	a.Required("type", "attributes")
})

var notificationAttributes = a.Type("NotificationAttributes", func() {
	a.Description(`JSONAPI store for all the "attributes" of a Notification. See also see http://jsonapi.org/format/#document-resource-object-attributes`)
	a.Attribute("type", d.String, "The notification type", func() {
		a.Example("workitem.create")
	})
	a.Attribute("id", d.String, "ID of the main resource that was created/changed", func() {
		a.Example("8bccc228-bba7-43ad-b077-15fbb9148f7f")
	})
	a.Attribute("custom", a.HashOf(d.String, d.Any), "custom information to be passed as a json", func() {
		a.Example(map[string]interface{}{"verification_url": "https://auth.openshift.io/api/user/verifyEmail?code=bb9148f7"})
	})
	a.Attribute("revisionID", d.UUID, "ID of revision that triggered this notification", func() {
		a.Example("1646e178-b90e-48b4-85e4-204d1c4c416f")
	})
	a.Required("type", "id")
})

var notificationSingle = JSONSingle(
	"notification", "Holds a single notification",
	notification,
	nil)

var _ = a.Resource("notify", func() {
	a.BasePath("/notify")
	a.Action("send", func() {
		a.Security("jwt")
		a.Routing(
			a.POST(""),
		)
		a.Payload(notificationSingle)
		a.Description("Register a new notification.")
		a.Response(d.Accepted)
		a.Response(d.BadRequest, JSONAPIErrors)
		a.Response(d.InternalServerError, JSONAPIErrors)
		a.Response(d.Unauthorized, JSONAPIErrors)
	})
})
