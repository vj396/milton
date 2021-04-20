package interrupts

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/vj396/milton/src/slack"
	"github.com/vj396/milton/src/types"
)

type queueCmd struct{}

func init() {
	q := new(queueCmd)
	err := slack.Register("!queue", q)
	if err != nil {
		fmt.Fprintf(os.Stderr, "could not register slack command, err: %+v", err)
	}
}

func (qc *queueCmd) ProcessMessage(b types.Backend, message *types.MessageMetadata) *types.MessageMetadata {
	fields := strings.Fields(message.Message)
	if len(fields) == 1 {
		return qc.getChannelQueue(b, message)
	}
	return qc.updateChannelQueue(b, message)
}

func (qc *queueCmd) getChannelQueue(b types.Backend, m *types.MessageMetadata) *types.MessageMetadata {
	i := types.Interrupt{
		ChannelId: m.ChannelID,
	}
	elements, err := b.GetInterruptRecordsForChannel(&i)
	if err != nil {
		m.Message = "error getting queue for channel"
		return m
	}
	if elements == nil {
		m.Message = "queue for the channel is empty"
		return m
	}
	var payload []string
	for idx := range elements {
		p := fmt.Sprintf("%+s : In Queue Since: %+s", elements[idx].Item, time.Since(time.Unix(elements[idx].SubmittedAt, 0)))
		payload = append(payload, p)
	}
	m.Message = strings.Join(payload, "\n")
	return m
}

func (qc *queueCmd) updateChannelQueue(b types.Backend, m *types.MessageMetadata) *types.MessageMetadata {
	fields := strings.Fields(m.Message)
	fields = fields[1:]
	m.Message = ""
	var errs []error
	for idx := range fields {
		if strings.HasPrefix(fields[idx], "//") {
			m.Message += " Comments: "
			m.Message += strings.Join(fields[idx:], " ")
			break
		}
		i := &types.Interrupt{
			SubmittedBy: m.UserID,
			Item:        fields[idx],
			SubmittedAt: time.Now().Unix(),
			ChannelId:   m.ChannelID,
		}
		err := b.CreateInterruptRecord(i)
		if err != nil {
			errs = append(errs, err)
			continue
		}
		m.Message += fmt.Sprintf("%+s added to queue. ", fields[idx])
	}
	if errs != nil {
		m.Message += " Errors:"
		m.Message += fmt.Sprintf(" %+v\n", errs)
	}
	return m
}

func (qc *queueCmd) Usage() string {
	return "<!queue D123456> || <!queue D123456 D234567 https://example.com/D345678> || <!queue D123 //comment>"
}
