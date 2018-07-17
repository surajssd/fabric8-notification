package email

import (
	"context"

	"github.com/fabric8-services/fabric8-notification/collector"
	"github.com/fabric8-services/fabric8-notification/template"

	"github.com/fabric8-services/fabric8-wit/log"
	"github.com/goadesign/goa/uuid"
)

type contextualNotification struct {
	context      context.Context
	notification Notification
}

type Notification struct {
	Type             string
	ID               string
	RevisionID       uuid.UUID
	CustomAttributes map[string]interface{}
	Resolver         collector.ReceiverResolver
	Template         template.Template
}

type Notifier interface {
	Send(context.Context, Notification)
}

func NewAsyncWorkerNotifier(sender Sender, concurrency int) Notifier {
	notifier := &AsyncWorkerNotifier{
		Sender:      sender,
		concurrency: concurrency,
		tasks:       make(chan contextualNotification)}
	notifier.start()
	return notifier
}

type CallbackNotifier struct {
	Callback func(ctx context.Context, notification Notification)
}

func (c *CallbackNotifier) Send(ctx context.Context, notification Notification) {
	if c.Callback != nil {
		c.Callback(ctx, notification)
	}
}

type AsyncWorkerNotifier struct {
	Sender      Sender
	concurrency int
	tasks       chan contextualNotification
}

func (a *AsyncWorkerNotifier) start() {
	for i := 0; i < a.concurrency; i++ {
		go a.work()
	}
}

func (a *AsyncWorkerNotifier) work() {
	for task := range a.tasks {
		log.Debug(task.context, map[string]interface{}{
			"type": task.notification.ID,
			"id":   task.notification.Type,
		}, "working on new notificaiton send of notification")

		a.do(task)

	}
}

func (a *AsyncWorkerNotifier) do(cn contextualNotification) {
	ctx, notification := cn.context, cn.notification

	receivers, vars, err := notification.Resolver(ctx, notification.ID, notification.RevisionID)
	if err != nil {
		log.Error(ctx, map[string]interface{}{
			"type": notification.ID,
			"id":   notification.Type,
			"err":  err,
		}, "failed to resolve receivers")

		return
	}

	if vars == nil {
		vars = map[string]interface{}{}
	}

	vars["custom"] = notification.CustomAttributes

	subject, body, headers, err := notification.Template.Render(vars)
	if err != nil {
		log.Error(ctx, map[string]interface{}{
			"type": notification.ID,
			"id":   notification.Type,
			"err":  err,
		}, "failed to render template")

		return
	}
	a.Sender.Send(ctx, subject, body, headers, receivers)
}

func (a *AsyncWorkerNotifier) Send(ctx context.Context, notification Notification) {
	a.tasks <- contextualNotification{context: ctx, notification: notification}

	log.Debug(ctx, map[string]interface{}{
		"type": notification.ID,
		"id":   notification.Type,
	}, "scheduled send of notification")
}
