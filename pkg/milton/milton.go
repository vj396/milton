package milton

import (
	"context"
	"encoding/json"
	"fmt"
	"regexp"
	"strings"
	"sync"

	"github.com/slack-go/slack"
	"github.com/slack-go/slack/slackevents"
	"github.com/slack-go/slack/socketmode"
	"github.com/vj396/milton/src/backend"
	slackPkg "github.com/vj396/milton/src/slack"
	"github.com/vj396/milton/src/types"
	"go.uber.org/zap"

	_ "github.com/vj396/milton/src/plugins/help"
	_ "github.com/vj396/milton/src/plugins/interrupts"
	_ "github.com/vj396/milton/src/plugins/opsgenie"
)

const (
	defaultProcessQueueSize = 100
)

type bot struct {
	backend types.Backend
	slack   *socketmode.Client

	cmdRegex      *regexp.Regexp
	processQueue  chan *types.MessageMetadata
	responseQueue chan *types.MessageMetadata

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
	cmds := []string{}
	for cmd := range slackPkg.GetRegistry() {
		cmds = append(cmds, cmd)
	}
	b.cmdRegex = regexp.MustCompile(fmt.Sprintf("^(%s)", strings.Join(cmds, "|")))
	b.processQueue = make(chan *types.MessageMetadata, defaultProcessQueueSize)
	b.responseQueue = make(chan *types.MessageMetadata, defaultProcessQueueSize)
	defer func() {
		close(b.processQueue)
		close(b.responseQueue)
		b.backend.Close()
	}()
	ctx, cancel := context.WithCancel(context.TODO())
	var wg sync.WaitGroup
	wg.Add(1)
	go b.processQueueChannel(ctx, &wg)
	wg.Add(1)
	go b.responseQueueChannel(ctx, &wg)
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
		b.logger.Debug("event", zap.Any("dump", event))
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
			b.logger.Debug("connected to slack in socketmode")
		case socketmode.EventTypeErrorBadMessage:
			event.Request = new(socketmode.Request)
			m := make(map[string]interface{})
			json.Unmarshal(event.Data.(*socketmode.ErrorBadMessage).Message, &m)
			event.Request.EnvelopeID = m["envelope_id"].(string)
			b.slack.Ack(*event.Request)
		case socketmode.EventTypeEventsAPI:
			b.slack.Ack(*event.Request)
			eventsAPIEvent, ok := event.Data.(slackevents.EventsAPIEvent)
			if !ok {
				b.logger.Debug("ignoring", zap.Any("event", event))
				continue
			}
			switch eventsAPIEvent.Type {
			case slackevents.CallbackEvent:
				switch ev := eventsAPIEvent.InnerEvent.Data.(type) {
				case *slackevents.MessageEvent:
					if ev.BotID != "" {
						continue
					}
					if !b.cmdRegex.MatchString(ev.Text) {
						continue
					}
					ts := ev.TimeStamp
					if ev.ThreadTimeStamp != "" {
						ts = ev.ThreadTimeStamp
					}
					msg := types.MessageMetadata{ChannelID: ev.Channel, UserID: ev.User, Message: ev.Text, Timestamp: ts}
					b.processQueue <- &msg
				}
			}
		default:
			b.logger.Info("", zap.Any("event", event))
		}
	}
}

func (b *bot) processQueueChannel(ctx context.Context, wg *sync.WaitGroup) {
	defer wg.Done()
	for {
		select {
		case <-ctx.Done():
			return
		case e := <-b.processQueue:
			s := strings.Fields(e.Message)
			p := slackPkg.GetRegistry()[strings.ToLower(s[0])]
			r, err := p.ProcessMessage(b.backend, e)
			if err != nil {
				b.logger.Error(err.Error())
				continue
			}
			b.responseQueue <- r
		}
	}
}

func (b *bot) responseQueueChannel(ctx context.Context, wg *sync.WaitGroup) {
	defer wg.Done()
	for {
		select {
		case <-ctx.Done():
			return
		case r := <-b.responseQueue:
			if r != nil {
				_, _, err := slackPkg.ApiClient().PostMessage(r.ChannelID, slack.MsgOptionText(r.Message, false), slack.MsgOptionTS(r.Timestamp))
				if err != nil {
					b.logger.Error("failed to post message", zap.Error(err))
				}
			}
		}
	}
}
