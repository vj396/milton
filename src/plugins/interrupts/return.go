package interrupts

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/vj396/milton/src/slack"
	"github.com/vj396/milton/src/types"
)

type returnCmd struct{}

func init() {
	r := new(returnCmd)
	err := slack.Register("!return", r)
	if err != nil {
		fmt.Fprintf(os.Stderr, "could not register slack command, err: %+v", err)
	}
}

func (rc *returnCmd) ProcessMessage(b types.Backend, m *types.MessageMetadata) *types.MessageMetadata {
	fields := strings.Fields(m.Message)
	if len(fields) == 1 {
		m.Message = "not enough arguments provided"
		return m
	}
	fields = fields[1:]
	m.Message = ""
	message := ""
	work := make(map[string][]string)
	var errs []error
	for idx := range fields {
		if strings.HasPrefix(fields[idx], "//") {
			message += "\nComments: "
			message += strings.Join(fields[idx:], " ")
			break
		}
		i := &types.Interrupt{
			SubmittedBy: m.UserID,
			Item:        fields[idx],
			SubmittedAt: time.Now().Unix(),
			ChannelId:   m.ChannelID,
		}
		err := b.DeleteInterruptRecord(i)
		if err != nil {
			errs = append(errs, err)
			continue
		}
		if _, ok := work[i.SubmittedBy]; !ok {
			work[i.SubmittedBy] = []string{}
		}
		work[i.SubmittedBy] = append(work[i.SubmittedBy], fmt.Sprintf("%+s returned. ", fields[idx]))
	}
	if errs != nil {
		message += "\nErrors:"
		message += fmt.Sprintf(" %+v\n", errs)
	}
	for u := range work {
		m.Message += fmt.Sprintf("<@%+s>: %+v\n", u, work[u])
	}
	m.Message += message
	return m
}

func (rc *returnCmd) Usage() string {
	return "<!return D123456> || <!return D123456 D234567 https://example.com/D345678> || <!return D1234 D3242 //comment>"
}
