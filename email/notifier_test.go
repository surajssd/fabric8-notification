package email

import (
	"context"
	"testing"

	"github.com/fabric8-services/fabric8-notification/collector"
	"github.com/fabric8-services/fabric8-notification/template"
	"github.com/fabric8-services/fabric8-notification/types"
)

type TestSender struct {
	callback chan bool
}

func (t *TestSender) Send(ctx context.Context, subject string, body string, headers map[string]string, receivers []collector.Receiver) {
	t.callback <- true
}

func TestAsyncWorkerNotifier(t *testing.T) {

	resolver := func(ctx context.Context, id string) ([]collector.Receiver, map[string]interface{}, error) {
		return []collector.Receiver{}, nil, nil
	}

	callback := make(chan bool)

	sender := &TestSender{callback: callback}
	notifier := NewAsyncWorkerNotifier(sender, 1)

	notifier.Send(context.Background(), Notification{ID: "TEST", CustomAttributes: map[string]interface{}{}, Type: string(types.WorkitemCreate), Resolver: resolver, Template: template.Template{}})

	<-sender.callback
}
