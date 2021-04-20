package milton

import (
	"context"
	"sync"

	"github.com/slack-go/slack"
	"github.com/slack-go/slack/slackevents"
	"github.com/slack-go/slack/socketmode"
	"github.com/vj396/milton/src/backend"
	slackPkg "github.com/vj396/milton/src/slack"
	"github.com/vj396/milton/src/types"
	"go.uber.org/zap"

	_ "github.com/vj396/milton/src/plugins/interrupts"
	_ "github.com/vj396/milton/src/plugins/opsgenie"
)

type bot struct {
	backend types.Backend
	slack   *socketmode.Client

	id     string
	name   string
	logger *zap.Logger
}

func Start(done chan struct{}, logger *zap.Logger, conf *types.Config, modelsDir string) {
	var err error
	b := new(bot)
	b.logger = logger
	b.backend, err = backend.Connect(logger, conf.Backend.Type, conf.Backend, modelsDir)
	if err != nil {
		b.logger.Fatal(err.Error())
	}
	b.slack, err = slackPkg.Client(conf.Slack)
	if err != nil {
		b.logger.Fatal(err.Error())
	}
	ctx, cancel := context.WithCancel(context.TODO())
	var wg sync.WaitGroup
	wg.Add(1)
	go b.run(ctx, &wg)
	<-done
	cancel()
}

func (b *bot) run(ctx context.Context, wg *sync.WaitGroup) {
	defer wg.Done()
	go func() {
		if err := b.slack.Run(); err != nil {
			b.logger.Error("could not run bot in socket-mode", zap.Error(err))
		}
	}()
	for event := range b.slack.Events {
		select {
		case <-ctx.Done():
			return
		default:
		}
		switch event.Type {
		case socketmode.EventTypeConnecting:
			b.logger.Debug("connecting to slack in socketmode")
		case socketmode.EventTypeConnectionError:
			b.logger.Error("connection failed. retrying later...")
		case socketmode.EventTypeConnected:
			b.logger.Debug("connected to slack in socketmode", zap.String("id", b.id), zap.String("name", b.name))
		case socketmode.EventTypeEventsAPI:
			eventsAPIEvent, ok := event.Data.(slackevents.EventsAPIEvent)
			if !ok {
				b.logger.Debug("ignoring", zap.Any("event", event))
				continue
			}
			b.slack.Ack(*event.Request)
			switch eventsAPIEvent.Type {
			case slackevents.CallbackEvent:
				switch ev := eventsAPIEvent.InnerEvent.Data.(type) {
				case *slackevents.AppMentionEvent:
					_, _, err := slackPkg.ApiClient().PostMessage(ev.Channel, slack.MsgOptionText("Yes, hello.", false))
					if err != nil {
						b.logger.Error("failed to post message", zap.Error(err))
					}
				case *slackevents.MemberJoinedChannelEvent:
					b.logger.Info("member joined channel", zap.String("user", ev.User), zap.String("channel", ev.Channel))
				}
			default:
				b.logger.Debug("unsupported events api event received", zap.String("event", eventsAPIEvent.InnerEvent.Type))
			}
		default:
			b.logger.Info("", zap.Any("event", event))
		}
	}
}
